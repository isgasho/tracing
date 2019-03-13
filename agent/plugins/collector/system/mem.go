package system

import (
	"fmt"
	"time"

	"github.com/imdevlab/tracing/agent/misc"
	"github.com/imdevlab/tracing/agent/service"
	"github.com/imdevlab/tracing/util"
)

// MemStats  ...
type MemStats struct {
	ps        PS
	stop      chan bool
	appname   string
	Interval  int    `toml:"interval"`
	TransName string `toml:"trans_name"`
}

// Gather ...
func (mem *MemStats) Gather() ([]*util.Metric, error) {
	vm, err := mem.ps.VMStat()
	if err != nil {
		return nil, fmt.Errorf("error getting virtual memory info: %s", err)
	}

	if vm.Total == 0 {
		return nil, fmt.Errorf("mem total is zeor")
	}

	fields := map[string]interface{}{
		"available":         int64(vm.Available),
		"available_percent": int64(100 * float64(vm.Available) / float64(vm.Total)),
	}

	metric := &util.Metric{
		Name:     "mem",
		Fields:   fields,
		Time:     time.Now().Unix(),
		Interval: misc.Conf.System.Interval,
		Tags:     make(map[string]string),
	}

	return []*util.Metric{metric}, nil
}

// Init ...
func (mem *MemStats) Init() error {

	return nil
}

func init() {
	service.AddCollector("mem", &MemStats{
		stop: make(chan bool, 1),
		ps:   &systemPS{},
	})
}
