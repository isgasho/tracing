package stats

import (
	"github.com/imdevlab/tracing/pkg/metric"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
)

// RuntimeStats runtime 计算
type RuntimeStats struct {
	Type     int              // runtime 类型
	JVMStats *metric.JVMStats // jvm
}

// NewRuntimeStats ...
func NewRuntimeStats() *RuntimeStats {
	return &RuntimeStats{
		JVMStats: metric.NewJVMStats(),
	}
}

// JVMCounter ...
func (r *RuntimeStats) JVMCounter(agentState *pinpoint.TAgentStat) error {
	jvm, ok := r.JVMStats.Agents[agentState.GetAgentId()]
	if !ok {
		jvm = metric.NewJVMInfo()
		r.JVMStats.Agents[agentState.GetAgentId()] = jvm
	}

	jvm.CPU.Count++
	jvm.CPU.Jvm += agentState.CpuLoad.GetJvmCpuLoad()
	jvm.CPU.System += agentState.CpuLoad.GetSystemCpuLoad()

	jvm.Memory.Count++
	jvm.Memory.HeapUsed += agentState.Gc.GetJvmMemoryHeapUsed()
	jvm.Memory.NonHeap += agentState.Gc.GetJvmMemoryNonHeapUsed()
	return nil
}
