package service

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"go.uber.org/zap"

	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/collector/stats"
	"github.com/imdevlab/tracing/pkg/metric"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
)

// 服务统计数据只实时计算1分钟的点，不做任何滑动窗口
// 通过nats或者其他mq将1分钟的数据发送给聚合计算服务，在聚合服务上做告警策略

// App 单个服务信息
type App struct {
	sync.RWMutex
	appType    int32                          // 服务类型
	taskID     int64                          // 定时任务ID
	name       string                         // 服务名称
	agents     map[string]*Agent              // agent集合
	stopC      chan bool                      // 停止通道
	tickerC    chan bool                      // 定时任务通道
	spanC      chan *trace.TSpan              // span类型通道
	spanChunkC chan *trace.TSpanChunk         // span chunk类型通道
	statC      chan *pinpoint.TAgentStat      // jvm状态类型通道
	statBatchC chan *pinpoint.TAgentStatBatch // 批量jvm状态类型通道
	apis       map[string]struct{}            // 接口信息
	orderlyKey stats.OrderlyKey               // 排序打点
	points     map[int64]*stats.Stats         // 计算点集合
	srvMapKey  int64                          // 拓扑计算当前计算key
	srvmap     map[int64]*metric.SrvMapStats  // 服务拓扑图
	apiCallKey int64                          // API被调用计算当前计算key
	apiCall    map[int64]*metric.APICallStats // API被调用
}

func newApp(name string, appType int32) *App {
	app := &App{
		appType:    appType,
		name:       name,
		agents:     make(map[string]*Agent),
		tickerC:    make(chan bool, 10),
		spanC:      make(chan *trace.TSpan, 200),
		spanChunkC: make(chan *trace.TSpanChunk, 200),
		statC:      make(chan *pinpoint.TAgentStat, 200),
		statBatchC: make(chan *pinpoint.TAgentStatBatch, 200),
		apis:       make(map[string]struct{}),
		points:     make(map[int64]*stats.Stats),
		srvmap:     make(map[int64]*metric.SrvMapStats),
		apiCall:    make(map[int64]*metric.APICallStats),
	}
	app.start()
	return app
}

// isExist agent是否存在
func (a *App) isExist(agentid string) bool {
	a.RLock()
	_, ok := a.agents[agentid]
	a.RUnlock()
	if !ok {
		return false
	}
	return true
}

// storeAgent 保存agent
func (a *App) storeAgent(agentid string, startTime int64) {
	a.RLock()
	agent, ok := a.agents[agentid]
	a.RUnlock()
	if ok {
		return
	}

	agent = newAgent(agentid, startTime)
	a.Lock()
	a.agents[agentid] = agent
	a.Unlock()

	return
}

// stats 计算模块
func (a *App) stats() {
	for {
		select {
		case <-a.tickerC:
			// 链路统计信息计算
			a.tickerTrace()
			// 拓扑信息计算
			a.reportSrvMap()
			// 被调用统计
			a.reportCall()
			break
		case span, ok := <-a.spanC:
			if ok {
				if err := a.statsSpan(span); err != nil {
					g.L.Warn("stats span", zap.String("error", err.Error()))
				}
			}
			break
		case spanChunk, ok := <-a.spanChunkC:
			if ok {
				if err := a.statsSpanChunk(spanChunk); err != nil {
					g.L.Warn("stats span", zap.String("error", err.Error()))
				}
			}
			break
		case <-a.statC:
			break
		case <-a.statBatchC:
			break
		case <-a.stopC:
			return
		}
	}
}

