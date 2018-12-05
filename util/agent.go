package util

// AgentInfo ...
type AgentInfo struct {
	AppName        string `db:"app_name" json:"applicationName"  msg:"applicationName"`
	ServiceType    int32  `db:"ser_type" json:"serviceType"  msg:"serviceType"`
	SocketID       int32  `db:"socket_id" json:"socketId" msg:"socketId"`
	HostName       string `db:"host_name" json:"hostName" msg:"hostName"`
	AgentID        string `db:"agent_id" json:"agentId"  msg:"agentId"`
	IP4S           string `db:"ip" json:"ip" msg:"ip"`
	Pid            int32  `db:"pid" json:"pid" msg:"pid"`
	Version        string `db:"version" json:"version" msg:"version"`
	StartTimestamp int64  `db:"start_time" json:"startTimestamp" msg:"startTimestamp"`
	EndTimestamp   int64  `db:"end_time" json:"end_time" msg:"end_time"`
	IsLive         bool   `db:"is_live" json:"is_live" msg:"is_live"`
	IsContainer    bool   `db:"is_container" json:"is_container" msg:"is_container"`
}

// NewAgentInfo ...
func NewAgentInfo() *AgentInfo {
	return &AgentInfo{}
}
