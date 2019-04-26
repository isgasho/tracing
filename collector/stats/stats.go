package stats

import (
	"fmt"
	"strings"

	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/metric"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
)

// Stats 计算结果
type Stats struct {
	APIStats        *metric.APIStats        // api计算统计
	MethodStats     *metric.MethodStats     // 接口计算统计
	SQLStats        *metric.SQLStats        // sql语句计算统计
	ExceptionsStats *metric.ExceptionsStats // 异常计算统计
	RespCodes       map[int]struct{}        //
	SrvMap          *metric.SrvMap          // 服务拓扑图
	// ServerMap       *metric.SrvMapStats     // 服务拓扑图
}

// NewStats ....
func NewStats(respCodes map[int]struct{}) *Stats {
	return &Stats{
		APIStats:        metric.NewAPIStats(),
		MethodStats:     metric.NewMethodStats(),
		SQLStats:        metric.NewSQLStats(),
		ExceptionsStats: metric.NewExceptionsStats(),
		// ServerMap:       metric.NewSrvMapStats(),
		RespCodes: respCodes,
		SrvMap:    metric.NewSrvMap(),
	}
}

// SpanCounter 计算
func (s *Stats) SpanCounter(span *trace.TSpan, apiCall *metric.APICallStats) error {
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
		s.parentMapCounter(span)
	}

	// 计算method、sql、exceptions
	{
		s.eventsCounter(span)
	}

	return nil
}

// SpanChunkCounter counter 计算
func (s *Stats) SpanChunkCounter(spanChunk *trace.TSpanChunk, apiCall *metric.APICallStats) error {
	// 计算method、sql、exceptions
	{
		s.eventsCounterSpanChunk(spanChunk)
	}
	return nil
}

// eventsCounterSpanChunk  计算method、sql、exceptions,
func (s *Stats) eventsCounterSpanChunk(spanChunk *trace.TSpanChunk) {
	// 这里获取不到api str 是否要抛弃该数据包， 或者method就不要放在api这个key下面
	if len(s.MethodStats.APIStr) == 0 {
		// continue
	}

	for _, event := range spanChunk.GetSpanEventList() {
		isErr := false
		// 是否有异常抛出
		if event.GetExceptionInfo() != nil {
			isErr = true
		}
		// 计算method
		s.methodCount(event.GetApiId(), int(event.GetServiceType()), event.EndElapsed, isErr)
		// 计算sql
		annotations := event.GetAnnotations()
		for _, annotation := range annotations {
			// 20为数据库类型
			if annotation.GetKey() == constant.SQL_ID {
				s.sqlCount(annotation.Value.GetIntStringStringValue().GetIntValue(), event.EndElapsed, isErr)
			}
		}
		// 异常计算
		s.exceptionCount(event.GetApiId(), event)
	}
}

// apiCallCounter 接口被哪些服务调用计算
func apiCallCounter(apiCall *metric.APICallStats, span *trace.TSpan) {
	// spanID 为-1的情况该服务就是父节点，查不到被谁调用，这里可以考虑能不能抓到请求者到IP
	if span.ParentSpanId == -1 {
		return
	}

	api, ok := apiCall.APIS[span.GetApiId()]
	if !ok {
		api = metric.NewAPI()
		apiCall.APIS[span.GetApiId()] = api
	}
	parent, ok := api.Parents[span.GetParentApplicationName()]
	if !ok {
		parent = metric.NewParentInfo()
		parent.Type = span.GetParentApplicationType()
		api.Parents[span.GetParentApplicationName()] = parent
	}
	parent.Count++
	parent.Duration += span.Elapsed
	if span.GetErr() != 0 {
		parent.ExceptionCount++
	}
}

// parentMapCounter 计算服务拓扑图
func (s *Stats) parentMapCounter(span *trace.TSpan) {
	// spanID 为-1的情况该服务就是父节点，请求者应该是没接入监控
	if span.ParentSpanId == -1 {
		s.SrvMap.UnknowParent.TargetCount++
		if span.GetErr() != 0 {
			s.SrvMap.UnknowParent.TargetErrCount++
		}
		return
	}

	s.SrvMap.AppType = span.GetServiceType()
	parent, ok := s.SrvMap.Parents[span.GetParentApplicationName()]
	if !ok {
		parent = metric.NewParent()
		parent.Type = span.GetParentApplicationType()
		s.SrvMap.Parents[span.GetParentApplicationName()] = parent
	}

	parent.TargetCount++
	if span.GetErr() != 0 {
		parent.TargetErrCount++
	}
}

func getip(destinationID string) (string, error) {
	strs := strings.Split(destinationID, ":")
	if len(strs) != 2 {
		return "", fmt.Errorf("unknow addr")
	}
	if len(strs[0]) == 0 {
		return "", fmt.Errorf("error ip")
	}
	return strs[0], nil
}

// 计算child拓扑图
// func (s *Stats) childMapCounter(event *trace.TSpanEvent, isErr bool) {
// if event.ServiceType == constant.DUBBO_CONSUMER ||
// 	event.ServiceType == constant.HTTP_CLIENT_4 ||
// 	event.ServiceType == constant.MYSQL_EXECUTE_QUERY ||
// 	event.ServiceType == constant.REDIS ||
// 	event.ServiceType == constant.ORACLE_EXECUTE_QUERY ||
// 	event.ServiceType == constant.MARIADB_EXECUTE_QUERY {