// stats 计算模块
func (a *App) statsSpan(span *trace.TSpan) error {
	// api缓存并入库
	if !a.apiIsExist(span.GetRPC()) {
		if err := gCollector.storage.StoreAPI(span); err != nil {
			g.L.Warn("store api", zap.String("error", err.Error()))
			return err
		}
		a.storeAPI(span.GetRPC())
	}

	// 计算当前span时间范围点
	t, err := utils.MSToTime(span.StartTime)
	if err != nil {
		g.L.Warn("ms to time", zap.Int64("time", span.StartTime), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳并将其精确到分钟
	spanKey := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请
	lstats, ok := a.points[spanKey]
	if !ok {
		lstats = stats.NewStats()
		a.points[spanKey] = lstats
	}
	// 计算
	// 计算服务拓扑图，api被调用需要将spanKey加上一个时间范围
	if a.srvMapKey == 0 {
		a.srvMapKey = spanKey + misc.Conf.Stats.MapRange
	} else {
		// 需要更新计算下标key
		if a.srvMapKey+misc.Conf.Stats.MapRange > spanKey {
			a.srvMapKey = spanKey + misc.Conf.Stats.MapRange
		}
	}

	// 计算服务拓扑图，api被调用需要将spanKey加上一个时间范围
	if a.apiCallKey == 0 {
		a.apiCallKey = spanKey + misc.Conf.Stats.APICallRang
	} else {
		// 需要更新计算下标key
		if a.apiCallKey+misc.Conf.Stats.APICallRang > spanKey {
			a.apiCallKey = spanKey + misc.Conf.Stats.APICallRang
		}
	}

	// 获取拓扑计算点
	srvMap, ok := a.srvmap[a.srvMapKey]
	if !ok {
		// 新点保存
		srvMap = metric.NewSrvMapStats()
		a.srvmap[a.srvMapKey] = srvMap
	}

	// 获取Apicall计算节点
	apiCall, ok := a.apiCall[a.apiCallKey]
	if !ok {
		apiCall = metric.NewAPICallStats()
		a.apiCall[a.apiCallKey] = apiCall
	}

	lstats.SpanCounter(span, srvMap, apiCall)

	return nil
}

// statsSpanChunk 计算模块
func (a *App) statsSpanChunk(spanChunk *trace.TSpanChunk) error {

	// 计算当前spanChunk时间范围点
	t, err := utils.MSToTime(spanChunk.GetKeyTime())
	if err != nil {
		g.L.Warn("ms to time", zap.Int64("time", spanChunk.GetKeyTime()), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳
	spanKey := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请
	lstats, ok := a.points[spanKey]
	if !ok {
		lstats = stats.NewStats()
		a.points[spanKey] = lstats
	}

	// 计算
	// 计算服务拓扑图，api被调用需要将spanKey加上一个时间范围
	if a.srvMapKey == 0 {
		a.srvMapKey = spanKey + misc.Conf.Stats.MapRange
	} else {
		// 需要更新计算下标key
		if a.srvMapKey+misc.Conf.Stats.MapRange > spanKey {
			a.srvMapKey = spanKey + misc.Conf.Stats.MapRange
		}
	}

	// 计算服务拓扑图，api被调用需要将spanKey加上一个时间范围
	if a.apiCallKey == 0 {
		a.apiCallKey = spanKey + misc.Conf.Stats.APICallRang
	} else {
		// 需要更新计算下标key
		if a.apiCallKey+misc.Conf.Stats.APICallRang > spanKey {
			a.apiCallKey = spanKey + misc.Conf.Stats.APICallRang
		}
	}

	// 获取拓扑计算点
	srvMap, ok := a.srvmap[a.srvMapKey]
	if !ok {
		// 新点保存
		srvMap = metric.NewSrvMapStats()
		a.srvmap[a.srvMapKey] = srvMap
	}

	// 获取Apicall计算节点
	apiCall, ok := a.apiCall[a.apiCallKey]
	if !ok {
		apiCall = metric.NewAPICallStats()
		a.apiCall[a.apiCallKey] = apiCall
	}

	lstats.SpanChunkCounter(spanChunk, srvMap, apiCall)

	return nil
}

func (a *App) start() {
	// 获取任务ID
	a.taskID = gCollector.tickers.NewID()
	g.L.Info("app start", zap.String("name", a.name), zap.Int64("taskID", a.taskID))
	// 加入定时模块
	gCollector.tickers.AddTask(a.taskID, a.tickerC)
	// 启动计算模块
	go a.stats()
}

// apiIsExist 检查api是否缓存
func (a *App) apiIsExist(api string) bool {
	a.RLock()
	_, isExist := a.apis[api]
	a.RUnlock()
	return isExist
}

// storeAPI 缓存api
func (a *App) storeAPI(api string) {
	a.Lock()
	a.apis[api] = struct{}{}
	a.Unlock()
}

func (a *App) recvSpan(appName, agentID string, span *trace.TSpan) error {
	a.spanC <- span
	return nil
}

func (a *App) recvSpanChunk(appName, agentID string, spanChunk *trace.TSpanChunk) error {
	a.spanChunkC <- spanChunk
	return nil
}

// tickerTrace 链路接口等计算上报
func (a *App) tickerTrace() error {
	// 清空之前节点
	a.orderlyKey = a.orderlyKey[:0]

	// 赋值
	for key := range a.points {
		a.orderlyKey = append(a.orderlyKey, key)
	}

	// 排序
	sort.Sort(a.orderlyKey)

	// 如果没有计算节点直接返回
	if a.orderlyKey.Len() <= 0 {
		return nil
	}

	key := a.orderlyKey[0]

	// 延迟计算，防止defer 时间内span未上报
	if time.Now().Unix() < key+misc.Conf.Stats.DeferTime {
		return nil
	}

	for apiStr, apiInfo := range a.points[key].APIStats.APIs {
		gCollector.storage.InsertAPIStats(a.name, key, apiStr, apiInfo)
	}

	for methodID, methodInfo := range a.points[key].MethodStats.Methods {
		gCollector.storage.InsertMethodStats(a.name, key, a.points[key].MethodStats.APIStr, methodID, methodInfo)
	}

	for sqlID, sqlInfo := range a.points[key].SQLStats.SQLs {
		gCollector.storage.InsertSQLStats(a.name, key, sqlID, sqlInfo)
	}

	for methodID, exceptions := range a.points[key].ExceptionsStats.MethodEx {
		gCollector.storage.InsertExceptionStats(a.name, key, methodID, exceptions.Exceptions)
	}

	// 上报打点信息并删除该时间点信息
	delete(a.points, key)
	return nil
}

// 各类拓扑图定时计算上报
func (a *App) reportSrvMap() error {

	// 清空之前节点
	a.orderlyKey = a.orderlyKey[:0]
	// 赋值
	for key := range a.srvmap {
		a.orderlyKey = append(a.orderlyKey, key)
	}
	// 排序
	sort.Sort(a.orderlyKey)
	// 如果没有计算节点直接返回
	if a.orderlyKey.Len() <= 0 {
		return nil
	}

	key := a.orderlyKey[0]
	// 延迟计算，防止defer 时间内span未上报
	if time.Now().Unix() < key+misc.Conf.Stats.MapRange {
		return nil
	}

	for parentName, parentInfo := range a.srvmap[key].SrvMaps {
		gCollector.storage.InsertServiceMap(a.name, a.appType, key, parentName, parentInfo)
	}

	for dbType, dbInfo := range a.srvmap[key].DBMaps {
		gCollector.storage.InsertDBMap(a.name, a.appType, key, int32(dbType), dbInfo)
	}

	// 上报打点信息并删除该时间点信息
	delete(a.srvmap, key)

	return nil
}

// reportCall api被调用情况
func (a *App) reportCall() error {
	// log.Println("apiCall", a.apiCall)
	// log.Println("apiCall", a.apiCall)
	// 清空之前节点
	a.orderlyKey = a.orderlyKey[:0]
	// 赋值
	for key := range a.apiCall {
		a.orderlyKey = append(a.orderlyKey, key)
	}
	// 排序
	sort.Sort(a.orderlyKey)
	// 如果没有计算节点直接返回
	if a.orderlyKey.Len() <= 0 {
		return nil
	}

	key := a.orderlyKey[0]
	// 延迟计算，防止defer 时间内span未上报
	if time.Now().Unix() < key+misc.Conf.Stats.APICallRang {
		return nil
	}

	for apiID, callInfo := range a.apiCall[key].APIS {
		for parentName, parentInfo := range callInfo.Parents {
			log.Println("服务被调用统计", apiID, parentName, parentInfo)
		}
	}

	// 上报打点信息并删除该时间点信息
	delete(a.apiCall, key)

	return nil
}
