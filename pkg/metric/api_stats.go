package metric

// APIStats api统计
type APIStats struct {
	APIs map[string]*APIInfo
}

// Get 获取Api信息
func (a *APIStats) Get(apiStr string) (*APIInfo, bool) {
	info, ok := a.APIs[apiStr]
	return info, ok
}

// Store 存储Api信息
func (a *APIStats) Store(apiStr string, info *APIInfo) {
	a.APIs[apiStr] = info
}

// NewAPIStats ...
func NewAPIStats() *APIStats {
	return &APIStats{
		APIs: make(map[string]*APIInfo),
	}
}

// APIInfo 统计信息
type APIInfo struct {
	Duration          int32 // 总耗时
	Count             int   // 发生次数
	ErrCount          int   // 错误次数
	MinDuration       int32 // 最小耗时
	MaxDuration       int32 // 最大耗时
	SatisfactionCount int   // 满意次数
	TolerateCount     int   // 可容忍次数
}

// NewAPIInfo ...
func NewAPIInfo() *APIInfo {
	return &APIInfo{}
}
