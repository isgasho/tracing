package metric

// SQLStats 接口计算统计
type SQLStats struct {
	SQLs map[int32]*SQLInfo
}

// NewSQLStats ...
func NewSQLStats() *SQLStats {
	return &SQLStats{
		SQLs: make(map[int32]*SQLInfo),
	}
}

// Get 获取sql信息
func (s *SQLStats) Get(sqlID int32) (*SQLInfo, bool) {
	info, ok := s.SQLs[sqlID]
	return info, ok
}

// Store 存储sql信息
func (s *SQLStats) Store(sqlID int32, info *SQLInfo) {
	s.SQLs[sqlID] = info
}

// SQLInfo 统计信息
type SQLInfo struct {
	TotalElapsed int32 // 总耗时
	Count        int   // 发生次数
	ErrCount     int   // 错误次数
	MinElapsed   int32 // 最小耗时
	MaxElapsed   int32 // 最大耗时
}

// NewSQLInfo ...
func NewSQLInfo() *SQLInfo {
	return &SQLInfo{}
}
