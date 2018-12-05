package storage

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"go.uber.org/zap"

	"github.com/mafanr/vgo/vgo/misc"

	"github.com/gocql/gocql"
)

// Storage ...
type Storage struct {
	session *gocql.Session
	// jvmC    chan *util.JVMS
	// spansC  chan []*util.Span
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		// jvmC:   make(chan *util.JVMS, misc.Conf.Storage.JVMCacheLen),
		// spansC: make(chan []*util.Span, misc.Conf.Storage.SpanCacheLen),
	}
}

// Start ...
func (storage *Storage) Start() error {
	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
	cluster.Keyspace = misc.Conf.Storage.Keyspace
	cluster.Consistency = gocql.Quorum
	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = misc.Conf.Storage.NumConns

	session, err := cluster.CreateSession()
	if err != nil {
		g.L.Fatal("Start:cluster.CreateSession", zap.String("error", err.Error()))
	}
	storage.session = session
	return nil
}

// Close ...
func (storage *Storage) Close() error {
	if storage.session != nil {
		storage.session.Close()
	}
	return nil
}

// AgentStore ...
func (storage *Storage) AgentStore(agentInfo *util.AgentInfo) error {
	if err := storage.session.Query(util.AgentInfoInsert,
		agentInfo.AppName,
		agentInfo.AgentID,
		agentInfo.ServiceType,
		agentInfo.SocketID,
		agentInfo.HostName,
		agentInfo.IP4S,
		agentInfo.Pid,
		agentInfo.Version,
		agentInfo.StartTimestamp,
		agentInfo.IsLive,
		agentInfo.IsContainer,
		agentInfo.EndTimestamp).Exec(); err != nil {
		g.L.Warn("agentStore:storage.session.Query", zap.String("error", err.Error()), zap.String("sql", util.AgentInfoInsert))
		return err
	}
	return nil
}
