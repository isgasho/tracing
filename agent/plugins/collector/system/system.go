package system

import (
	"log"
	"time"

	"github.com/mafanr/g"
	"go.uber.org/zap"

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

func (sys *SystemStats) Gather() error {
	loadavg, err := load.Avg()
	if err != nil {
		return err
	}

	hostinfo, err := host.Info()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"load1":  loadavg.Load1,
		"uptime": int64(hostinfo.Uptime),
	}

	log.Println(fields)
	// metric := &system.Metric{
	// 	Name:     "system",
	// 	Fields:   fields,
	// 	Time:     time.Now().Unix(),
	// 	Interval: sys.Interval,
	// 	Tags:     make(map[string]string),
	// }

	// metric.Tags["app"] = sys.appname

	// log.Println("采集4 tags", metric.Tags)
	// log.Println("采集4 fields", metric.Fields)
	// agent.Writer(sys.TransName, []*system.Metric{metric})

	return nil
}

func (sys *SystemStats) start() {
	// sys.appname = agent.Name()
	// ticker := time.NewTicker(time.Duration(sys.Interval) * time.Second)
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		agent.Logger.Error("sys init", zap.Any("err", err))
	// 	}
	// }()
	// defer ticker.Stop()

	// for {
	// 	select {
	// 	case <-sys.stop:
	// 		return nil
	// 	case <-ticker.C:
	// 		sys.Gather()
	// 		continue
	// 	}
	// }

	for {
		time.Sleep(1 * time.Second)
		g.L.Info("plugins", zap.String("name", "system"))
	}

}

// Init ...
func (sys *SystemStats) Init() error {
	go sys.start()
	return nil
}

// Start ...
func (sys *SystemStats) Start() error {
	return nil
}

// Close ...
func (sys *SystemStats) Close() error {
	sys.stop <- true
	return nil
}

func init() {
	service.AddCollector("system", &SystemStats{
		stop: make(chan bool, 1),
	})
}
