package service

import (
	"github.com/mafanr/g"
	"go.uber.org/zap"

	"github.com/mafanr/vgo/vgo/misc"

	"github.com/gocql/gocql"
)

// Storage ...
type Storage struct {
	session *gocql.Session
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{}
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
	// log.Println("start")
	// t := time.Now()
	// for i := 0; i < 100; i++ {
	// 	q := session.Query(`INSERT INTO bigrow (rowname, iplist) VALUES (?,?)`, fmt.Sprintf("name_%d", i), fmt.Sprintf("ip_%d", i))
	// 	q.Exec()
	// }
	// log.Println("end", time.Now().Sub(t).Nanoseconds())
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
