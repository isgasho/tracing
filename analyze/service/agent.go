package service

import "sync"

// Agent ...
type Agent struct {
	sync.RWMutex
	AgentID       string
	startTime     int64
	lastPointTime int64
	isLive        bool
}

// NewAgent ...
func NewAgent(agentID string) *Agent {
	return &Agent{
		AgentID: agentID,
	}
}
