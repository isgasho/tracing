package service

import (
	"sort"
	"sync"
	"time"

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
	appType      int32                          // 服务类型
	taskID       int64                          // 定时任务ID
	name         string                         // 服务名称
	agents       map[string]*Agent              // agent集合
	stopC        chan bool                      // 停止通道
	tickerC      chan bool                      // 定时任务通道
	spanC        chan *trace.TSpan              // span类型通道
	spanChunkC   chan *trace.TSpanChunk         // span chunk类型通道
	statC        chan *pinpoint.TAgentStat      // jvm状态类型通道
	apis         map[string]struct{}            // 接口信息
	orderlyKey   stats.OrderlyKey               // 排序打点
	points       map[int64]*stats.Stats         // 计算点集合
	srvMapKey    int64                          // 拓扑计算当前计算key
	srvmap       map[int64]*metric.SrvMapStats  // 服务拓扑图
	apiCallKey   int64                          // API被调用计算当前计算key
	apiCall      map[int64]*metric.APICallStats // API被调用
	runtimeStats map[int64]*stats.RuntimeStats  // runtime stats
	runtimeKey   int64                          // runtime key
}

func newApp(name string, appType int32) *App {
	app := &App{
		appType:      appType,
		name:         name,
		agents:       make(map[string]*Agent),
		tickerC:      make(chan bool, 10),
		spanC:        make(chan *trace.TSpan, 200),
		spanChunkC:   make(chan *trace.TSpanChunk, 200),
		statC:        make(chan *pinpoint.TAgentStat, 200),
		apis:         make(map[string]struct{}),
		points:       make(map[int64]*stats.Stats),
		srvmap:       make(map[int64]*metric.SrvMapStats),
		apiCall:      make(map[int64]*metric.APICallStats),
		runtimeStats: make(map[int64]*stats.RuntimeStats),
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
			a.linkTrace()
			// 拓扑信息计算
			a.reportSrvMap()
			// 被调用统计
			a.reportCall()
			// runtime 统计
			a.reportRuntime()
			break
		case span, ok := <-a.spanC:
			if ok {
				if err := a.statsSpan(span); err != nil {
					logger.Warn("stats span", zap.String("error", err.Error()))
				}
			}
			break
		case spanChunk, ok := <-a.spanChunkC:
			if ok {
				if err := a.statsSpanChunk(spanChunk); err != nil {
					logger.Warn("stats span", zap.String("error", err.Error()))
				}
			}
			break
		case agentStat, ok := <-a.statC:
			if ok {
				if err := a.statsAgentStat(agentStat); err != nil {
					logger.Warn("stats agent stat", zap.String("error", err.Error()))
				}
			}
			break
		// case <-a.statBatchC:
		// 	break
		case <-a.stopC:
			return
		}
	}
}

