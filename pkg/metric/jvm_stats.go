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
	CPU    *JVMCPULoad `json:"cpuload"`
	Memory *JvmMemory  `json:"memory"`
	GC     *JVMGC      `json:"gc"`
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
	HeapUsed int64 `json:"heapused"`
	NonHeap  int64 `json:"nonheap"`
	Count    int   `json:"count"`
}

// JVMCPULoad ...
type JVMCPULoad struct {
	Jvm    float64 `json:"jvm"`
	System float64 `json:"system"`
	Count  int     `json:"count"`
}

// JVMGC @TODO
type JVMGC struct {
}
