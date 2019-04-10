package service

import (
	"github.com/shaocongcong/tracing/collector/misc"
	"github.com/shaocongcong/tracing/pkg/proto/pinpoint/thrift/trace"
	"github.com/shaocongcong/tracing/pkg/proto/stats"
	"github.com/shaocongcong/tracing/pkg/proto/ttype"
)

// Stats 计算结果
type Stats struct {
	APIStats        *stats.APIStats        // api计算统计
	MethodStats     *stats.MethodStats     // 接口计算统计
	SQLStats        *stats.SQLStats        // sql语句计算统计
	ExceptionsStats *stats.ExceptionsStats // 异常计算统计
}

// newStats ....
func newStats() *Stats {
	return &Stats{
		APIStats:        stats.NewAPIStats(),
		MethodStats:     stats.NewMethodStats(),
		SQLStats:        stats.NewSQLStats(),
		ExceptionsStats: stats.NewExceptionsStats(),
	}
}

// counter 计算
func (s *Stats) spanCounter(span *trace.TSpan, srvMap *stats.SrvMapStats, apiCall *stats.APICallStats) error {
	// 计算API信息
	{
		s.apiCounter(span)
	}
	// 计算API被哪些服务调用
	{
		apiCallCounter(apiCall, span)
	}

	// 计算服务拓扑图
	{
		srvMapCounter(srvMap, span)
	}
	// 计算method、sql、exceptions
	{
		s.eventsCounter(span, srvMap)
	}

	return nil
}

// apiCallCounter 接口被哪些服务调用计算
func apiCallCounter(apiCall *stats.APICallStats, span *trace.TSpan) {
	api, ok := apiCall.APIS[span.GetApiId()]
	if !ok {
		api = stats.NewAPI()
		apiCall.APIS[span.GetApiId()] = api
	}
	parent, ok := api.Parents[span.GetParentApplicationName()]
	if !ok {
		parent = stats.NewParentInfo()
		parent.Type = span.GetParentApplicationType()
		api.Parents[span.GetParentApplicationName()] = parent
	}
	parent.Count++
	parent.Totalelapsed += span.Elapsed
	if span.GetErr() != 0 {
		parent.ErrCount++
	}
}

// srvMapCounter 计算服务拓扑图
func srvMapCounter(srvMap *stats.SrvMapStats, span *trace.TSpan) {
	srvMap.AppType = span.GetServiceType()
	srv, ok := srvMap.SrvMaps[span.GetParentApplicationName()]
	if !ok {
		srv = stats.NewParentInfo()
		srv.Type = span.GetParentApplicationType()
		srvMap.SrvMaps[span.GetParentApplicationName()] = srv
	}

	srv.Count++
	srv.Totalelapsed += span.Elapsed
	if span.GetErr() != 0 {
		srv.ErrCount++
	}
}

// 计算sql拓扑图
func sqlMapCounter(srvMap *stats.SrvMapStats, event *trace.TSpanEvent) {
	isDB := false
	if event.ServiceType == ttype.MYSQL_EXECUTE_QUERY {
		isDB = true
	} else if event.ServiceType == ttype.REDIS {
		isDB = true
	} else if event.ServiceType == ttype.ORACLE_EXECUTE_QUERY {
		isDB = true
	} else if event.ServiceType == ttype.POSTGRESQL_EXECUTE_QUERY {
		isDB = true
	}
	if isDB {
		dbInfo, ok := srvMap.DBMaps[event.ServiceType]
		if !ok {
			dbInfo = stats.NewDBInfo()
			srvMap.DBMaps[event.ServiceType] = dbInfo
		}

		dbInfo.Count++
		if event.GetExceptionInfo() != nil {
			dbInfo.ErrCount++
		}
		dbInfo.Totalelapsed += event.EndElapsed
	}
}

// apiCounter 计算api信息
func (s *Stats) apiCounter(span *trace.TSpan) {
	apiInfo, ok := s.APIStats.Get(span.GetRPC())
	if !ok {
		apiInfo = stats.NewAPIInfo()
		s.APIStats.Store(span.GetRPC(), apiInfo)
	}

	apiInfo.TotalElapsed += span.Elapsed
	apiInfo.Count++

	// 耗时小于满意时间满意次数加1
	if span.Elapsed < misc.Conf.Stats.SatisfactionTime {
		apiInfo.SatisfactionCount++
		// 耗时小于可容忍时间，可容忍次数加一， 其他都为沮丧次数
	} else if span.Elapsed < misc.Conf.Stats.TolerateTime {
		apiInfo.TolerateCount++
	}
	// 当前时间大于最大时间，更新最大耗时
	if span.Elapsed > apiInfo.MaxElapsed {
		apiInfo.MaxElapsed = span.Elapsed
	}
	// 最小耗时为0或者小于最小耗时，更新最小耗时
	if apiInfo.MinElapsed == 0 || apiInfo.MinElapsed > span.Elapsed {
		apiInfo.MinElapsed = span.Elapsed
	}

	// 获取是否有错误
	if span.GetErr() != 0 {
		apiInfo.ErrCount++
	}
}