func (a *App) statsAgentStat(agentStat *pinpoint.TAgentStat) error {
	// 计算当前TAgentStat时间范围点
	t, err := utils.MSToTime(agentStat.GetTimestamp())
	if err != nil {
		logger.Warn("ms to time", zap.Int64("time", agentStat.GetTimestamp()), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳并将其精确到秒
	nowSecond := t.Unix()
	if a.runtimeKey == 0 {
		a.runtimeKey = nowSecond + misc.Conf.Stats.RuntimeDefer
	} else {
		// 需要更新计算下标key
		if nowSecond > a.runtimeKey+misc.Conf.Stats.RuntimeDefer {
			a.runtimeKey = nowSecond + misc.Conf.Stats.RuntimeDefer
		}
	}
	runtimeStat, ok := a.runtimeStats[a.runtimeKey]
	if !ok {
		runtimeStat = stats.NewRuntimeStats()
		a.runtimeStats[a.runtimeKey] = runtimeStat
	}

	runtimeStat.JVMCounter(agentStat)

	return nil
}

// stats 计算模块
func (a *App) statsSpan(span *trace.TSpan) error {
	// api缓存并入库
	if !a.apiIsExist(span.GetRPC()) {
		if err := gCollector.storage.StoreAPI(span); err != nil {
			logger.Warn("store api", zap.String("error", err.Error()))
			return err
		}
		a.storeAPI(span.GetRPC())
	}

	// 计算当前span时间范围点
	t, err := utils.MSToTime(span.StartTime)
	if err != nil {
		logger.Warn("ms to time", zap.Int64("time", span.StartTime), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳并将其精确到分钟
	nowSecond := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请
	lstats, ok := a.points[nowSecond]
	if !ok {
		lstats = stats.NewStats()
		a.points[nowSecond] = lstats
	}
	// 计算
	// 计算服务拓扑图，api被调用需要将spanKey加上一个时间范围
	if a.srvMapKey == 0 {
		a.srvMapKey = nowSecond + misc.Conf.Stats.MapDefer
	} else {
		// 需要更新计算下标key
		if nowSecond > a.srvMapKey+misc.Conf.Stats.MapDefer {
			a.srvMapKey = nowSecond + misc.Conf.Stats.MapDefer
		}
	}

	// 计算服务拓扑图，api被调用需要将nowSecond加上一个时间范围
	if a.apiCallKey == 0 {
		a.apiCallKey = nowSecond + misc.Conf.Stats.APICallDefer
	} else {
		// 需要更新计算下标key
		if nowSecond > a.apiCallKey+misc.Conf.Stats.APICallDefer {
			a.apiCallKey = nowSecond + misc.Conf.Stats.APICallDefer
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
		logger.Warn("ms to time", zap.Int64("time", spanChunk.GetKeyTime()), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳
	nowSecond := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请
	lstats, ok := a.points[nowSecond]
	if !ok {
		lstats = stats.NewStats()
		a.points[nowSecond] = lstats
	}

	// 计算
	// 计算服务拓扑图，api被调用需要将nowSecond加上一个时间范围
	if a.srvMapKey == 0 {
		a.srvMapKey = nowSecond + misc.Conf.Stats.MapDefer
	} else {
		// 需要更新计算下标key
		if nowSecond > a.srvMapKey+misc.Conf.Stats.MapDefer {
			a.srvMapKey = nowSecond + misc.Conf.Stats.MapDefer
		}
	}

	// 计算服务拓扑图，api被调用需要将nowSecond加上一个时间范围
	if a.apiCallKey == 0 {
		a.apiCallKey = nowSecond + misc.Conf.Stats.APICallDefer
	} else {
		// 需要更新计算下标key
		if nowSecond > a.apiCallKey+misc.Conf.Stats.APICallDefer {
			a.apiCallKey = nowSecond + misc.Conf.Stats.APICallDefer
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
	a.taskID = gCollector.ticker.NewID()
	logger.Info("app start", zap.String("name", a.name), zap.Int64("taskID", a.taskID))
	// 加入定时模块
	gCollector.ticker.AddTask(a.taskID, a.tickerC)
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

func (a *App) recvAgentStat(appName, agentID string, agentStat *pinpoint.TAgentStat) error {
	a.statC <- agentStat
	return nil
}

// linkTrace 链路接口等计算上报
func (a *App) linkTrace() error {
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
	inputDate := key // + misc.Conf.Stats.DeferTime
	for apiStr, apiInfo := range a.points[key].APIStats.APIs {
		gCollector.storage.InsertAPIStats(a.name, inputDate, apiStr, apiInfo)
	}

	for methodID, methodInfo := range a.points[key].MethodStats.Methods {
		gCollector.storage.InsertMethodStats(a.name, inputDate, a.points[key].MethodStats.APIStr, methodID, methodInfo)
	}

	for sqlID, sqlInfo := range a.points[key].SQLStats.SQLs {
		gCollector.storage.InsertSQLStats(a.name, inputDate, sqlID, sqlInfo)
	}

	for methodID, exceptions := range a.points[key].ExceptionsStats.MethodEx {
		gCollector.storage.InsertExceptionStats(a.name, inputDate, methodID, exceptions.Exceptions)
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
	if time.Now().Unix() < key+misc.Conf.Stats.MapDefer {
		return nil
	}

	inputDate := key // + misc.Conf.Stats.MapDefer

	for parentName, parentInfo := range a.srvmap[key].Parents {
		gCollector.storage.InsertParentMap(a.name, a.appType, inputDate, parentName, parentInfo)
	}

	for childType, child := range a.srvmap[key].Childs {
		for destinationStr, destination := range child.Destinations {
			gCollector.storage.InsertChildMap(a.name, a.appType, inputDate, int32(childType), destinationStr, destination)
		}
	}

	unknowParent := a.srvmap[key].UnknowParent
	// 只有被调用才可以入库
	if unknowParent.Count > 0 {
		gCollector.storage.InsertUnknowParentMap(a.name, a.appType, inputDate, unknowParent)
	}

	// 上报打点信息并删除该时间点信息
	delete(a.srvmap, key)

	return nil
}

// reportCall api被调用情况
func (a *App) reportCall() error {
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
	if time.Now().Unix() < key+misc.Conf.Stats.APICallDefer {
		return nil
	}
	inputDate := key //+ misc.Conf.Stats.APICallDefer
	for apiID, apiInfo := range a.apiCall[key].APIS {
		for parentName, parentInfo := range apiInfo.Parents {
			gCollector.storage.InsertAPICallStats(a.name, a.appType, inputDate, apiID, parentName, parentInfo)
		}
	}

	// 上报打点信息并删除该时间点信息
	delete(a.apiCall, key)

	return nil
}

// reportRuntime runtime计算结果上报&存储
func (a *App) reportRuntime() error {
	// 清空之前节点
	a.orderlyKey = a.orderlyKey[:0]
	// 赋值
	for key := range a.runtimeStats {
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
	if time.Now().Unix() < key+misc.Conf.Stats.RuntimeDefer {
		return nil
	}

	inputDate := key //+ misc.Conf.Stats.APICallDefer

	for agentID, info := range a.runtimeStats[key].JVMStats.Agents {
		gCollector.storage.InsertRuntimeStats(a.name, inputDate, agentID, info)
	}

	// 上报打点信息并删除该时间点信息
	delete(a.runtimeStats, key)

	return nil
}
