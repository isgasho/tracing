package service

import (
	"log"
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
							g.L.Warn("jvmStore:msgpack.Unmarshal", zap.String("error", err.Error()))
							continue
						}

						batchInsert.Query(`INSERT INTO jvms (app_name, instance_id, report_time, value) VALUES (?,?,?,?)`, value.AppName, value.InstanceID, value.Time, body)
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
						g.L.Warn("jvmStore:msgpack.Unmarshal", zap.String("error", err.Error()))
						continue
					}
					batchInsert.Query(`INSERT INTO jvms (app_name, instance_id, report_time, value) VALUES (?,?,?,?)`, value.AppName, value.InstanceID, value.Time, body)
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

// spanStore ...
func (storage *Storage) spanStore() {
	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
	var spanQueue []*util.Span
	for {
		select {
		case spans, ok := <-storage.spansC:
			if ok {
				spanQueue = append(spanQueue, spans...)
				if len(spanQueue) > misc.Conf.Storage.SpanStoreLen {
					// 插入
					batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
					for _, value := range spanQueue {
						batchInsert.Query(`INSERT INTO traces (trace_id, span_id, app_id, instance_id, span_type, span_layer, start_time, end_time, parent_span_id, operation_id, is_error, refs, tags, logs) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
							value.TraceID, value.SpanID, value.AppID, value.InstanceID, value.SpanType, value.SpanLayer, value.StartTime, value.EndTime, value.ParentSpanID, value.OperationNameID, value.IsError, value.Refs, value.Tags, value.Logs)
						// TraceID         string                `msg:"tid"`
						//           int32                 `msg:"sid"`
						//            int32                 `msg:"aid"`
						//       int32                 `msg:"inid"`
						//         SpanType              `msg:"sty"`
						//        SpanLayer             `msg:"sly"`
						//             []*SpanRef            `msg:"rfs"`
						//        int64                 `msg:"st"`
						//          int64                 `msg:"et"`
						//     int32                 `msg:"pid"`
						//  int32                 `msg:"oid"`
						//          bool                  `msg:"ie"`
						//             []*KeyWithStringValue `msg:"tags"`
						//             []*LogMessage         `msg:"logs"`

						log.Println("TraceID", value.TraceID)
						log.Println("AppID", value.AppID)
						log.Println("InstanceID", value.InstanceID)
						log.Println("SpanType", value.SpanType)
						log.Println("SpanLayer", value.SpanLayer)
						log.Println("Refs", value.Refs)
						log.Println("StartTime", value.StartTime)
						log.Println("EndTime", value.EndTime)
						log.Println("ParentSpanID", value.ParentSpanID)
						log.Println("OperationNameID", value.OperationNameID)
						log.Println("IsError", value.IsError)
						log.Println("Tags", value.Tags)
						log.Println("Logs", value.Logs)
					}
					if err := storage.session.ExecuteBatch(batchInsert); err != nil {
						g.L.Warn("spanStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
					}
					// 清空缓存
					spanQueue = spanQueue[:0]
				}
			}
			break
		case <-ticker.C:
			if len(spanQueue) > 0 {
				// 插入
				batchInsert := storage.session.NewBatch(gocql.UnloggedBatch)
				for _, value := range spanQueue {
					batchInsert.Query(`INSERT INTO traces (trace_id, span_id, app_id, instance_id, span_type, span_layer, start_time, end_time, parent_span_id, operation_id, is_error, refs, tags, logs) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
						value.TraceID, value.SpanID, value.AppID, value.InstanceID, int32(value.SpanType), int32(value.SpanLayer), value.StartTime, value.EndTime, value.ParentSpanID, value.OperationNameID, value.IsError, value.Refs, value.Tags, value.Logs)
					log.Println("TraceID", value.TraceID)
					log.Println("AppID", value.AppID)
					log.Println("InstanceID", value.InstanceID)
					log.Println("SpanType", value.SpanType)
					log.Println("SpanLayer", value.SpanLayer)
					log.Println("Refs", value.Refs)
					log.Println("StartTime", value.StartTime)
					log.Println("EndTime", value.EndTime)
					log.Println("ParentSpanID", value.ParentSpanID)
					log.Println("OperationNameID", value.OperationNameID)
					log.Println("IsError", value.IsError)
					log.Println("Tags", value.Tags)
					log.Println("Logs", value.Logs)

					for _, tag := range value.Tags {
						log.Println("tag", tag.Key, tag.Value)
					}
				}
				if err := storage.session.ExecuteBatch(batchInsert); err != nil {
					g.L.Warn("spanStore:storage.session.ExecuteBatch", zap.String("error", err.Error()))
				}
				// 清空缓存
				spanQueue = spanQueue[:0]
			}
			break
		}
	}
}
