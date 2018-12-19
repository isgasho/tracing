package service

import (
	"encoding/json"
	"time"

	"github.com/gocql/gocql"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/pinpoint"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"github.com/mafanr/vgo/util"
	"github.com/mafanr/vgo/vgo/misc"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

// Storage ...
type Storage struct {
	session       *gocql.Session
	spanChan      chan *trace.TSpan
	spanChunkChan chan *trace.TSpanChunk
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		spanChan:      make(chan *trace.TSpan, misc.Conf.Storage.SpanCacheLen+500),
		spanChunkChan: make(chan *trace.TSpanChunk, misc.Conf.Storage.SpanChunkCacheLen+500),
	}
}

// Init ...
func (s *Storage) Init() error {

	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
	cluster.Keyspace = misc.Conf.Storage.Keyspace
	cluster.Consistency = gocql.Quorum
	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = misc.Conf.Storage.NumConns

	session, err := cluster.CreateSession()
	if err != nil {
		g.L.Warn("Start:cluster.CreateSession", zap.String("error", err.Error()))
		return err
	}
	s.session = session
	return nil
}

// Start ...
func (s *Storage) Start() error {

	go s.spanStore()
	go s.spanChunkStore()

	return nil
}

// Close ...
func (s *Storage) Close() error {
	return nil
}

