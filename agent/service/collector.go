package service

import (
	"sync"

	"github.com/mafanr/vgo/agent/misc"

	"github.com/mafanr/g"
	"go.uber.org/zap"
)

// SystemCollector 系统信息采集
type SystemCollector struct {
}

// NewSystemCollector ...
func NewSystemCollector() *SystemCollector {
	return &SystemCollector{}
}

// Start ...
func (s *SystemCollector) Start() {
	// 系统采集可以关闭
	if !misc.Conf.System.OnOff {
		return
	}
	// 系统信息采集
	for name, collector := range Collectors {
		if err := collector.Init(); err != nil {
			g.L.Fatal("collector Init", zap.Error(err), zap.String("name", name))
		}
		if err := collector.Start(); err != nil {
			g.L.Fatal("collector Start", zap.Error(err), zap.String("name", name))
		}
		g.L.Info("collector start", zap.String("name", name))
	}
}

// Collectors ...
var Collectors map[string]Collector
var clock sync.RWMutex

// Collector ...
type Collector interface {
	Init() error
	Start() error
	Close() error
}

// AddCollector ....
func AddCollector(name string, collector Collector) {
	clock.Lock()
	defer clock.Unlock()

	if Collectors == nil {
		Collectors = make(map[string]Collector)
	}
	Collectors[name] = collector
}
