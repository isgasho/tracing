package storage

import (
	"github.com/gocql/gocql"
	"github.com/imdevlab/tracing/alert/misc"
	"go.uber.org/zap"
)

// Storage 存储
type Storage struct {
	cql    *gocql.Session
	logger *zap.Logger
}

// NewStorage 新建存储
func NewStorage(logger *zap.Logger) *Storage {
	return &Storage{
		logger: logger,
	}
}

// init 初始化存储
func (s *Storage) init() error {
	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
	cluster.Keyspace = misc.Conf.Storage.Keyspace
	cluster.Consistency = gocql.Quorum
	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = misc.Conf.Storage.NumConns

	session, err := cluster.CreateSession()
	if err != nil {
		s.logger.Warn("create session", zap.String("error", err.Error()))
		return err
	}
	s.cql = session
	return nil
}
