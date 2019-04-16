package metric

// ExceptionsStats 异常计算统计
type ExceptionsStats struct {
	MethodEx map[int32]*MethodExceptions
}

// NewExceptionsStats ...
func NewExceptionsStats() *ExceptionsStats {
	return &ExceptionsStats{
		MethodEx: make(map[int32]*MethodExceptions),
	}
}

// Get 获取Method异常信息
func (a *ExceptionsStats) Get(methodID int32) (*MethodExceptions, bool) {
	info, ok := a.MethodEx[methodID]
	return info, ok
}

// Store 存储methodID异常信息
func (a *ExceptionsStats) Store(methodID int32, info *MethodExceptions) {
	a.MethodEx[methodID] = info
}

// MethodExceptions 异常
type MethodExceptions struct {
	Exceptions map[string]*ExceptionInfo
}

// NewAPIExceptions ...
func NewAPIExceptions() *MethodExceptions {
	return &MethodExceptions{
		Exceptions: make(map[string]*ExceptionInfo),
	}
}

// NewExceptionInfo  ....
func NewExceptionInfo() *ExceptionInfo {
	return &ExceptionInfo{}
}

// ExceptionInfo 异常信息
type ExceptionInfo struct {
	Type         int   // 服务类型
	TotalElapsed int32 // 总耗时
	Count        int   // 发生次数
	MinElapsed   int32 // 最小耗时
	MaxElapsed   int32 // 最大耗时
}
