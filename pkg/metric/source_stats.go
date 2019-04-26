package metric

// SrvMap 应用拓扑
type SrvMap struct {
	AppType      int16                        // 本服务服务类型
	Parents      map[string]*Parent           // 展示所有父节点
	UnknowParent *UnknowParent                // 未接入监控的请求者
	Targets      map[int16]map[string]*Target // 子节点拓扑图
}

// Parent 父节点访问子节点信息
type Parent struct {
	Type           int16
	TargetCount    int // 目标应用收到请求总数
	TargetErrCount int // 目标应用内部异常数
}

// NewParent ....
func NewParent() *Parent {
	return &Parent{}
}

// Target ...
type Target struct {
	AccessCount    int   // 访问总数
	AccessErrCount int   // 访问错误数
	AccessDuration int32 // 访问总耗时
}

// NewTarget ...
func NewTarget() *Target {
	return &Target{}
}

// NewSrvMap ...
func NewSrvMap() *SrvMap {
	return &SrvMap{
		UnknowParent: NewUnknowParent(),
		Parents:      make(map[string]*Parent),
		Targets:      make(map[int16]map[string]*Target),
	}
}
