package metric

// Info 统计信息
type Info struct {
	TotalElapsed int32 // 总耗时
	Count        int   // 发生次数
	ErrCount     int   // 错误次数
	MinElapsed   int32 // 最小耗时
	MaxElapsed   int32 // 最大耗时
}

// NewInfo ...
func NewInfo() *Info {
	return &Info{}
}
