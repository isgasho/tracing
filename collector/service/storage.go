package service

// import (
// 	"encoding/json"
// 	"time"

// 	"github.com/gocql/gocql"
// 	"github.com/imdevlab/g"
// 	"github.com/imdevlab/tracing/collector/misc"
// 	"github.com/imdevlab/tracing/pkg/network"
// 	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
// 	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
// 	"github.com/imdevlab/tracing/pkg/sql"
// 	"github.com/sunface/talent"
// 	"go.uber.org/zap"
// )

// // Storage 存储
// type Storage struct {
// 	cql           *gocql.Session
// 	spanChan      chan *trace.TSpan
// 	spanChunkChan chan *trace.TSpanChunk
// 	// metricsChan   chan *util.MetricData
// }

// // newStorage 新建存储
// func newStorage() *Storage {
// 	return &Storage{
// 		spanChan:      make(chan *trace.TSpan, misc.Conf.Storage.SpanCacheLen+500),
// 		spanChunkChan: make(chan *trace.TSpanChunk, misc.Conf.Storage.SpanChunkCacheLen+500),
// 		// metricsChan:   make(chan *util.MetricData, misc.Conf.Storage.MetricCacheLen+500),
// 	}
// }

// // init 初始化存储
// func (s *Storage) init() error {
// 	// connect to the cluster
// 	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
// 	cluster.Keyspace = misc.Conf.Storage.Keyspace
// 	cluster.Consistency = gocql.Quorum
// 	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
// 	cluster.NumConns = misc.Conf.Storage.NumConns

// 	session, err := cluster.CreateSession()
// 	if err != nil {
// 		g.L.Warn("create session", zap.String("error", err.Error()))
// 		return err
// 	}
// 	s.cql = session
// 	return nil
// }

// // Start ...
// func (s *Storage) Start() error {
// 	if err := s.init(); err != nil {
// 		g.L.Warn("storage init", zap.String("error", err.Error()))
// 		return err
// 	}

// 	go s.spanStore()
// 	go s.spanChunkStore()
// 	// go s.systemStore()
// 	return nil
// }

// // Close ...
// func (s *Storage) Close() error {
// 	return nil
// }