// AgentStore ...
func (s *Storage) AgentStore(agentInfo *util.AgentInfo) error {

	agentInsert := `INSERT INTO agents (app_name, agent_id, ser_type, socket_id, host_name, ip,
		pid, version, start_time, end_time, is_live, is_container, operating_env) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	if err := s.session.Query(
		agentInsert,
		agentInfo.AppName,
		agentInfo.AgentID,
		agentInfo.ServiceType,
		agentInfo.SocketID,
		agentInfo.HostName,
		agentInfo.IP4S,
		agentInfo.Pid,
		agentInfo.Version,
		agentInfo.StartTimestamp,
		agentInfo.EndTimestamp,
		agentInfo.IsLive,
		agentInfo.IsContainer,
		agentInfo.OperatingEnv,
	).Exec(); err != nil {
		g.L.Warn("AgentStore error", zap.String("error", err.Error()), zap.String("SQL", agentInsert))
		return err
	}

	return nil
}

// AgentInfoStore ...
func (s *Storage) AgentInfoStore(appName, agentID string, startTime int64, agentInfo []byte) error {
	agentInofInsert := `INSERT INTO agents (app_name, agent_id, start_time, agent_info) 
	VALUES ( ?, ?, ?, ?);`
	if err := s.session.Query(
		agentInofInsert,
		appName,
		agentID,
		startTime,
		string(agentInfo),
	).Exec(); err != nil {
		g.L.Warn("AgentInfoStore error", zap.String("error", err.Error()), zap.String("SQL", agentInofInsert))
		return err
	}
	return nil
}

// AgentOffline ...
func (s *Storage) AgentOffline(appName, agentID string, startTime, endTime int64, isLive bool) error {

	agentOfflineInsert := `INSERT INTO agents (app_name, agent_id, start_time, end_time, is_live) 
	VALUES ( ?, ?, ?, ?, ?);`

	if err := s.session.Query(
		agentOfflineInsert,
		appName,
		agentID,
		startTime,
		endTime,
		isLive,
	).Exec(); err != nil {
		g.L.Warn("AgentOffline error", zap.String("error", err.Error()), zap.String("SQL", agentOfflineInsert))
		return err
	}
	return nil
}

// AgentAPIStore ...
func (s *Storage) AgentAPIStore(appName, agentID string, apiInfo *trace.TApiMetaData) error {

	agentAPIInsert := `INSERT INTO agent_apis (app_name, agent_id, api_id, start_time, api_info, line, type) 
	VALUES (?, ?, ?, ?, ?, ?, ?);`
	if err := s.session.Query(
		agentAPIInsert,
		appName,
		agentID,
		apiInfo.ApiId,
		apiInfo.AgentStartTime,
		apiInfo.ApiInfo,
		apiInfo.GetLine(),
		apiInfo.GetType(),
	).Exec(); err != nil {
		g.L.Warn("AgentAPIStore error", zap.String("error", err.Error()), zap.String("SQL", agentAPIInsert))
		return err
	}

	return nil
}

// AgentSQLStore ...
func (s *Storage) AgentSQLStore(appName, agentID string, sqlInfo *trace.TSqlMetaData) error {
	newSQL := g.B64.EncodeToString(talent.String2Bytes(sqlInfo.Sql))
	agentSQLInsert := `INSERT INTO agent_sqls (app_name, agent_id, sql_id, start_time, sql_info) 
	VALUES (?, ?, ?, ?, ?);`
	if err := s.session.Query(
		agentSQLInsert,
		appName,
		agentID,
		sqlInfo.SqlId,
		sqlInfo.AgentStartTime,
		newSQL,
	).Exec(); err != nil {
		g.L.Warn("AgentSQLStore error", zap.String("error", err.Error()), zap.String("SQL", agentSQLInsert))
		return err
	}

	return nil
}

// AgentStringStore ...
func (s *Storage) AgentStringStore(appName, agentID string, strInfo *trace.TStringMetaData) error {
	agentStrInsert := `INSERT INTO agent_strs (app_name, agent_id, str_id, start_time, str_info) 
	VALUES (?, ?, ?, ?, ?);`
	if err := s.session.Query(
		agentStrInsert,
		appName,
		agentID,
		strInfo.StringId,
		strInfo.AgentStartTime,
		strInfo.StringValue,
	).Exec(); err != nil {
		g.L.Warn("AgentStringStore error", zap.String("error", err.Error()), zap.String("SQL", agentStrInsert))
		return err
	}

	return nil
}

// spanStore ...
func (s *Storage) spanStore() {
	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
	var spansQueue []*trace.TSpan
	for {
		select {
		case span, ok := <-s.spanChan:
			if ok {
				spansQueue = append(spansQueue, span)
				if len(spansQueue) >= misc.Conf.Storage.SpanCacheLen {
					// 插入
					for _, qSapn := range spansQueue {
						if err := s.WriteSpan(qSapn); err != nil {
							g.L.Warn("writeSpan error", zap.String("error", err.Error()))
							continue
						}
					}
					// 清空缓存
					spansQueue = spansQueue[:0]
				}
			}
			break
		case <-ticker.C:
			if len(spansQueue) > 0 {
				// 插入
				for _, sapn := range spansQueue {
					if err := s.WriteSpan(sapn); err != nil {
						g.L.Warn("writeSpan error", zap.String("error", err.Error()))
						continue
					}
				}
				// 清空缓存
				spansQueue = spansQueue[:0]
			}
			break
		}
	}
}

// spanChunkStore ...
func (s *Storage) spanChunkStore() {
	ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SpanStoreInterval) * time.Millisecond)
	var spansChunkQueue []*trace.TSpanChunk
	for {
		select {
		case spanChunk, ok := <-s.spanChunkChan:
			if ok {
				spansChunkQueue = append(spansChunkQueue, spanChunk)
				if len(spansChunkQueue) >= misc.Conf.Storage.SpanChunkCacheLen {
					// 插入
					for _, qSapnChunk := range spansChunkQueue {
						if err := s.writeSpanChunk(qSapnChunk); err != nil {
							g.L.Warn("writeSpanChunk error", zap.String("error", err.Error()))
							continue
						}
					}
					// 清空缓存
					spansChunkQueue = spansChunkQueue[:0]
				}
			}
			break
		case <-ticker.C:
			if len(spansChunkQueue) > 0 {
				// 插入
				for _, sapnChunk := range spansChunkQueue {
					if err := s.writeSpanChunk(sapnChunk); err != nil {
						g.L.Warn("writeSpanChunk error", zap.String("error", err.Error()))
						continue
					}
				}
				// 清空缓存
				spansChunkQueue = spansChunkQueue[:0]
			}
			break
		}
	}
}

// WriteSpan ...
func (s *Storage) WriteSpan(span *trace.TSpan) error {
	if err := s.writeSpan(span); err != nil {
		g.L.Warn("WriteSpan error", zap.String("error", err.Error()))
		return err
	}

	if err := s.writeIndexes(span); err != nil {
		g.L.Warn("WriteSpan error", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// writeSpan ...
func (s *Storage) writeSpan(span *trace.TSpan) error {
	insertSpan := `
	INSERT
	INTO traces(trace_id, span_id, agent_id, app_name, agent_start_time, parent_id,
		start_time, elapsed, rpc, service_type, end_point, remote_addr, annotations, err,
		span_event_list, parent_app_name, parent_app_type, acceptor_host, app_service_type, exception_info, api_id)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	annotations, _ := json.Marshal(span.GetAnnotations())
	spanEvenlist, _ := json.Marshal(span.GetSpanEventList())
	exceptioninfo, _ := json.Marshal(span.GetExceptionInfo())

	if err := s.session.Query(
		insertSpan,
		span.TransactionId,
		span.SpanId,
		span.AgentId,
		span.ApplicationName,
		span.AgentStartTime,
		span.ParentSpanId,
		span.StartTime,
		span.Elapsed,
		span.RPC,
		span.ServiceType,
		span.GetEndPoint(),
		span.GetRemoteAddr(),
		annotations, //span.GetAnnotations(), // 转码
		span.GetErr(),
		spanEvenlist, // span.GetSpanEventList(), // 转码
		span.GetParentApplicationName(),
		span.GetParentApplicationType(),
		span.GetAcceptorHost(),
		span.GetApplicationServiceType(),
		exceptioninfo, // span.GetExceptionInfo(), // 转码
		span.GetApiId(),
	).Exec(); err != nil {
		g.L.Warn("writeSpan error", zap.String("error", err.Error()), zap.String("SQL", insertSpan))
		return err
	}

	return nil
}

