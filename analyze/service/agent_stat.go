package service

import (
	"fmt"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
	"github.com/mafanr/vgo/proto/pinpoint/thrift"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/pinpoint"
	"go.uber.org/zap"
)

// AgentStats ...
type AgentStats struct {
	agentID string
	cpus    map[int64]*CPULoad
	memorys map[int64]*JvmMemory
}

// NewAgentStats ...
func NewAgentStats() *AgentStats {
	return &AgentStats{
		cpus:    make(map[int64]*CPULoad),
		memorys: make(map[int64]*JvmMemory),
	}
}

func (agentStats *AgentStats) statsCounter(agentID string, index int64, stats []*pinpoint.TAgentStat) error {
	if len(agentStats.agentID) == 0 {
		agentStats.agentID = agentID
	}
	for _, stat := range stats {
		cpu, ok := agentStats.cpus[index]
		if !ok {
			cpu = NewCPULoad()
			agentStats.cpus[index] = cpu
		}

		memory, ok := agentStats.memorys[index]
		if !ok {
			memory = NewJvmMemory()
			agentStats.memorys[index] = memory
		}

		cpu.Jvm += stat.CpuLoad.GetJvmCpuLoad()
		cpu.System += stat.CpuLoad.GetSystemCpuLoad()
		cpu.count++
		memory.HeapUsed += stat.Gc.GetJvmMemoryHeapUsed()
		memory.NonHeap += stat.Gc.GetJvmMemoryNonHeapUsed()
		memory.count++
	}

	return nil
}

// sqlRecord ...
func (agentStats *AgentStats) statRecord(app *App, recordTime int64) error {

	for index, cpu := range agentStats.cpus {
		query := gAnalyze.cql.Session.Query(misc.InsertCPULoadRecord,
			app.AppName,
			agentStats.agentID,
			index,
			cpu.Jvm,
			cpu.System,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("statRecord error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		}
	}

	for index, memory := range agentStats.memorys {
		query := gAnalyze.cql.Session.Query(misc.InsertJVMMemoryRecord,
			app.AppName,
			agentStats.agentID,
			index,
			memory.HeapUsed,
			memory.NonHeap,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("statRecord error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		}
	}
	return nil
}

// CPULoad ...
type CPULoad struct {
	Jvm    float64
	System float64
	count  int
}

// NewCPULoad ...
func NewCPULoad() *CPULoad {
	return &CPULoad{}
}

// JvmMemory ...
type JvmMemory struct {
	HeapUsed int64
	NonHeap  int64
	count    int
}

// NewJvmMemory ...
func NewJvmMemory() *JvmMemory {
	return &JvmMemory{}
}

// statsCounter ...
func statsCounter(app *App, startTime, endTime int64, es map[int64]*Element) error {
	for _, agent := range app.Agents {
		iterAgentStat := gAnalyze.appStore.cql.Session.Query(misc.QueryAgentStat, app.AppName, agent.AgentID, startTime, endTime).Iter()
		var timestamp int64
		var statInfo []byte
		for iterAgentStat.Scan(&timestamp, &statInfo) {
			index, _ := ModMs2Min(timestamp)
			if e, ok := es[index]; ok {
				tStruct := thrift.Deserialize(statInfo)
				switch m := tStruct.(type) {
				case *pinpoint.TAgentStat:
					e.stats.statsCounter(agent.AgentID, index, []*pinpoint.TAgentStat{m})
					break
				case *pinpoint.TAgentStatBatch:
					e.stats.statsCounter(agent.AgentID, index, m.AgentStats)
					break
				default:
					g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", m)))
				}
			}
		}

		if err := iterAgentStat.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}
	for key, e := range es {
		cpu, ok := e.stats.cpus[key]
		if ok {
			if cpu != nil {
				if cpu.count != 0 {
					cpu.Jvm = cpu.Jvm / float64(cpu.count)
					cpu.System = cpu.System / float64(cpu.count)
				}
			}
		}
		memcury, ok := e.stats.memorys[key]
		if ok {
			if memcury != nil {
				if memcury.count != 0 {
					memcury.HeapUsed = memcury.HeapUsed / int64(memcury.count)
					memcury.NonHeap = memcury.NonHeap / int64(memcury.count)
				}
			}
		}
	}
	return nil
}
