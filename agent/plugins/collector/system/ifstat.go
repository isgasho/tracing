package system

import (
	"fmt"
	"time"

	"github.com/imdevlab/vgo/agent/misc"
	"github.com/imdevlab/vgo/agent/service"
	"github.com/imdevlab/vgo/util"
	"github.com/shirou/gopsutil/net"
)

// IfStat ...
type IfStat struct {
	ps              PS
	LastIfStats     map[string]net.IOCountersStat
	LastCollectTime time.Time
	skipChecks      bool
	Interfaces      []string
	stop            chan bool
	appname         string
	Interval        int `toml:"interval"`
}

// Gather ...
func (s *IfStat) Gather() ([]*util.Metric, error) {
	netio, err := s.ps.NetIO()
	if err != nil {
		return nil, fmt.Errorf("error getting netif info: %s", err)
	}

	metrics := make([]*util.Metric, 0)

	now := time.Now()

	// first time,just record the stats
	if s.LastIfStats == nil {
		s.LastIfStats = make(map[string]net.IOCountersStat)
		for _, v := range netio {
			s.LastIfStats[v.Name] = v
		}
		s.LastCollectTime = now
	}

	for _, io := range netio {
		// 采集指定网卡的信息
		if len(misc.Conf.Ifstat.Interfaces) != 0 {
			var found bool
			for _, name := range misc.Conf.Ifstat.Interfaces {
				if name == io.Name {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		tags := map[string]string{
			"interface": io.Name,
		}

		lio := s.LastIfStats[io.Name]
		duration := now.Sub(s.LastCollectTime).Seconds()

		if duration == 0 {
			continue
		}
		fields := map[string]interface{}{
			"out_bytes": float64(io.BytesSent-lio.BytesSent) / duration,
			"in_bytes":  float64(io.BytesRecv-lio.BytesRecv) / duration,

			"out_packets": float64(io.PacketsSent-lio.PacketsSent) / duration,
			"in_packets":  float64(io.PacketsRecv-lio.PacketsRecv) / duration,

			"out_errors": float64(io.Errout-lio.Errout) / duration,
			"in_errors":  float64(io.Errin-lio.Errin) / duration,

			"out_dropped": float64(io.Dropout-lio.Dropout) / duration,
			"in_dropped":  float64(io.Dropin-lio.Dropin) / duration,
		}

		metric := &util.Metric{
			Name:     "ifstat",
			Fields:   fields,
			Time:     time.Now().Unix(),
			Interval: misc.Conf.System.Interval,
			Tags:     tags,
		}

		metrics = append(metrics, metric)

	}

	for _, v := range netio {
		s.LastIfStats[v.Name] = v
	}
	s.LastCollectTime = now

	return metrics, nil
}

// Init ...
func (s *IfStat) Init() error {

	return nil
}

func init() {
	service.AddCollector("ifstat", &IfStat{
		stop: make(chan bool, 1),
		ps:   &systemPS{},
	})
}