// writeSpanChunk ...
func (s *Storage) writeSpanChunk(spanChunk *trace.TSpanChunk) error {
	insertSpanChunk := `
	INSERT
	INTO traces_chunk(trace_id, span_id, agent_id, app_name, service_type, end_point,
		span_event_list, app_service_type, key_time, version)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	spanEvenlist, _ := json.Marshal(spanChunk.GetSpanEventList())
	if err := s.session.Query(
		insertSpanChunk,
		spanChunk.TransactionId,
		spanChunk.SpanId,
		spanChunk.AgentId,
		spanChunk.ApplicationName,
		spanChunk.GetServiceType(),
		spanChunk.GetEndPoint(),
		spanEvenlist,
		spanChunk.GetApplicationServiceType(),
		spanChunk.GetKeyTime(),
		spanChunk.Version,
	).Exec(); err != nil {
		g.L.Warn("writeSpanChunk error", zap.String("error", err.Error()), zap.String("SQL", insertSpanChunk))
		return err
	}

	return nil
}

// writeIndexes ...
func (s *Storage) writeIndexes(span *trace.TSpan) error {
	if err := s.saveAppNameAndAPIID(span); err != nil {
		g.L.Warn("saveAppNameAndAPIID error", zap.String("error", err.Error()))
		return err
	}
	if err := s.appOperationIndex(span); err != nil {
		g.L.Warn("appOperationIndex error", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// saveAppNameAndAPIID ...
func (s *Storage) saveAppNameAndAPIID(span *trace.TSpan) error {
	if !gVgo.appStore.checkAndSaveAgent(span.ApplicationName, span.AgentId) {
		insertAppAndAgentID := `
		INSERT
		INTO app_names(app_name, agent_id)
		VALUES (?, ?)`
		if err := s.session.Query(
			insertAppAndAgentID,
			span.ApplicationName,
			span.AgentId,
		).Exec(); err != nil {
			g.L.Warn("inster app_names error", zap.String("error", err.Error()), zap.String("SQL", insertAppAndAgentID))
			return err
		}
	}

	if !gVgo.appStore.checkAndSaveAPIID(span.ApplicationName, span.AgentId, span.GetApiId()) {
		insertAPIID := `
	INSERT
	INTO operation_apis(app_name, agent_id, api_id, start_time)
	VALUES (?, ?, ?, ?)`
		if err := s.session.Query(
			insertAPIID,
			span.ApplicationName,
			span.AgentId,
			span.GetApiId(),
			span.StartTime,
		).Exec(); err != nil {
			g.L.Warn("inster operation_apis error", zap.String("error", err.Error()), zap.String("SQL", insertAPIID))
			return err
		}
	}
	return nil
}

// appOperationIndex ...
func (s *Storage) appOperationIndex(span *trace.TSpan) error {
	insertOperIndex := `
	INSERT
	INTO app_operation_index(app_name, agent_id, api_id, start_time, trace_id)
	VALUES (?, ?, ?, ?, ?)`
	if err := s.session.Query(
		insertOperIndex,
		span.ApplicationName,
		span.AgentId,
		span.GetApiId(),
		span.StartTime,
		span.TransactionId,
	).Exec(); err != nil {
		g.L.Warn("inster app_operation_index error", zap.String("error", err.Error()), zap.String("SQL", insertOperIndex))
		return err
	}

	return nil
}

// writeAgentStatBatch ....
func (s *Storage) writeAgentStatBatch(appName, agentID string, agentStatBatch *pinpoint.TAgentStatBatch, infoB []byte) error {
	batchInsert := s.session.NewBatch(gocql.UnloggedBatch)
	var insertAgentStatBatch string
	if misc.Conf.Storage.AgentStatUseTTL {
		insertAgentStatBatch = `
		INSERT
		INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
		VALUES (?, ?, ?, ?, ?) USING TTL ?;`
		for _, stat := range agentStatBatch.AgentStats {
			batchInsert.Query(
				insertAgentStatBatch,
				appName,
				agentID,
				stat.GetStartTimestamp(),
				stat.GetTimestamp(),
				infoB,
				misc.Conf.Storage.AgentStatTTL)
		}
	} else {
		insertAgentStatBatch = `
		INSERT
		INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
		VALUES (?, ?, ?, ?, ?) ;`
		for _, stat := range agentStatBatch.AgentStats {
			batchInsert.Query(
				insertAgentStatBatch,
				appName,
				agentID,
				stat.GetStartTimestamp(),
				stat.GetTimestamp(),
				infoB)
		}
	}

	if err := s.session.ExecuteBatch(batchInsert); err != nil {
		g.L.Warn("writeAgentStatBatch", zap.String("error", err.Error()), zap.String("SQL", insertAgentStatBatch))
		return err
	}

	return nil
}

// appOperationIndex ...
func (s *Storage) writeAgentStat(appName, agentID string, agentStat *pinpoint.TAgentStat, infoB []byte) error {
	if misc.Conf.Storage.AgentStatUseTTL {
		insertAgentStat := `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?) USING TTL ?;`
		if err := s.session.Query(
			insertAgentStat,
			appName,
			agentID,
			agentStat.GetStartTimestamp(),
			agentStat.GetTimestamp(),
			infoB,
			misc.Conf.Storage.AgentStatTTL,
		).Exec(); err != nil {
			g.L.Warn("inster writeAgentStat error", zap.String("error", err.Error()), zap.String("SQL", insertAgentStat))
			return err
		}
	} else {
		insertAgentStat := `
	INSERT
	INTO agent_stats(app_name, agent_id, start_time, timestamp, stat_info)
	VALUES (?, ?, ?, ?, ?);`
		if err := s.session.Query(
			insertAgentStat,
			appName,
			agentID,
			agentStat.GetStartTimestamp(),
			agentStat.GetTimestamp(),
			infoB,
		).Exec(); err != nil {
			g.L.Warn("inster writeAgentStat error", zap.String("error", err.Error()), zap.String("SQL", insertAgentStat))
			return err
		}
	}
	return nil
}
