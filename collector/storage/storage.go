package storage

import (
	"encoding/json"
	"time"

	"github.com/gocql/gocql"
	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/metric"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/sql"
	"github.com/sunface/talent"
	"go.uber.org/zap"
)

// Storage 存储
type Storage struct {
	cql           *gocql.Session
	spanChan      chan *trace.TSpan
	spanChunkChan chan *trace.TSpanChunk
	logger        *zap.Logger
	// metricsChan   chan *util.MetricData
}

// NewStorage 新建存储
func NewStorage(logger *zap.Logger) *Storage {
	return &Storage{
		spanChan:      make(chan *trace.TSpan, misc.Conf.Storage.SpanCacheLen+500),
		spanChunkChan: make(chan *trace.TSpanChunk, misc.Conf.Storage.SpanChunkCacheLen+500),
		logger:        logger,
		// metricsChan:   make(chan *util.MetricData, misc.Conf.Storage.MetricCacheLen+500),
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

// Start ...
func (s *Storage) Start() error {
	if err := s.init(); err != nil {
		s.logger.Warn("storage init", zap.String("error", err.Error()))
		return err
	}

	go s.spanStore()
	go s.spanChunkStore()
	// go s.systemStore()
	return nil
}

// SpanStore span存储
func (s *Storage) SpanStore(span *trace.TSpan) {
	s.spanChan <- span
}

// SpanChunkStore spanChunk存储
func (s *Storage) SpanChunkStore(span *trace.TSpanChunk) {
	s.spanChunkChan <- span
}

// Close ...
func (s *Storage) Close() error {
	return nil
}

// AgentStore agent信息存储
func (s *Storage) AgentStore(agentInfo *network.AgentInfo, islive bool) error {
	query := s.cql.Query(
		sql.InsertAgent,
		agentInfo.AppName,
		agentInfo.AgentID,
		agentInfo.ServiceType,
		agentInfo.HostName,
		agentInfo.IP4S,
		agentInfo.StartTimestamp,
		agentInfo.EndTimestamp,
		agentInfo.IsContainer,
		agentInfo.OperatingEnv,
		misc.Conf.Collector.Addr,
		islive,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("agent store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// UpdateAgentState agent在线状态更新
func (s *Storage) UpdateAgentState(appname string, agentid string, islive bool) error {
	query := s.cql.Query(
		sql.UpdateAgentState,
		islive,
		appname,
		agentid,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("update agent state error", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// AppNameStore 存储Appname
func (s *Storage) AppNameStore(name string) error {
	query := s.cql.Query(
		sql.InsertApp,
		name,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert app name error", zap.String("SQL", query.String()), zap.String("error", err.Error()), zap.String("appName", name))
		return err
	}
	return nil
}

// AgentInfoStore ...
func (s *Storage) AgentInfoStore(appName, agentID string, startTime int64, agentInfo []byte) error {
	query := s.cql.Query(
		sql.InsertAgentInfo,
		appName,
		agentID,
		startTime,
		string(agentInfo),
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("agent info store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}
	return nil
}

// AgentOffline ...
func (s *Storage) AgentOffline(appName, agentID string, startTime, endTime int64, isLive bool) error {
	// query := s.cql.Query(
	// 	misc.AgentOfflineInsert,
	// 	appName,
	// 	agentID,
	// 	endTime,
	// 	isLive,
	// )
	// if err := query.Exec(); err != nil {
	// 	s.logger.Warn("AgentOffline error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
	// 	return err
	// }
	return nil
}

// AppMethodStore ...
func (s *Storage) AppMethodStore(appName string, apiInfo *trace.TApiMetaData) error {
	query := s.cql.Query(
		sql.InsertMethod,
		appName,
		apiInfo.ApiId,
		apiInfo.ApiInfo,
		apiInfo.GetLine(),
		apiInfo.GetType(),
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("api store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}
	return nil
}

// AppSQLStore sql语句存储，sql语句需要base64转码，防止sql注入
func (s *Storage) AppSQLStore(appName string, sqlInfo *trace.TSqlMetaData) error {
	newSQL := g.B64.EncodeToString(talent.String2Bytes(sqlInfo.Sql))
	query := s.cql.Query(
		sql.InsertSQL,
		appName,
		sqlInfo.SqlId,
		newSQL,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("sql store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// AppStringStore ...
func (s *Storage) AppStringStore(appName string, strInfo *trace.TStringMetaData) error {
	query := s.cql.Query(
		sql.InsertString,
		appName,
		strInfo.StringId,
		strInfo.StringValue,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("string store", zap.String("SQL", query.String()), zap.String("error", err.Error()))
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
							s.logger.Warn("write span", zap.String("error", err.Error()))
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
				for _, span := range spansQueue {
					if err := s.WriteSpan(span); err != nil {
						s.logger.Warn("write span", zap.String("error", err.Error()))
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
							s.logger.Warn("write spanChunk", zap.String("error", err.Error()))
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
						s.logger.Warn("write spanChunk", zap.String("error", err.Error()))
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
		s.logger.Warn("write span", zap.String("error", err.Error()))
		return err
	}

	if err := s.writeIndexes(span); err != nil {
		s.logger.Warn("write span index", zap.String("error", err.Error()))
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
		sql.InsertSpan,
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
		s.logger.Warn("write span", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// writeSpanChunk ...
func (s *Storage) writeSpanChunk(spanChunk *trace.TSpanChunk) error {

	spanEvenlist, _ := json.Marshal(spanChunk.GetSpanEventList())
	query := s.cql.Query(
		sql.InsertSpanChunk,
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
		s.logger.Warn("write spanChunk", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// writeIndexes ...
func (s *Storage) writeIndexes(span *trace.TSpan) error {
	if err := s.appOperationIndex(span); err != nil {
		s.logger.Warn("appOperationIndex error", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// appOperationIndex ...
func (s *Storage) appOperationIndex(span *trace.TSpan) error {

	query := s.cql.Query(
		sql.InsertOperIndex,
		span.ApplicationName,
		span.AgentId,
		span.GetApiId(),
		span.StartTime,
		span.GetElapsed(),
		span.TransactionId,
		span.GetRPC(),
		span.GetSpanId(),
		span.GetErr(),
		span.GetRemoteAddr(),
	)

	if err := query.Exec(); err != nil {
		s.logger.Warn("inster app_operation_index error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}

	return nil
}

// WriteAgentStatBatch ....
func (s *Storage) WriteAgentStatBatch(appName, agentID string, agentStatBatch *pinpoint.TAgentStatBatch, infoB []byte) error {
	batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)

	for _, agentStat := range agentStatBatch.AgentStats {
		jvmInfo := metric.NewJVMInfo()
		jvmInfo.CPULoad.Jvm = agentStat.CpuLoad.GetJvmCpuLoad()
		jvmInfo.CPULoad.System = agentStat.CpuLoad.GetSystemCpuLoad()
		jvmInfo.GC.Type = agentStat.Gc.GetType()
		jvmInfo.GC.HeapUsed = agentStat.Gc.GetJvmMemoryHeapUsed()
		jvmInfo.GC.HeapMax = agentStat.Gc.GetJvmMemoryHeapMax()
		jvmInfo.GC.NonHeapUsed = agentStat.Gc.GetJvmMemoryNonHeapUsed()
		jvmInfo.GC.NonHeapMax = agentStat.Gc.GetJvmMemoryHeapMax()
		jvmInfo.GC.GcOldCount = agentStat.Gc.GetJvmGcOldCount()
		jvmInfo.GC.JvmGcOldTime = agentStat.Gc.GetJvmGcOldTime()
		jvmInfo.GC.JvmGcNewCount = agentStat.Gc.GetJvmGcDetailed().GetJvmGcNewCount()
		jvmInfo.GC.JvmGcNewTime = agentStat.Gc.GetJvmGcDetailed().GetJvmGcNewTime()
		jvmInfo.GC.JvmPoolCodeCacheUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolCodeCacheUsed()
		jvmInfo.GC.JvmPoolNewGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolNewGenUsed()
		jvmInfo.GC.JvmPoolOldGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolOldGenUsed()
		jvmInfo.GC.JvmPoolSurvivorSpaceUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolSurvivorSpaceUsed()
		jvmInfo.GC.JvmPoolPermGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolPermGenUsed()
		jvmInfo.GC.JvmPoolMetaspaceUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolMetaspaceUsed()

		body, err := json.Marshal(jvmInfo)
		if err != nil {
			s.logger.Warn("json marshal", zap.String("error", err.Error()))
			continue
		}

		t, err := utils.MSToTime(agentStat.GetTimestamp())
		if err != nil {
			s.logger.Warn("ms to time", zap.Int64("time", agentStat.GetTimestamp()), zap.String("error", err.Error()))
			continue
		}

		batchInsert.Query(
			sql.InsertRuntimeStat,
			appName,
			agentID,
			t.Unix(),
			body,
			1)
	}
	if err := s.cql.ExecuteBatch(batchInsert); err != nil {
		s.logger.Warn("agent stat batch", zap.String("error", err.Error()), zap.String("SQL", sql.InsertRuntimeStat))
		return err
	}

	return nil
}

// WriteAgentStat  ...
func (s *Storage) WriteAgentStat(appName, agentID string, agentStat *pinpoint.TAgentStat, infoB []byte) error {
	jvmInfo := metric.NewJVMInfo()
	jvmInfo.CPULoad.Jvm = agentStat.CpuLoad.GetJvmCpuLoad()
	jvmInfo.CPULoad.System = agentStat.CpuLoad.GetSystemCpuLoad()
	jvmInfo.GC.Type = agentStat.Gc.GetType()
	jvmInfo.GC.HeapUsed = agentStat.Gc.GetJvmMemoryHeapUsed()
	jvmInfo.GC.HeapMax = agentStat.Gc.GetJvmMemoryHeapMax()
	jvmInfo.GC.NonHeapUsed = agentStat.Gc.GetJvmMemoryNonHeapUsed()
	jvmInfo.GC.NonHeapMax = agentStat.Gc.GetJvmMemoryHeapMax()
	jvmInfo.GC.GcOldCount = agentStat.Gc.GetJvmGcOldCount()
	jvmInfo.GC.JvmGcOldTime = agentStat.Gc.GetJvmGcOldTime()
	jvmInfo.GC.JvmGcNewCount = agentStat.Gc.GetJvmGcDetailed().GetJvmGcNewCount()
	jvmInfo.GC.JvmGcNewTime = agentStat.Gc.GetJvmGcDetailed().GetJvmGcNewTime()
	jvmInfo.GC.JvmPoolCodeCacheUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolCodeCacheUsed()
	jvmInfo.GC.JvmPoolNewGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolNewGenUsed()
	jvmInfo.GC.JvmPoolOldGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolOldGenUsed()
	jvmInfo.GC.JvmPoolSurvivorSpaceUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolSurvivorSpaceUsed()
	jvmInfo.GC.JvmPoolPermGenUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolPermGenUsed()
	jvmInfo.GC.JvmPoolMetaspaceUsed = agentStat.Gc.GetJvmGcDetailed().GetJvmPoolMetaspaceUsed()

	body, err := json.Marshal(jvmInfo)
	if err != nil {
		s.logger.Warn("json marshal", zap.String("error", err.Error()))
		return err
	}

	t, err := utils.MSToTime(agentStat.GetTimestamp())
	if err != nil {
		s.logger.Warn("ms to time", zap.Int64("time", agentStat.GetTimestamp()), zap.String("error", err.Error()))
		return err
	}

	query := s.cql.Query(
		sql.InsertRuntimeStat,
		appName,
		agentID,
		t.Unix(),
		body,
		1,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("inster agentstat", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// StoreAPI 存储API信息
func (s *Storage) StoreAPI(span *trace.TSpan) error {
	query := s.cql.Query(
		sql.InsertAPIs,
		span.ApplicationName,
		span.RPC,
	)

	if err := query.Exec(); err != nil {
		s.logger.Warn("store api", zap.String("SQL", query.String()), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// StoreSrvType 存储服务类型
func (s *Storage) StoreSrvType() error {
	batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)
	for svrID, info := range constant.ServiceType {
		batchInsert.Query(
			sql.InsertSrvType,
			svrID,
			info)
	}
	if err := s.cql.ExecuteBatch(batchInsert); err != nil {
		s.logger.Warn("insert server type", zap.String("SQL", sql.InsertSrvType), zap.String("error", err.Error()))
		return err
	}
	return nil
}

// InsertAPIStats ...
func (s *Storage) InsertAPIStats(appName string, inputTime int64, apiStr string, apiInfo *metric.APIInfo) error {
	query := s.cql.Query(sql.InsertAPIStats,
		appName,
		inputTime,
		apiStr,
		apiInfo.TotalElapsed,
		apiInfo.MaxElapsed,
		apiInfo.MinElapsed,
		apiInfo.Count,
		apiInfo.ErrCount,
		apiInfo.SatisfactionCount,
		apiInfo.TolerateCount,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("inster api stats error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		return err
	}

	return nil
}

// InsertMethodStats ...
func (s *Storage) InsertMethodStats(appName string, inputTime int64, apiStr string, methodID int32, methodInfo *metric.MethodInfo) error {
	query := s.cql.Query(sql.InsertMethodStats,
		appName,
		apiStr,
		inputTime,
		methodID,
		methodInfo.Type,
		methodInfo.TotalElapsed,
		methodInfo.MaxElapsed,
		methodInfo.MinElapsed,
		methodInfo.Count,
		methodInfo.ErrCount,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert method error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}
	return nil
}

// InsertExceptionStats ...
func (s *Storage) InsertExceptionStats(appName string, inputTime int64, methodID int32, exceptions map[int32]*metric.ExceptionInfo) error {
	for classID, exinfo := range exceptions {
		query := s.cql.Query(sql.InsertExceptionStats,
			appName,
			methodID,
			classID,
			inputTime,
			exinfo.TotalElapsed,
			exinfo.MaxElapsed,
			exinfo.MinElapsed,
			exinfo.Count,
			exinfo.Type,
		)
		if err := query.Exec(); err != nil {
			s.logger.Warn("insert exception error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
			return err
		}
	}
	return nil
}

// InsertParentMap ...
func (s *Storage) InsertParentMap(appName string, appType int32, inputTime int64, parentName string, parentInfo *metric.ParentInfo) error {
	query := s.cql.Query(sql.InsertParentMap,
		appName,
		inputTime,
		appType,
		parentName,
		parentInfo.Type,
		parentInfo.Count,
		parentInfo.ErrCount,
		parentInfo.Totalelapsed,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert parent map error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		return err
	}

	return nil
}

// InsertChildMap ...
func (s *Storage) InsertChildMap(appName string, appType int32, inputTime int64, childType int32, destinationStr string, destination *metric.Destination) error {
	query := s.cql.Query(sql.InsertChildMap,
		appName,
		inputTime,
		appType,
		childType,
		destinationStr,
		destination.Count,
		destination.ErrCount,
		destination.Totalelapsed,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert child map error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		return err
	}

	return nil
}

// InsertUnknowParentMap ...
func (s *Storage) InsertUnknowParentMap(appName string, appType int32, inputTime int64, unknowParent *metric.UnknowParent) error {
	query := s.cql.Query(sql.InsertUnknowParentMap,
		appName,
		inputTime,
		appType,
		unknowParent.Count,
		unknowParent.ErrCount,
		unknowParent.Totalelapsed,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert unknow parent map error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		return err
	}

	return nil
}

// InsertAPICallStats Api被调用统计信息
func (s *Storage) InsertAPICallStats(appName string, appType int32, inputTime int64, apiID int32, parentname string, parentInfo *metric.ParentInfo) error {
	query := s.cql.Query(sql.InsertAPICallStats,
		appName,
		inputTime,
		appType,
		apiID,
		parentname,
		parentInfo.Count,
		parentInfo.ErrCount,
		parentInfo.Totalelapsed,
	)
	if err := query.Exec(); err != nil {
		s.logger.Warn("insert api call stats error", zap.String("error", err.Error()), zap.String("sql", query.String()))
		return err
	}

	return nil
}

// InsertSQLStats ...
func (s *Storage) InsertSQLStats(appName string, inputTime int64, sqlID int32, sqlInfo *metric.SQLInfo) error {
	query := s.cql.Query(sql.InsertSQLStats,
		appName,
		sqlID,
		inputTime,
		sqlInfo.TotalElapsed,
		sqlInfo.MaxElapsed,
		sqlInfo.MinElapsed,
		sqlInfo.Count,
		sqlInfo.ErrCount,
	)

	if err := query.Exec(); err != nil {
		s.logger.Warn("sql stats insert error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}
	return nil
}

// InsertRuntimeStats ...
func (s *Storage) InsertRuntimeStats(appName string, agentID string, inputTime int64, runtimeType int, info *metric.JVMInfo) error {
	body, err := json.Marshal(info)
	if err != nil {
		s.logger.Warn("InsertRuntimeStats", zap.String("error", err.Error()))
		return err
	}
	query := s.cql.Query(sql.InsertRuntimeStat,
		appName,
		agentID,
		inputTime,
		body,
	)

	if err := query.Exec(); err != nil {
		s.logger.Warn("runtime stats insert error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return err
	}
	return nil
}

// systemStore ...
func (s *Storage) systemStore() {
	// ticker := time.NewTicker(time.Duration(misc.Conf.Storage.SystemStoreInterval) * time.Millisecond)
	// var metricQueue []*util.MetricData
	// for {
	// 	select {
	// 	case metric, ok := <-s.metricsChan:
	// 		if ok {
	// 			metricQueue = append(metricQueue, metric)
	// 			if len(metricQueue) >= misc.Conf.Storage.MetricCacheLen {
	// 				// 插入
	// 				if err := s.WriteMetric(metricQueue); err != nil {
	// 					s.logger.Warn("writeMetric error", zap.String("error", err.Error()))
	// 				}
	// 				// 清空缓存
	// 				metricQueue = metricQueue[:0]
	// 			}
	// 		}
	// 		break
	// 	case <-ticker.C:
	// 		if len(metricQueue) > 0 {
	// 			// 插入
	// 			if err := s.WriteMetric(metricQueue); err != nil {
	// 				s.logger.Warn("writeMetric error", zap.String("error", err.Error()))
	// 			}
	// 			// 清空缓存
	// 			metricQueue = metricQueue[:0]
	// 		}
	// 		break
	// 	}
	// }
}

// WriteMetric ...
// func (s *Storage) WriteMetric(metrics []*util.MetricData) error {
// batchInsert := s.cql.NewBatch(gocql.UnloggedBatch)

// for _, metric := range metrics {
// 	b, err := json.Marshal(&metric.Payload)
// 	if err != nil {
// 		s.logger.Warn("json", zap.String("error", err.Error()), zap.Any("data", metric.Payload))
// 		continue
// 	}
// 	batchInsert.Query(misc.InsertSystems,
// 		metric.AppName,
// 		metric.AgentID,
// 		metric.Time,
// 		b)

// }

// if err := s.cql.ExecuteBatch(batchInsert); err != nil {
// 	s.logger.Warn("insert metric", zap.String("error", err.Error()), zap.String("SQL", misc.InsertSystems))
// 	return err
// }
// 	return nil
// }
