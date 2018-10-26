package util

import "sync"

// App ...
type App struct {
	AppID  int32  `db:"id" json:"id" msg:"id"`
	Name   string `db:"name" json:"name" msg:"name"`
	Agents sync.Map
	Apis   sync.Map
}

// NewApp ...
func NewApp() *App {
	return &App{}
}

// AgentInfo ...
type AgentInfo struct {
	AppID        int32  `db:"app_id" json:"app_id"  msg:"app_id"`
	InstanceID   int32  `db:"instance_id" json:"instance_id" msg:"instance_id"`
	AgentUUID    string `db:"agent_uuid" json:"agent_uuid" msg:"agent_uuid"`
	AppName      string `db:"app_name" json:"app_name"  msg:"app_id"`
	OsName       string `db:"os_name" json:"os_name" msg:"os_name"`
	Ipv4S        string `db:"ipv4s" json:"ipv4s" msg:"ipv4s"`
	RegisterTime int64  `db:"register_time" json:"register_time" msg:"register_time"`
	ProcessID    int32  `db:"process_id" json:"process_id" msg:"process_id"`
	HostName     string `db:"host_name" json:"host_name" msg:"host_name"`
}

// NewAgentInfo ...
func NewAgentInfo() *AgentInfo {
	return &AgentInfo{}
}
