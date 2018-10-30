package service

import (
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"

	"github.com/mafanr/vgo/vgo/misc"

	"github.com/gocql/gocql"
)

// Storage ...
type Storage struct {
	session *gocql.Session
	jvmC    chan *util.JVMS
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		jvmC: make(chan *util.JVMS, misc.Conf.Storage.JVMCacheLen),
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

	go storage.jvmStore()

	return nil
}

// Close ...
func (storage *Storage) Close() error {
	if storage.session != nil {
		storage.session.Close()
	}
	return nil
}

func (storage *Storage) jvmStore() {
	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.JVMStoreInterval) * time.Millisecond)
	var jvmsQueue []*util.JVMS
	for {
		select {
		case jvms, ok := <-storage.jvmC:
			if ok {
				jvmsQueue = append(jvmsQueue, jvms)
				if len(jvmsQueue) > misc.Conf.Storage.JVMStoreLen {
					// 插入
					batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
					for _, value := range jvmsQueue {
						body, err := msgpack.Marshal(value.JVMs)
						if err != nil {
							g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
							continue
						}

						batchInsert.Query(`INSERT INTO jvm (app_name, instance_id, report_time, value) VALUES (?,?,?,?)`, value.AppName, value.InstanceID, value.Time, body)
					}
					if err := storage.session.ExecuteBatch(batchInsert); err != nil {
						g.L.Warn("jvmStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
					}
					// 清空缓存
					jvmsQueue = jvmsQueue[:0]
				}
			}
			break
		case <-ticker.C:
			if len(jvmsQueue) > 0 {
				// 插入
				batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
				for _, value := range jvmsQueue {
					body, err := msgpack.Marshal(value.JVMs)
					if err != nil {
						g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
						continue
					}
					batchInsert.Query(`INSERT INTO jvm (app_name, instance_id, report_time, value) VALUES (?,?,?,?)`, value.AppName, value.InstanceID, value.Time, body)
				}
				if err := storage.session.ExecuteBatch(batchInsert); err != nil {
					g.L.Warn("jvmStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
				}
				// 清空缓存
				jvmsQueue = jvmsQueue[:0]
			}
			break
		}
	}
}
