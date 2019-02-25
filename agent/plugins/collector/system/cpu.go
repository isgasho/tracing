package system

import (
	"fmt"
	"time"

	"github.com/imdevlab/vgo/agent/misc"
	"github.com/imdevlab/vgo/agent/service"

	"github.com/imdevlab/vgo/util"

	"github.com/shirou/gopsutil/cpu"
)

// CPUStats ....
type CPUStats struct {
	ps        PS
	lastStats []cpu.TimesStat
	stop      chan bool
	PerCPU    bool   `toml:"percpu"`
	TotalCPU  bool   `toml:"totalcpu"`
	Interval  int    `toml:"interval"`
	TransName string `toml:"trans_name"`
	appname   string
}

// NewCPUStats ...
func NewCPUStats(ps PS) *CPUStats {
	return &CPUStats{
		ps: ps,
	}
}

// Gather ....
func (cpu *CPUStats) Gather() ([]*util.Metric, error) {
	times, err := cpu.ps.CPUTimes(cpu.PerCPU, cpu.TotalCPU)
	if err != nil {
		return nil, fmt.Errorf("error getting CPU info: %s", err)
	}
	now := time.Now()

	metrics := make([]*util.Metric, 0)

	for i, cts := range times {
		tags := map[string]string{
			"cpu": cts.CPU,
		}

		total := totalCPUTime(cts)
		fields := map[string]interface{}{}
		// Add in percentage
		if len(cpu.lastStats) == 0 {
			metric := &util.Metric{
				Name:     "cpu",
				Tags:     tags,
				Fields:   fields,
				Time:     now.Unix(),
				Interval: misc.Conf.System.Interval,
			}

			metrics = append(metrics, metric)
			continue
		}
		lastCts := cpu.lastStats[i]
		lastTotal := totalCPUTime(lastCts)
		totalDelta := total - lastTotal

		if totalDelta < 0 {
			cpu.lastStats = times
			return nil, fmt.Errorf("Error: current total CPU time is less than previous total CPU time")
		}

		if totalDelta == 0 {
			continue
		}

		fields["usage_user"] = 100 * (cts.User - lastCts.User) / totalDelta
		fields["usage_idle"] = 100 * (cts.Idle - lastCts.Idle) / totalDelta
		fields["usage_iowait"] = 100 * (cts.Iowait - lastCts.Iowait) / totalDelta
		metric := &util.Metric{
			Name:     "cpu",
			Tags:     tags,
			Fields:   fields,
			Time:     now.Unix(),
			Interval: misc.Conf.System.Interval,
		}
		metrics = append(metrics, metric)
	}
	cpu.lastStats = times

	return metrics, nil
}

// Init ...
func (cpu *CPUStats) Init() error {

	return nil
}

func totalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal +
		t.Guest + t.GuestNice + t.Idle
	return total
}

func init() {
	service.AddCollector("cpu", &CPUStats{
		stop:     make(chan bool, 1),
		PerCPU:   true,
		TotalCPU: true,
		ps:       &systemPS{},
	})
}
