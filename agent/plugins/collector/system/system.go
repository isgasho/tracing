package system

import (
	"time"

	"github.com/imdevlab/tracing/agent/misc"
	"github.com/imdevlab/tracing/agent/service"

	"github.com/imdevlab/tracing/util"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
)

// SystemStats ...
type SystemStats struct {
	stop      chan bool
	appname   string
	Interval  int    `toml:"interval"`
	TransName string `toml:"trans_name"`
}

// Gather ...
func (sys *SystemStats) Gather() ([]*util.Metric, error) {
	loadavg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	hostinfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	fields := map[string]interface{}{
		"load1":  loadavg.Load1,
		"uptime": int64(hostinfo.Uptime),
	}

	metric := &util.Metric{
		Name:     "system",
		Fields:   fields,
		Time:     time.Now().Unix(),
		Interval: misc.Conf.System.Interval,
		Tags:     make(map[string]string),
	}

	return []*util.Metric{metric}, nil
}

// Init ...
func (sys *SystemStats) Init() error {
	return nil
}

func init() {
	service.AddCollector("system", &SystemStats{
		stop: make(chan bool, 1),
	})
}
