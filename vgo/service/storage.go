package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/mafanr/g"
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
	return nil
	query := fmt.Sprintf(util.AgentInsert,
		agentInfo.AppName, agentInfo.AgentID, agentInfo.ServiceType, agentInfo.HostName,
		agentInfo.IP4S, agentInfo.Pid, agentInfo.Version, agentInfo.StartTimestamp, agentInfo.EndTimestamp,
		agentInfo.IsLive, agentInfo.IsContainer, agentInfo.OperatingEnv)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentInfoStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf(util.AgentUpdate, agentInfo.ServiceType, agentInfo.HostName,
			agentInfo.IP4S, agentInfo.Pid, agentInfo.Version, agentInfo.StartTimestamp, agentInfo.EndTimestamp,
			agentInfo.IsLive, agentInfo.IsContainer, agentInfo.OperatingEnv, agentInfo.AppName, agentInfo.AgentID)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentInfoStore ...
func (s *Storage) AgentInfoStore(appName, agentID string, agentInfo []byte) error {
	return nil
	query := fmt.Sprintf(util.AgentInfoInsert, agentInfo, appName, agentID)
	_, err := g.DB.Query(query)
	if err != nil {
		g.L.Warn("AgentInfoStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		return err
	}
	return nil
}

// AgentOffline ...
func (s *Storage) AgentOffline(appName, agentID string, endTime int64, isLive bool) error {
	return nil
	query := fmt.Sprintf(util.AgentOffLine, isLive, endTime, appName, agentID)
	_, err := g.DB.Query(query)
	if err != nil {
		g.L.Warn("AgentOffline:g.DB.Query", zap.String("query", query), zap.Error(err))
		return err
	}
	return nil
}

// AgentAPIStore ...
func (s *Storage) AgentAPIStore(appName, agentID string, apiInfo *trace.TApiMetaData) error {
	return nil
	query := fmt.Sprintf("insert into agent_api (app_name, agent_id, api_id, api_info, agent_start_time) values ('%s','%q','%d','%s',%d);",
		appName, agentID, apiInfo.ApiId, apiInfo.ApiInfo, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentAPIStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_api set api_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and api_id='%d';",
			apiInfo.ApiInfo, apiInfo.AgentStartTime, appName, agentID, apiInfo.ApiId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentAPIStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentSQLStore ...
func (s *Storage) AgentSQLStore(appName, agentID string, apiInfo *trace.TSqlMetaData) error {
	return nil
	newSQL := g.B64.EncodeToString(talent.String2Bytes(apiInfo.Sql))
	query := fmt.Sprintf("insert into agent_sql (app_name, agent_id, sql_id, sql_info, agent_start_time) values ('%s','%s','%d','%q',%d);",
		appName, agentID, apiInfo.SqlId, newSQL, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentSQLStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_sql set sql_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and sql_id='%d';",
			newSQL, apiInfo.AgentStartTime, appName, agentID, apiInfo.SqlId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentSQLStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
	}
	return nil
}

// AgentStringStore ...
func (s *Storage) AgentStringStore(appName, agentID string, apiInfo *trace.TStringMetaData) error {
	return nil
	query := fmt.Sprintf("insert into agent_str (app_name, agent_id, str_id, str_info, agent_start_time) values ('%s','%s','%d','%q',%d);",
		appName, agentID, apiInfo.StringId, apiInfo.StringValue, apiInfo.AgentStartTime)
	if _, err := g.DB.Query(query); err != nil {
		g.L.Warn("AgentStringStore:g.DB.Query", zap.String("query", query), zap.Error(err))
		query = fmt.Sprintf("update  agent_str set str_info='%q', agent_start_time=%d where app_name='%s' and agent_id='%s' and str_id='%d';",
			apiInfo.StringValue, apiInfo.AgentStartTime, appName, agentID, apiInfo.StringId)
		if _, err := g.DB.Query(query); err != nil {
			g.L.Warn("AgentStringStore:g.DB.Query", zap.String("query", query), zap.Error(err))
			return err
		}
		return nil
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
	if !false {
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

	insertAPIID := `
	INSERT
	INTO operation_apis(app_name, agent_id, api_id)
	VALUES (?, ?, ?)`
	if err := s.session.Query(
		insertAPIID,
		span.ApplicationName,
		span.AgentId,
		span.GetApiId(),
	).Exec(); err != nil {
		g.L.Warn("inster operation_apis error", zap.String("error", err.Error()), zap.String("SQL", insertAPIID))
		return err
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
