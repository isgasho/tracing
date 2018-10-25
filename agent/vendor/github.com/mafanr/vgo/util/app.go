package util

import (
	"sync"
)

// App ...
type App struct {
	sync.RWMutex
	Code   int32  `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Agents map[int32]*AgentInfo
}

// NewApp ...
func NewApp() *App {
	return &App{
		Agents: make(map[int32]*AgentInfo),
	}
}

// AgentInfo ...
type AgentInfo struct {
	ID           int32  `db:"id" json:"id" msg:"id"`
	AgentUUID    string `db:"agent_uuid" json:"agent_uuid" msg:"aid"`
	AppCode      int32  `db:"app_code" json:"app_code"  msg:"ac"`
	AppName      string `db:"app_name" json:"app_name"  msg:"an"`
	OsName       string `db:"os_name" json:"os_name" msg:"on"`
	Ipv4S        string `db:"ipv4s" json:"ipv4s" msg:"ipv4s"`
	RegisterTime int64  `db:"register_time" json:"register_time" msg:"rt"`
	ProcessID    int32  `db:"process_id" json:"process_id" msg:"pid"`
	HostName     string `db:"host_name" json:"host_name" msg:"hn"`
}

// NewAgentInfo ...
func NewAgentInfo() *AgentInfo {
	return &AgentInfo{}
}