// // AgentStore agent信息存储
// func (s *Storage) AgentStore(agentInfo *network.AgentInfo) error {
// 	query := s.cql.Query(
// 		sql.InsertAgent,
// 		agentInfo.AppName,
// 		agentInfo.AgentID,
// 		agentInfo.ServiceType,
// 		agentInfo.HostName,
// 		agentInfo.IP4S,
// 		agentInfo.StartTimestamp,
// 		agentInfo.EndTimestamp,
// 		agentInfo.IsLive,
// 		agentInfo.IsContainer,
// 		agentInfo.OperatingEnv,
// 		misc.Conf.Collector.Addr,
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("agent store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}

// 	return nil
// }

// // AgentInfoStore ...
// func (s *Storage) AgentInfoStore(appName, agentID string, startTime int64, agentInfo []byte) error {
// 	query := s.cql.Query(
// 		sql.InsertAgentInfo,
// 		appName,
// 		agentID,
// 		startTime,
// 		string(agentInfo),
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("agent info store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // AgentOffline ...
// func (s *Storage) AgentOffline(appName, agentID string, startTime, endTime int64, isLive bool) error {
// 	// query := s.cql.Query(
// 	// 	misc.AgentOfflineInsert,
// 	// 	appName,
// 	// 	agentID,
// 	// 	endTime,
// 	// 	isLive,
// 	// )
// 	// if err := query.Exec(); err != nil {
// 	// 	g.L.Warn("AgentOffline error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
// 	// 	return err
// 	// }
// 	return nil
// }

// // AppMethodStore ...
// func (s *Storage) AppMethodStore(appName string, apiInfo *trace.TApiMetaData) error {
// 	query := s.cql.Query(
// 		sql.InsertMethod,
// 		appName,
// 		apiInfo.ApiId,
// 		apiInfo.ApiInfo,
// 		apiInfo.GetLine(),
// 		apiInfo.GetType(),
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("api store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // AppSQLStore sql语句存储，sql语句需要base64转码，防止sql注入
// func (s *Storage) AppSQLStore(appName string, sqlInfo *trace.TSqlMetaData) error {
// 	newSQL := g.B64.EncodeToString(talent.String2Bytes(sqlInfo.Sql))
// 	query := s.cql.Query(
// 		sql.InsertSQL,
// 		appName,
// 		sqlInfo.SqlId,
// 		newSQL,
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("sql store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}

// 	return nil
// }

// // AppStringStore ...
// func (s *Storage) AppStringStore(appName string, strInfo *trace.TStringMetaData) error {
// 	query := s.cql.Query(
// 		sql.InsertString,
// 		appName,
// 		strInfo.StringId,
// 		strInfo.StringValue,
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("string store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // spanStore ...
// func (s *Storage) spanStore() {
// 	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
// 	var spansQueue []*trace.TSpan
// 	for {
// 		select {
// 		case span, ok := <-s.spanChan:
// 			if ok {
// 				spansQueue = append(spansQueue, span)
// 				if len(spansQueue) >= misc.Conf.Storage.SpanCacheLen {
// 					// 插入
// 					for _, qSpan := range spansQueue {
// 						if err := s.WriteSpan(qSpan); err != nil {
// 							g.L.Warn("write span", zap.String("error", err.Error()))
// 							continue
// 						}
// 					}
// 					// 清空缓存
// 					spansQueue = spansQueue[:0]
// 				}
// 			}
// 			break
// 		case <-ticker.C:
// 			if len(spansQueue) > 0 {
// 				// 插入
// 				for _, span := range spansQueue {
// 					if err := s.WriteSpan(span); err != nil {
// 						g.L.Warn("write span", zap.String("error", err.Error()))
// 						continue
// 					}
// 				}
// 				// 清空缓存
// 				spansQueue = spansQueue[:0]
// 			}
// 			break
// 		}
// 	}
// }

// // spanChunkStore ...
// func (s *Storage) spanChunkStore() {
// 	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
// 	var spansChunkQueue []*trace.TSpanChunk
// 	for {
// 		select {
// 		case spanChunk, ok := <-s.spanChunkChan:
// 			if ok {
// 				spansChunkQueue = append(spansChunkQueue, spanChunk)
// 				if len(spansChunkQueue) >= misc.Conf.Storage.SpanChunkCacheLen {
// 					// 插入
// 					for _, qSapnChunk := range spansChunkQueue {
// 						if err := s.writeSpanChunk(qSapnChunk); err != nil {
// 							g.L.Warn("write spanChunk", zap.String("error", err.Error()))
// 							continue
// 						}
// 					}
// 					// 清空缓存
// 					spansChunkQueue = spansChunkQueue[:0]
// 				}
// 			}
// 			break
// 		case <-ticker.C:
// 			if len(spansChunkQueue) > 0 {
// 				// 插入
// 				for _, sapnChunk := range spansChunkQueue {
// 					if err := s.writeSpanChunk(sapnChunk); err != nil {
// 						g.L.Warn("write spanChunk", zap.String("error", err.Error()))
// 						continue
// 					}
// 				}
// 				// 清空缓存
// 				spansChunkQueue = spansChunkQueue[:0]
// 			}
// 			break
// 		}
// 	}
// }

// // WriteSpan ...
// func (s *Storage) WriteSpan(span *trace.TSpan) error {
// 	if err := s.writeSpan(span); err != nil {
// 		g.L.Warn("write span", zap.String("error", err.Error()))
// 		return err
// 	}

// 	if err := s.writeIndexes(span); err != nil {
// 		g.L.Warn("write span index", zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // writeSpan ...
// func (s *Storage) writeSpan(span *trace.TSpan) error {

// 	// @TODO 转码优化
// 	annotations, _ := json.Marshal(span.GetAnnotations())
// 	spanEvenlist, _ := json.Marshal(span.GetSpanEventList())
// 	exceptioninfo, _ := json.Marshal(span.GetExceptionInfo())

// 	query := s.cql.Query(
// 		sql.InsertSpan,
// 		span.TransactionId,
// 		span.SpanId,
// 		span.AgentId,
// 		span.ApplicationName,
// 		span.AgentStartTime,
// 		span.ParentSpanId,
// 		span.StartTime,
// 		span.Elapsed,
// 		span.RPC,
// 		span.ServiceType,
// 		span.GetEndPoint(),
// 		span.GetRemoteAddr(),
// 		annotations,
// 		span.GetErr(),
// 		spanEvenlist,
// 		span.GetParentApplicationName(),
// 		span.GetParentApplicationType(),
// 		span.GetAcceptorHost(),
// 		span.GetApplicationServiceType(),
// 		exceptioninfo,
// 		span.GetApiId(),
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("write span", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}

// 	return nil
// }

// // writeSpanChunk ...
// func (s *Storage) writeSpanChunk(spanChunk *trace.TSpanChunk) error {

// 	spanEvenlist, _ := json.Marshal(spanChunk.GetSpanEventList())
// 	query := s.cql.Query(
// 		sql.InsertSpanChunk,
// 		spanChunk.TransactionId,
// 		spanChunk.SpanId,
// 		spanChunk.AgentId,
// 		spanChunk.ApplicationName,
// 		spanChunk.GetServiceType(),
// 		spanChunk.GetEndPoint(),
// 		spanEvenlist,
// 		spanChunk.GetApplicationServiceType(),
// 		spanChunk.GetKeyTime(),
// 		spanChunk.Version,
// 	)
// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("write spanChunk", zap.String("error", err.Error()), zap.String("SQL", query.String()))
// 		return err
// 	}

// 	return nil
// }

// // writeIndexes ...
// func (s *Storage) writeIndexes(span *trace.TSpan) error {
// 	if err := s.appOperationIndex(span); err != nil {
// 		g.L.Warn("appOperationIndex error", zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

// // appOperationIndex ...
// func (s *Storage) appOperationIndex(span *trace.TSpan) error {

// 	query := s.cql.Query(
// 		sql.InsertOperIndex,
// 		span.ApplicationName,
// 		span.AgentId,
// 		span.GetApiId(),
// 		span.StartTime,
// 		span.GetElapsed(),
// 		span.TransactionId,
// 		span.GetRPC(),
// 		span.GetSpanId(),
// 		span.GetErr(),
// 	)

// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("inster app_operation_index error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
// 		return err
// 	}

// 	return nil
// }

// // writeAgentStatBatch ....
// func (s *Storage) writeAgentStatBatch(appName, agentID string, agentStatBatch *pinpoint.TAgentStatBatch, infoB []byte) error {
// 	batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)
// 	var insertAgentStatBatch string
// 	if misc.Conf.Storage.AgentStatUseTTL {
// 		insertAgentStatBatch = `
// 		INSERT
// 		INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
// 		VALUES (?, ?, ?, ?, ?) USING TTL ?;`
// 		for _, stat := range agentStatBatch.AgentStats {
// 			batchInsert.Query(
// 				insertAgentStatBatch,
// 				appName,
// 				agentID,
// 				stat.GetStartTimestamp(),
// 				stat.GetTimestamp(),
// 				infoB,
// 				misc.Conf.Storage.AgentStatTTL)
// 		}
// 	} else {
// 		insertAgentStatBatch = `
// 		INSERT
// 		INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
// 		VALUES (?, ?, ?, ?, ?) ;`
// 		for _, stat := range agentStatBatch.AgentStats {
// 			batchInsert.Query(
// 				insertAgentStatBatch,
// 				appName,
// 				agentID,
// 				stat.GetStartTimestamp(),
// 				stat.GetTimestamp(),
// 				infoB)
// 		}
// 	}

// 	if err := s.cql.ExecuteBatch(batchInsert); err != nil {
// 		g.L.Warn("agent stat batch", zap.String("error", err.Error()), zap.String("SQL", insertAgentStatBatch))
// 		return err
// 	}

// 	return nil
// }

// // appOperationIndex ...
// func (s *Storage) writeAgentStat(appName, agentID string, agentStat *pinpoint.TAgentStat, infoB []byte) error {
// 	if misc.Conf.Storage.AgentStatUseTTL {
// 		query := s.cql.Query(
// 			sql.InsertAgentStatWithTTL,
// 			appName,
// 			agentID,
// 			agentStat.GetStartTimestamp(),
// 			agentStat.GetTimestamp(),
// 			infoB,
// 			misc.Conf.Storage.AgentStatTTL,
// 		)
// 		if err := query.Exec(); err != nil {
// 			g.L.Warn("inster agentstat", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 			return err
// 		}
// 	} else {
// 		query := s.cql.Query(
// 			sql.InsertAgentStat,
// 			appName,
// 			agentID,
// 			agentStat.GetStartTimestamp(),
// 			agentStat.GetTimestamp(),
// 			infoB,
// 		)
// 		if err := query.Exec(); err != nil {
// 			g.L.Warn("inster agentstat", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 			return err
// 		}
// 	}
// 	return nil
// }

// // storeAPI ...
// func (s *Storage) storeAPI(span *trace.TSpan) error {
// 	query := s.cql.Query(
// 		sql.InsertAPIs,
// 		span.ApplicationName,
// 		span.RPC,
// 	)

// 	if err := query.Exec(); err != nil {
// 		g.L.Warn("store api", zap.String("SQL", query.String()), zap.String("error", err.Error()))
// 		return err
// 	}

// 	return nil
// }

// func (s *Storage) storeServiceType() error {
// 	// batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)
// 	// inputServiceType := `
// 	// 	INSERT
// 	// 	INTO service_type(service_type, info)
// 	// 	VALUES (?, ?) ;`
// 	// for svrID, info := range util.ServiceType {
// 	// 	batchInsert.Query(
// 	// 		inputServiceType,
// 	// 		svrID,
// 	// 		info)
// 	// }

// 	// if err := s.cql.ExecuteBatch(batchInsert); err != nil {
// 	// 	g.L.Warn("storeServiceType", zap.String("error", err.Error()), zap.String("SQL", inputServiceType))
// 	// 	return err
// 	// }
// 	return nil
// }

// // systemStore ...
// func (s *Storage) systemStore() {
// 	// ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SystemStoreInterval) * time.Millisecond)
// 	// var metricQueue []*util.MetricData
// 	// for {
// 	// 	select {
// 	// 	case metric, ok := <-s.metricsChan:
// 	// 		if ok {
// 	// 			metricQueue = append(metricQueue, metric)
// 	// 			if len(metricQueue) >= misc.Conf.Storage.MetricCacheLen {
// 	// 				// 插入
// 	// 				if err := s.WriteMetric(metricQueue); err != nil {
// 	// 					g.L.Warn("writeMetric error", zap.String("error", err.Error()))
// 	// 				}
// 	// 				// 清空缓存
// 	// 				metricQueue = metricQueue[:0]
// 	// 			}
// 	// 		}
// 	// 		break
// 	// 	case <-ticker.C:
// 	// 		if len(metricQueue) > 0 {
// 	// 			// 插入
// 	// 			if err := s.WriteMetric(metricQueue); err != nil {
// 	// 				g.L.Warn("writeMetric error", zap.String("error", err.Error()))
// 	// 			}
// 	// 			// 清空缓存
// 	// 			metricQueue = metricQueue[:0]
// 	// 		}
// 	// 		break
// 	// 	}
// 	// }
// }

// // WriteMetric ...
// // func (s *Storage) WriteMetric(metrics []*util.MetricData) error {
// // batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)

// // for _, metric := range metrics {
// // 	b, err := json.Marshal(&metric.Payload)
// // 	if err != nil {
// // 		g.L.Warn("json", zap.String("error", err.Error()), zap.Any("data", metric.Payload))
// // 		continue
// // 	}
// // 	batchInsert.Query(misc.InsertSystems,
// // 		metric.AppName,
// // 		metric.AgentID,
// // 		metric.Time,
// // 		b)

// // }

// // if err := s.cql.ExecuteBatch(batchInsert); err != nil {
// // 	g.L.Warn("insert metric", zap.String("error", err.Error()), zap.String("SQL", misc.InsertSystems))
// // 	return err
// // }
// // 	return nil
// // }