// 	destinationID := event.GetDestinationId()
// 	if len(destinationID) <= 0 {
// 		return
// 	}
// 	child, ok := s.ServerMap.Childs[event.ServiceType]
// 	if !ok {
// 		child = metric.NewChild()
// 		s.ServerMap.Childs[event.ServiceType] = child
// 	}

// 	// http&&dubbo做特殊处理
// 	if event.ServiceType == constant.HTTP_CLIENT_4 || event.ServiceType == constant.DUBBO_CONSUMER {
// 		ip, err := getip(destinationID)
// 		if err == nil {
// 			appName, ok := misc.AddrStore.Get(ip)
// 			if ok {
// 				destinationID = appName
// 			}
// 		}
// 	}

// 	destination, ok := child.Destinations[destinationID]
// 	if !ok {
// 		destination = metric.NewDestination()
// 		child.Destinations[destinationID] = destination
// 	}

// 	if event.ServiceType == constant.HTTP_CLIENT_4 {
// 		for _, annotation := range event.GetAnnotations() {
// 			if annotation.GetKey() == constant.HTTP_STATUS_CODE {
// 				if _, ok := s.RespCodes[int(annotation.Value.GetIntValue())]; !ok {
// 					destination.AccessErrCount++
// 				}
// 			}
// 		}
// 	}

// 	destination.Count++
// 	if isErr {
// 		destination.ExceptionCount++
// 	}
// 	destination.Duration += event.EndElapsed
// }
// }

// apiCounter 计算api信息
func (s *Stats) apiCounter(span *trace.TSpan) {
	apiStr := span.GetRPC()
	if len(apiStr) <= 0 {
		return
	}
	apiInfo, ok := s.APIStats.Get(apiStr)
	if !ok {
		apiInfo = metric.NewAPIInfo()
		s.APIStats.Store(apiStr, apiInfo)
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

// 计算child拓扑图
func (s *Stats) targetMapCounter(event *trace.TSpanEvent) {
	if event.ServiceType == constant.DUBBO_CONSUMER ||
		event.ServiceType == constant.HTTP_CLIENT_4 ||
		event.ServiceType == constant.MYSQL_EXECUTE_QUERY ||
		event.ServiceType == constant.REDIS ||
		event.ServiceType == constant.ORACLE_EXECUTE_QUERY ||
		event.ServiceType == constant.MARIADB_EXECUTE_QUERY {

		destinationID := event.GetDestinationId()
		if len(destinationID) <= 0 {
			return
		}

		targets, ok := s.SrvMap.Targets[event.ServiceType]
		if !ok {
			targets = make(map[string]*metric.Target)
			s.SrvMap.Targets[event.ServiceType] = targets
		}

		// http&&dubbo做特殊处理
		if event.ServiceType == constant.HTTP_CLIENT_4 || event.ServiceType == constant.DUBBO_CONSUMER {
			ip, err := getip(destinationID)
			if err == nil {
				appName, ok := misc.AddrStore.Get(ip)
				if ok {
					destinationID = appName
				}
			}
		}

		target, ok := targets[destinationID]
		if !ok {
			target = metric.NewTarget()
			targets[destinationID] = target
		}

		if event.ServiceType == constant.HTTP_CLIENT_4 {
			for _, annotation := range event.GetAnnotations() {
				if annotation.GetKey() == constant.HTTP_STATUS_CODE {
					if _, ok := s.RespCodes[int(annotation.Value.GetIntValue())]; !ok {
						target.AccessErrCount++
					}
				}
			}
		}
		target.AccessCount++
		target.AccessDuration += event.EndElapsed
	}
}

// counterEvents 计算method、SQL信息
func (s *Stats) eventsCounter(span *trace.TSpan) {
	if len(s.MethodStats.APIStr) == 0 {
		s.MethodStats.APIStr = span.GetRPC()
	}

	for _, event := range span.GetSpanEventList() {
		isErr := false
		// 是否有异常抛出
		if event.GetExceptionInfo() != nil {
			isErr = true
		}

		// app后续服务拓扑图计算
		s.targetMapCounter(event)

		// 计算method
		s.methodCount(event.GetApiId(), int(event.GetServiceType()), event.EndElapsed, isErr)

		// 计算sql
		annotations := event.GetAnnotations()
		for _, annotation := range annotations {
			// 20为数据库类型
			if annotation.GetKey() == constant.SQL_ID {
				s.sqlCount(annotation.Value.GetIntStringStringValue().GetIntValue(), event.EndElapsed, isErr)
			}
		}
		// 异常计算
		s.exceptionCount(event.GetApiId(), event)
	}
}

func (s *Stats) methodCount(apiID int32, srvType int, elapsed int32, isErr bool) {
	methodInfo, ok := s.MethodStats.Get(apiID)
	if !ok {
		methodInfo = metric.NewMethodInfo()
		methodInfo.Type = srvType
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
		sqlInfo = metric.NewSQLInfo()
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
func (s *Stats) exceptionCount(methodID int32, event *trace.TSpanEvent) {
	// 参看是否存在异常，不存在直接返回
	exInfo := event.GetExceptionInfo()
	if exInfo == nil {
		return
	}

	apiEx, ok := s.ExceptionsStats.Get(methodID)
	if !ok {
		apiEx = metric.NewAPIExceptions()
		s.ExceptionsStats.Store(methodID, apiEx)
	}

	ex, ok := apiEx.Exceptions[exInfo.GetIntValue()]
	if !ok {
		ex = metric.NewExceptionInfo()
		apiEx.Exceptions[exInfo.GetIntValue()] = ex
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
