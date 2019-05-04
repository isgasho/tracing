package metric

// MethodStats 接口计算统计
type MethodStats struct {
	// sync.RWMutex
	APIStr  string                // 触发api接口的method
	Methods map[int32]*MethodInfo // method信息
}

// NewMethodStats ...
func NewMethodStats() *MethodStats {
	return &MethodStats{
		Methods: make(map[int32]*MethodInfo),
	}
}

// Get 获取medthod信息
func (m *MethodStats) Get(apiID int32) (*MethodInfo, bool) {
	// m.RLock()
	info, ok := m.Methods[apiID]
	// m.RUnlock()
	return info, ok
}

// Store 存储method信息
func (m *MethodStats) Store(apiID int32, info *MethodInfo) {
	// m.Lock()
	m.Methods[apiID] = info
	// m.Unlock()
}

// MethodInfo 统计信息
type MethodInfo struct {
	Type        int   // 服务类型
	Duration    int32 // 总耗时
	Count       int   // 发生次数
	ErrCount    int   // 错误次数
	MinDuration int32 // 最小耗时
	MaxDuration int32 // 最大耗时
}

// NewMethodInfo ...
func NewMethodInfo() *MethodInfo {
	return &MethodInfo{}
}
