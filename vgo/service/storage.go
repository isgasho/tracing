package service

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
	jvmC    chan *util.JVMS
	spansC  chan []*util.Span
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		jvmC:   make(chan *util.JVMS, misc.Conf.Storage.JVMCacheLen),
		spansC: make(chan []*util.Span, misc.Conf.Storage.SpanCacheLen),
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

	go storage.spanStore()

	return nil
}

// Close ...
func (storage *Storage) Close() error {
	if storage.session != nil {
		storage.session.Close()
	}
	return nil
}

// agentStore ...
func (storage *Storage) agentStore(agentInfo *util.AgentInfo) error {
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

func (storage *Storage) jvmStore() {
	// ticker := time.NewTicker(time.Duration(misc.Conf.Storage.JVMStoreInterval) * time.Millisecond)
	// var jvmsQueue []*util.JVMS
	// for {
	// 	select {
	// 	case jvms, ok := <-storage.jvmC:
	// 		if ok {
	// 			jvmsQueue = append(jvmsQueue, jvms)
	// 			if len(jvmsQueue) > misc.Conf.Storage.JVMStoreLen {
	// 				// 插入
	// 				batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
	// 				for _, value := range jvmsQueue {
	// 					batchInsert.Query(util.JVMInsert, value.AppName, value.InstanceID, value.Time, value.JVMs)
	// 				}
	// 				if err := storage.session.ExecuteBatch(batchInsert); err != nil {
	// 					g.L.Warn("jvmStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
	// 				}
	// 				// 清空缓存
	// 				jvmsQueue = jvmsQueue[:0]
	// 			}
	// 		}
	// 		break
	// 	case <-ticker.C:
	// 		if len(jvmsQueue) > 0 {
	// 			// 插入
	// 			batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
	// 			for _, value := range jvmsQueue {
	// 				batchInsert.Query(util.JVMInsert, value.AppName, value.InstanceID, value.Time, value.JVMs)
	// 			}
	// 			if err := storage.session.ExecuteBatch(batchInsert); err != nil {
	// 				g.L.Warn("jvmStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
	// 			}
	// 			// 清空缓存
	// 			jvmsQueue = jvmsQueue[:0]
	// 		}
	// 		break
	// 	}
	// }
}

// spanStore ...
func (storage *Storage) spanStore() {
	// ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
	// var spanQueue []*util.Span
	// for {
	// 	select {
	// 	case spans, ok := <-storage.spansC:
	// 		if ok {
	// 			spanQueue = append(spanQueue, spans...)
	// 			if len(spanQueue) > misc.Conf.Storage.SpanStoreLen {
	// 				// 插入
	// 				batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
	// 				for _, value := range spanQueue {
	// 					batchInsert.Query(util.SpanInsert,
	// 						value.TraceID,
	// 						value.TraceSegmentID,
	// 						value.SpanID,
	// 						value.AppID,
	// 						value.InstanceID,
	// 						int32(value.SpanType),
	// 						int32(value.SpanLayer),
	// 						value.StartTime,
	// 						value.EndTime,
	// 						value.ParentSpanID,
	// 						value.OperationNameID,
	// 						value.IsError,
	// 						value.Refs,
	// 						value.Tags,
	// 						value.Logs)
	// 					log.Println("入库的地方 value.TraceID", value.TraceID)
	// 					log.Println("入库的地方 value.AppID", value.AppID)
	// 					log.Println("入库的地方 value.InstanceID", value.InstanceID)
	// 				}
	// 				if err := storage.session.ExecuteBatch(batchInsert); err != nil {
	// 					g.L.Warn("spanStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
	// 				}
	// 				// 清空缓存
	// 				spanQueue = spanQueue[:0]
	// 			}
	// 		}
	// 		break
	// 	case <-ticker.C:
	// 		if len(spanQueue) > 0 {
	// 			// 插入
	// 			batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
	// 			for _, value := range spanQueue {
	// 				batchInsert.Query(util.SpanInsert,
	// 					value.TraceID,
	// 					value.TraceSegmentID,
	// 					value.SpanID,
	// 					value.AppID,
	// 					value.InstanceID,
	// 					int32(value.SpanType),
	// 					int32(value.SpanLayer),
	// 					value.StartTime,
	// 					value.EndTime,
	// 					value.ParentSpanID,
	// 					value.OperationNameID,
	// 					value.IsError,
	// 					value.Refs,
	// 					value.Tags,
	// 					value.Logs)

	// 				log.Println("----------------------------------------------------")
	// 				log.Println("------------------- start ---------------------------------")
	// 				log.Println("----------------------------------------------------")
	// 				log.Println("入库的地方 value.TraceID", value.TraceID)
	// 				log.Println("入库的地方 value.SpanID", value.SpanID)
	// 				log.Println("入库的地方 value.AppID", value.AppID)
	// 				log.Println("入库的地方 value.InstanceID", value.InstanceID)
	// 				log.Println("--------------------- end -------------------------------")
	// 				log.Println("----------------------------------------------------")
	// 				log.Println("----------------------------------------------------")
	// 			}
	// 			if err := storage.session.ExecuteBatch(batchInsert); err != nil {
	// 				g.L.Warn("spanStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
	// 			}
	// 			// 清空缓存
	// 			spanQueue = spanQueue[:0]
	// 		}
	// 		break
	// 	}
	// }
}
