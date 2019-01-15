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

// type Storage interface {
// 	Start() error
// 	Close() error
// 	AgentStore(agentInfo *util.AgentInfo) error
// }

// Storage ...
type Storage struct {
	cql           *gocql.Session
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
		g.L.Warn("CreateSession", zap.String("error", err.Error()))
		return err
	}
	s.cql = session
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
	query := s.cql.Query(
		misc.AgentInsert,
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
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AgentStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// AgentInfoStore ...
func (s *Storage) AgentInfoStore(appName, agentID string, startTime int64, agentInfo []byte) error {

	query := s.cql.Query(
		misc.AgentInofInsert,
		appName,
		agentID,
		startTime,
		string(agentInfo),
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AgentInfoStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}
	return nil
}

// AgentOffline ...
func (s *Storage) AgentOffline(appName, agentID string, startTime, endTime int64, isLive bool) error {
	query := s.cql.Query(
		misc.AgentOfflineInsert,
		appName,
		agentID,
		endTime,
		isLive,
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AgentOffline error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}
	return nil
}

// AppAPIStore ...
func (s *Storage) AppAPIStore(appName string, apiInfo *trace.TApiMetaData) error {
	query := s.cql.Query(
		misc.AppAPIInsert,
		appName,
		apiInfo.ApiId,
		apiInfo.ApiInfo,
		apiInfo.GetLine(),
		apiInfo.GetType(),
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AppAPIStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// APIStore ...
func (s *Storage) APIStore(apiInfo *trace.TApiMetaData) error {

	query := s.cql.Query(
		misc.APIInsert,
		apiInfo.ApiId,
		apiInfo.ApiInfo,
		apiInfo.GetLine(),
		apiInfo.GetType(),
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("APIStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// AppSQLStore ...
func (s *Storage) AppSQLStore(appName string, sqlInfo *trace.TSqlMetaData) error {
	newSQL := g.B64.EncodeToString(talent.String2Bytes(sqlInfo.Sql))
	query := s.cql.Query(
		misc.AppSQLInsert,
		appName,
		sqlInfo.SqlId,
		newSQL,
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AppSQLStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// AppStringStore ...
func (s *Storage) AppStringStore(appName string, strInfo *trace.TStringMetaData) error {
	query := s.cql.Query(
		misc.AgentStrInsert,
		appName,
		strInfo.StringId,
		strInfo.StringValue,
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("AgentStringStore error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
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
					for _, qSpan := range spansQueue {
						if err := s.WriteSpan(qSpan); err != nil {
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

	// @TODO 转码优化
	annotations, _ := json.Marshal(span.GetAnnotations())
	spanEvenlist, _ := json.Marshal(span.GetSpanEventList())
	exceptioninfo, _ := json.Marshal(span.GetExceptionInfo())

	query := s.cql.Query(
		misc.InsertSpan,
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
		annotations,
		span.GetErr(),
		spanEvenlist,
		span.GetParentApplicationName(),
		span.GetParentApplicationType(),
		span.GetAcceptorHost(),
		span.GetApplicationServiceType(),
		exceptioninfo,
		span.GetApiId(),
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("writeSpan error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// writeSpanChunk ...
func (s *Storage) writeSpanChunk(spanChunk *trace.TSpanChunk) error {

	spanEvenlist, _ := json.Marshal(spanChunk.GetSpanEventList())
	query := s.cql.Query(
		misc.InsertSpanChunk,
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
	)
	if err := query.Exec(); err != nil {
		g.L.Warn("writeSpanChunk error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// writeIndexes ...
func (s *Storage) writeIndexes(span *trace.TSpan) error {
	if err := s.appOperationIndex(span); err != nil {
		g.L.Warn("appOperationIndex error", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// appOperationIndex ...
func (s *Storage) appOperationIndex(span *trace.TSpan) error {

	query := s.cql.Query(
		misc.InsertOperIndex,
		span.ApplicationName,
		span.AgentId,
		span.GetApiId(),
		span.StartTime,
		span.TransactionId,
		span.GetRPC(),
		span.GetSpanId(),
	)

	if err := query.Exec(); err != nil {
		g.L.Warn("inster app_operation_index error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// writeAgentStatBatch ....
func (s *Storage) writeAgentStatBatch(appName, agentID string, agentStatBatch *pinpoint.TAgentStatBatch, infoB []byte) error {
	batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)
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

	if err := s.cql.ExecuteBatch(batchInsert); err != nil {
		g.L.Warn("writeAgentStatBatch", zap.String("error", err.Error()), zap.String("SQL", insertAgentStatBatch))
		return err
	}

	return nil
}

// appOperationIndex ...
func (s *Storage) writeAgentStat(appName, agentID string, agentStat *pinpoint.TAgentStat, infoB []byte) error {
	if misc.Conf.Storage.AgentStatUseTTL {
		query := s.cql.Query(
			misc.InsertAgentStatTTL,
			appName,
			agentID,
			agentStat.GetStartTimestamp(),
			agentStat.GetTimestamp(),
			infoB,
			misc.Conf.Storage.AgentStatTTL,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("inster writeAgentStat error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
			return err
		}
	} else {
		query := s.cql.Query(
			misc.InsertAgentStat,
			appName,
			agentID,
			agentStat.GetStartTimestamp(),
			agentStat.GetTimestamp(),
			infoB,
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("inster writeAgentStat error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
			return err
		}
	}
	return nil
}
