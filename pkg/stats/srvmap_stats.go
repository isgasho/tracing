package stats

// SrvMapStats 应用拓扑
type SrvMapStats struct {
	AppType int16                  // 本服务服务类型
	SrvMaps map[string]*ParentInfo // 父节点拓扑
	DBMaps  map[int16]*DBInfo      // 访问DB拓扑
}

// NewSrvMapStats ...
func NewSrvMapStats() *SrvMapStats {
	return &SrvMapStats{
		SrvMaps: make(map[string]*ParentInfo),
		DBMaps:  make(map[int16]*DBInfo),
	}
}

// ParentInfo 调用信息
type ParentInfo struct {
	Type         int16
	Totalelapsed int32
	Count        int
	ErrCount     int
}

// NewParentInfo ...
func NewParentInfo() *ParentInfo {
	return &ParentInfo{}
}

// DBInfo 数据库信息
type DBInfo struct {
	Totalelapsed int32
	Count        int
	ErrCount     int
}

// NewDBInfo ...
func NewDBInfo() *DBInfo {
	return &DBInfo{}
}
