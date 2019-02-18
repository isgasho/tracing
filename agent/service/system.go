package service

import (
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/agent/misc"
	"github.com/mafanr/vgo/util"
	"go.uber.org/zap"
)

// SystemCollector 系统信息采集
type SystemCollector struct {
	stop chan bool
}

// NewSystemCollector ...
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{
		stop: make(chan bool, 1),
	}
}

// Start ...
func (s *SystemCollector) Start() {
	// 是否启用系统信息采集
	if !misc.Conf.System.OnOff {
		return
	}

	// 初始化系统信息采集
	for name, collector := range Collectors {
		if err := collector.Init(); err != nil {
			g.L.Fatal("collector Init", zap.Error(err), zap.String("name", name))
		}
	}

	// 使用同一个计时器
	ticker := time.NewTicker(time.Duration(misc.Conf.System.Interval) * time.Second)
	defer func() {
		if err := recover(); err != nil {
			g.L.Error("cpu init", zap.Any("err", err))
		}
	}()

	defer ticker.Stop()

	metrics := make([]*util.Metric, 0)

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			for name, collector := range Collectors {
				metric, err := collector.Gather()
				if err != nil {
					g.L.Debug("system collector err", zap.String("name", name), zap.String("error", err.Error()))
					continue
				}
				metrics = append(metrics, metric)
			}

			// 发送

			// 清空缓存
			metrics = metrics[:0]

			continue
		}
	}
}

// Close ....
func (s *SystemCollector) Close() error {
	s.stop <- true
	return nil
}
