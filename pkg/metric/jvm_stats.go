package metric

// JVMStats jvm 信息计算统计
type JVMStats struct {
	Agents map[string]*JVMInfo
}

// NewJVMStats ...
func NewJVMStats() *JVMStats {
	return &JVMStats{
		Agents: make(map[string]*JVMInfo),
	}
}

// JVMInfo ...
type JVMInfo struct {
	CPU    *JVMCPULoad
	Memory *JvmMemory
	GC     *JVMGC
}

// NewJVMInfo ...
func NewJVMInfo() *JVMInfo {
	return &JVMInfo{
		CPU:    &JVMCPULoad{},
		Memory: &JvmMemory{},
		GC:     &JVMGC{},
	}
}

// JvmMemory ...
type JvmMemory struct {
	HeapUsed int64
	NonHeap  int64
	Count    int
}

// JVMCPULoad ...
type JVMCPULoad struct {
	Jvm    float64
	System float64
	Count  int
}

// JVMGC @TODO
type JVMGC struct {
}
