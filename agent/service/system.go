package service

import (
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/vgo/agent/misc"
	"github.com/imdevlab/vgo/util"
	"github.com/vmihailenco/msgpack"
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

	metrics := util.NewMetricData()
	metrics.AppName = gAgent.agentInfo.AppName
	metrics.AgentID = gAgent.agentInfo.AgentID
	packet := &util.VgoPacket{
		Type:       util.TypeOfSystem,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncNo,
		IsCompress: util.TypeOfCompressYes,
	}

	isGetApInfo := false
	if len(metrics.AppName) == 0 || len(metrics.AgentID) == 0 {
		metrics.AppName = gAgent.agentInfo.AppName
		metrics.AgentID = gAgent.agentInfo.AgentID
	}

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			if !isGetApInfo {
				metrics.AppName = gAgent.agentInfo.AppName
				metrics.AgentID = gAgent.agentInfo.AgentID
				if len(metrics.AppName) != 0 && len(metrics.AgentID) != 0 {
					isGetApInfo = true
				}
			}
			// 一次采集所有插件
			for name, collector := range Collectors {
				metric, err := collector.Gather()
				if err != nil {
					g.L.Debug("system collector err", zap.String("name", name), zap.String("error", err.Error()))
					continue
				}
				if metric == nil {
					continue
				}
				// 存放
				metrics.Payload = append(metrics.Payload, metric...)
			}
			metrics.Time = time.Now().Unix()
			// 编码
			payload, err := msgpack.Marshal(metrics)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				// 清空缓存
				metrics.Payload = metrics.Payload[:0]
				continue
			}
			packet.Payload = payload

			if len(metrics.Payload) == 0 {
				continue
			}

			// 发送
			if err := gAgent.client.WritePacket(packet); err != nil {
				g.L.Warn("WritePacket", zap.String("error", err.Error()))
			}

			// 清空缓存
			metrics.Payload = metrics.Payload[:0]
			continue
		}
	}
}

// Close ....
func (s *SystemCollector) Close() error {
	s.stop <- true
	return nil
}
