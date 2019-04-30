package service

// Agent 服务单例
type Agent struct {
	id        string // agent id
	startTime int64  // 启动时间
	isLive    bool
	hostName  string
}

func newAgent(id string, starttime int64) *Agent {
	return &Agent{
		id:        id,
		startTime: starttime,
	}
}
