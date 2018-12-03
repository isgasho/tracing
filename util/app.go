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
	AppName        string `db:"applicationName" json:"applicationName"  msg:"applicationName"`
	ServiceType    int32  `db:"serviceType" json:"serviceType"  msg:"serviceType"`
	SocketID       int32  `db:"socketId" json:"socketId" msg:"socketId"`
	HostName       string `db:"hostName" json:"hostName" msg:"hostName"`
	AgentID        string `db:"agentId" json:"agentId"  msg:"agentId"`
	IP4S           string `db:"ip" json:"ip" msg:"ip"`
	Pid            int32  `db:"pid" json:"pid" msg:"pid"`
	Version        string `db:"version" json:"version" msg:"version"`
	StartTimestamp int64  `db:"startTimestamp" json:"startTimestamp" msg:"startTimestamp"`
}

// NewAgentInfo ...
func NewAgentInfo() *AgentInfo {
	return &AgentInfo{}
}
