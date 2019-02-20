package system

import (
	"log"
	"time"

	"github.com/mafanr/vgo/util"

	"github.com/mafanr/vgo/agent/service"
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

func (sys *SystemStats) Gather() (*util.Metric, error) {
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
		Interval: sys.Interval,
		Tags:     make(map[string]string),
	}

	log.Println("采集 tags", metric.Tags)
	log.Println("采集 fields", metric.Fields)

	return metric, nil
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
