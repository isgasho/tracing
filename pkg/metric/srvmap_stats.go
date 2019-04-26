package metric

// SrvMapStats 应用拓扑
type SrvMapStats struct {
	AppType      int16                  // 本服务服务类型
	UnknowParent *UnknowParent          // 未接入监控的请求者
	Parents      map[string]*ParentInfo // 父节点拓扑
	Childs       map[int16]*Child       // 子节点拓扑图
}

// NewSrvMapStats ...
func NewSrvMapStats() *SrvMapStats {
	return &SrvMapStats{
		UnknowParent: NewUnknowParent(),
		Parents:      make(map[string]*ParentInfo),
		Childs:       make(map[int16]*Child),
	}
}

// ParentInfo 调用信息
type ParentInfo struct {
	Type           int16
	Duration       int32
	Count          int
	AccessErrCount int
	ExceptionCount int
}

// NewParentInfo ...
func NewParentInfo() *ParentInfo {
	return &ParentInfo{}
}

// Child ...
type Child struct {
	Destinations map[string]*Destination
}

// NewChild ....
func NewChild() *Child {
	return &Child{
		Destinations: make(map[string]*Destination),
	}
}

// NewDestination ...
func NewDestination() *Destination {
	return &Destination{}
}

// Destination 目标
type Destination struct {
	Duration       int32
	Count          int
	ExceptionCount int
	AccessErrCount int
}

// UnknowParent 未接入监控的服务，只能抓到访问地址
type UnknowParent struct {
	// Duration       int32
	TargetCount    int
	TargetErrCount int
}

// NewUnknowParent ...
func NewUnknowParent() *UnknowParent {
	return &UnknowParent{}
}
