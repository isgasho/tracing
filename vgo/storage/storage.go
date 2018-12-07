package storage

import "github.com/mafanr/vgo/util"

// Storage ...
type Storage interface {
	Init() error
	AgentStore(agentInfo *util.AgentInfo) error
}
