package stats

// ExceptionsStats 异常计算统计
type ExceptionsStats struct {
	APIEx map[int32]*APIExceptions
}

// NewExceptionsStats ...
func NewExceptionsStats() *ExceptionsStats {
	return &ExceptionsStats{
		APIEx: make(map[int32]*APIExceptions),
	}
}

// Get 获取Api异常信息
func (a *ExceptionsStats) Get(apiID int32) (*APIExceptions, bool) {
	info, ok := a.APIEx[apiID]
	return info, ok
}

// Store 存储Api异常信息
func (a *ExceptionsStats) Store(apiID int32, info *APIExceptions) {
	a.APIEx[apiID] = info
}

// APIExceptions 异常
type APIExceptions struct {
	Exceptions map[string]*ExceptionInfo
}

// NewAPIExceptions ...
func NewAPIExceptions() *APIExceptions {
	return &APIExceptions{
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