// counterEvents 计算method、SQL信息
func (s *Stats) eventsCounter(span *trace.TSpan, srvMap *stats.SrvMapStats) {
	if len(s.MethodStats.APIStr) == 0 {
		s.MethodStats.APIStr = span.GetRPC()
	}
	for _, event := range span.GetSpanEventList() {
		{
			isErr := false
			// 是否有异常抛出
			if event.GetExceptionInfo() != nil {
				isErr = true
			}
			// 计算method
			s.methodCount(event.GetApiId(), event.EndElapsed, isErr)
			// 计算sql
			annotations := event.GetAnnotations()
			for _, annotation := range annotations {
				// 20为数据库类型
				if annotation.GetKey() == 20 {
					s.sqlCount(annotation.Value.GetIntStringStringValue().GetIntValue(), event.EndElapsed, isErr)
				}
			}
			// 异常计算
			s.exceptionCount(event.GetApiId(), event)
		}
	}
}

func (s *Stats) methodCount(apiID int32, elapsed int32, isErr bool) {
	methodInfo, ok := s.MethodStats.Get(apiID)
	if !ok {
		methodInfo = stats.NewMethodInfo()
		s.MethodStats.Store(apiID, methodInfo)
	}

	methodInfo.TotalElapsed += elapsed
	methodInfo.Count++

	if elapsed > methodInfo.MaxElapsed {
		methodInfo.MaxElapsed = elapsed
	}

	if methodInfo.MinElapsed == 0 || methodInfo.MinElapsed > elapsed {
		methodInfo.MinElapsed = elapsed
	}

	// 是否有异常抛出
	if isErr {
		methodInfo.ErrCount++
	}
}

// sqlCount 计算sql
func (s *Stats) sqlCount(sqlID int32, elapsed int32, isErr bool) {
	sqlInfo, ok := s.SQLStats.Get(sqlID)
	if !ok {
		sqlInfo = stats.NewSQLInfo()
		s.SQLStats.Store(sqlID, sqlInfo)
	}

	sqlInfo.TotalElapsed += elapsed
	sqlInfo.Count++

	if elapsed > sqlInfo.MaxElapsed {
		sqlInfo.MaxElapsed = elapsed
	}

	if sqlInfo.MinElapsed == 0 || sqlInfo.MinElapsed > elapsed {
		sqlInfo.MinElapsed = elapsed
	}

	// 是否有异常抛出
	if isErr {
		sqlInfo.ErrCount++
	}
}

// exceptionCount 异常统计
func (s *Stats) exceptionCount(apiID int32, event *trace.TSpanEvent) {
	// 参看是否存在异常，不存在直接返回
	exInfo := event.GetExceptionInfo()
	if exInfo == nil {
		return
	}

	apiEx, ok := s.ExceptionsStats.Get(apiID)
	if !ok {
		apiEx = stats.NewAPIExceptions()
		s.ExceptionsStats.Store(apiID, apiEx)
	}

	ex, ok := apiEx.Exceptions[exInfo.GetStringValue()]
	if !ok {
		ex = stats.NewExceptionInfo()
		apiEx.Exceptions[exInfo.GetStringValue()] = ex
	}

	ex.TotalElapsed += event.GetEndElapsed()
	ex.Type = int(event.GetServiceType())
	ex.Count++

	if event.GetEndElapsed() > ex.MaxElapsed {
		ex.MaxElapsed = event.GetEndElapsed()
	}

	if ex.MinElapsed == 0 || ex.MinElapsed > event.GetEndElapsed() {
		ex.MinElapsed = event.GetEndElapsed()
	}
}

// OrderlyKey 排序工具
type OrderlyKey []int64

// Len OrderlyKey 长度
func (o OrderlyKey) Len() int {
	return len(o)
}

// Swap 交换
func (o OrderlyKey) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less 对比
func (o OrderlyKey) Less(i, j int) bool {
	return o[i] < o[j]
}
