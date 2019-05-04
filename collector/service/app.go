package service

import (
	"sort"
	"sync"
	"time"

	"github.com/imdevlab/g/utils"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"

	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/pkg/alert"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/metric"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/util"
)

// 服务统计数据只实时计算1分钟的点，不做任何滑动窗口
// 通过nats或者其他mq将1分钟的数据发送给聚合计算服务，在聚合服务上做告警策略

// App 单个服务信息
type App struct {
	sync.RWMutex
	appType    int32                     // 服务类型
	taskID     int64                     // 定时任务ID
	name       string                    // 服务名称
	orderKeys  OrderlyKeys               // 排序打点
	agents     map[string]*util.Agent    // agent集合
	apis       map[string]struct{}       // 接口信息
	points     map[int64]*Stats          // 计算点集合
	apiMapKey  int64                     // API被调用计算当前计算key
	apiMap     map[int64]*metric.APIMap  // API被调用
	stopC      chan bool                 // 停止通道
	tickerC    chan bool                 // 定时任务通道
	spanC      chan *trace.TSpan         // span类型通道
	spanChunkC chan *trace.TSpanChunk    // span chunk类型通道
	statC      chan *pinpoint.TAgentStat // jvm状态类型通道

	APIStats        *metric.APIStats        // api计算统计
	MethodStats     *metric.MethodStats     // 接口计算统计
	SQLStats        *metric.SQLStats        // sql语句计算统计
	ExceptionsStats *metric.ExceptionsStats // 异常计算统计
	respCodes       map[int]struct{}        //
	SrvMap          *metric.SrvMap          // 服务拓扑图
}

func newApp(name string) *App {
	app := &App{
		// appType:    appType,
		name:       name,
		agents:     make(map[string]*util.Agent),
		stopC:      make(chan bool, 1),
		tickerC:    make(chan bool, 10),
		spanC:      make(chan *trace.TSpan, 200),
		spanChunkC: make(chan *trace.TSpanChunk, 200),
		statC:      make(chan *pinpoint.TAgentStat, 200),
		apis:       make(map[string]struct{}),
		points:     make(map[int64]*Stats),
		apiMap:     make(map[int64]*metric.APIMap),
		respCodes:  make(map[int]struct{}),

		APIStats:        metric.NewAPIStats(),
		MethodStats:     metric.NewMethodStats(),
		SQLStats:        metric.NewSQLStats(),
		ExceptionsStats: metric.NewExceptionsStats(),
		SrvMap:          metric.NewSrvMap(),
	}
	// @TODO codes会从策略模版中去取
	// 默认200
	app.respCodes[200] = struct{}{}
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
func (a *App) storeAgent(agentID string, isLive bool) {
	a.RLock()
	agent, ok := a.agents[agentID]
	a.RUnlock()
	if !ok {
		agent = util.NewAgent()
		a.Lock()
		a.agents[agentID] = agent
		a.Unlock()
	}
	agent.IsLive = isLive
	return
}

// stats 计算模块
func (a *App) stats() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("app stats", zap.Any("msg", err), zap.String("name", a.name))
			return
		}
	}()

	for {
		select {
		case _, ok := <-a.tickerC:
			if ok {
				// 链路统计信息计算
				a.linkTrace()
				// 被调用统计
				a.reportAPIMap()
			}
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
		case <-a.statC:
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
	spanTime := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请, span统计的范围是分钟，所以这里直接用优化过后的spanTime
	stats, ok := a.points[spanTime]
	if !ok {
		stats = NewStats(a.respCodes)
		a.points[spanTime] = stats
	}

	// api被调用需要将nowSecond加上一个时间范围
	if a.apiMapKey == 0 {
		a.apiMapKey = spanTime + misc.Conf.Stats.APICallRange
	} else {
		// 需要更新计算下标key
		if spanTime > a.apiMapKey {
			a.apiMapKey = spanTime + misc.Conf.Stats.APICallRange
		}
	}
	// 获取Apicall计算节点
	apiMap, ok := a.apiMap[a.apiMapKey]
	if !ok {
		apiMap = metric.NewAPIMap()
		a.apiMap[a.apiMapKey] = apiMap
	}

	stats.SpanCounter(span, apiMap)

	return nil
}

// statsSpanChunk 计算模块
func (a *App) statsSpanChunk(spanChunk *trace.TSpanChunk) error {
	// 计算当前spanChunk时间范围点
	// SpanChunk 里面没有start信息，只能用当前时间来做
	t, err := utils.MSToTime(time.Now().Unix())
	if err != nil {
		logger.Warn("ms to time", zap.Int64("time", time.Now().Unix()), zap.String("error", err.Error()))
		return err
	}

	// 获取时间戳
	spanChunkTime := t.Unix() - int64(t.Second())

	// 查找时间点，不存在新申请
	stats, ok := a.points[spanChunkTime]
	if !ok {
		stats = NewStats(a.respCodes)
		a.points[spanChunkTime] = stats
	}

	stats.SpanChunkCounter(spanChunk)

	return nil
}

func (a *App) start() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("app start", zap.Any("msg", err))
			return
		}
	}()

	// 获取任务ID
	a.taskID = gCollector.ticker.NewID()
	logger.Info("app start", zap.String("name", a.name), zap.Int64("taskID", a.taskID))
	// 加入定时模块
	gCollector.ticker.AddTask(a.taskID, a.tickerC)
	// 启动计算模块
	go a.stats()
}

// close 退出
func (a *App) close() {
	// 获取任务ID
	logger.Info("app close", zap.String("name", a.name), zap.Int64("taskID", a.taskID))

	gCollector.ticker.RemoveTask(a.taskID)
	close(a.tickerC)
	close(a.stopC)
	close(a.spanC)
	close(a.spanChunkC)
	close(a.statC)
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
	a.orderKeys = a.orderKeys[:0]

	// 赋值
	for key := range a.points {
		a.orderKeys = append(a.orderKeys, key)
	}

	// 排序
	sort.Sort(a.orderKeys)

	// 如果没有计算节点直接返回
	if a.orderKeys.Len() <= 0 {
		return nil
	}

	inputDate := a.orderKeys[0]
	now := time.Now().Unix()

	if now < inputDate+misc.Conf.Stats.DeferTime {
		return nil
	}

	// 防止这个点就没有数据，导致程序崩溃
	if a.points[inputDate].SrvMap == nil {
		logger.Warn("srv map is nil", zap.Int64("time", inputDate))
		return nil
	}

	apis := alert.NewAPIs()
	for apiStr, apiInfo := range a.points[inputDate].APIStats.APIs {
		gCollector.storage.InsertAPIStats(a.name, inputDate, apiStr, apiInfo)

		api := &alert.API{
			Desc:     apiStr,
			Count:    apiInfo.Count,
			Errcount: apiInfo.ErrCount,
			Duration: apiInfo.Duration,
		}
		apis.APIS[apiStr] = api
	}
	// 有api数据发送给mq
	if len(apis.APIS) > 0 {
		data := alert.NewData()
		data.AppName = a.name
		data.Type = constant.ALERT_TYPE_API
		data.InputDate = inputDate
		payload, err := msgpack.Marshal(apis)
		if err != nil {
			logger.Warn("msgpack", zap.String("error", err.Error()))
		} else {
			data.Payload = payload
			// 推送
			gCollector.publish(data)
		}
	}

	for methodID, methodInfo := range a.points[inputDate].MethodStats.Methods {
		gCollector.storage.InsertMethodStats(a.name, inputDate, a.points[inputDate].MethodStats.APIStr, methodID, methodInfo)
	}

	sqls := alert.NewSQLs()
	for sqlID, sqlInfo := range a.points[inputDate].SQLStats.SQLs {
		gCollector.storage.InsertSQLStats(a.name, inputDate, sqlID, sqlInfo)
		sql := &alert.SQL{
			ID:       sqlID,
			Count:    sqlInfo.Count,
			Errcount: sqlInfo.ErrCount,
			Duration: sqlInfo.Duration,
		}
		sqls.SQLs[sqlID] = sql
	}

	// 有sql数据发送给mq
	if len(sqls.SQLs) > 0 {
		data := alert.NewData()
		data.AppName = a.name
		data.InputDate = inputDate
		data.Type = constant.ALERT_TYPE_SQL
		payload, err := msgpack.Marshal(sqls)
		if err != nil {
			logger.Warn("msgpack", zap.String("error", err.Error()))
		} else {
			data.Payload = payload
			// 推送
			gCollector.publish(data)
		}
	}

	for methodID, exceptions := range a.points[inputDate].ExceptionsStats.MethodEx {
		gCollector.storage.InsertExceptionStats(a.name, inputDate, methodID, exceptions.Exceptions)
	}

	for parentName, parent := range a.points[inputDate].SrvMap.Parents {
		gCollector.storage.InsertParentMap(a.name, a.appType, inputDate, parentName, parent)
	}

	for childType, targets := range a.points[inputDate].SrvMap.Targets {
		for targetName, target := range targets {
			gCollector.storage.InsertTargetMap(a.name, a.appType, inputDate, int32(childType), targetName, target)
		}
	}

	unknowParent := a.points[inputDate].SrvMap.UnknowParent
	// 只有被调用才可以入库
	if unknowParent.AccessCount > 0 {
		gCollector.storage.InsertUnknowParentMap(a.name, a.appType, inputDate, unknowParent)
	}

	// 上报打点信息并删除该时间点信息
	delete(a.points, inputDate)
	return nil
}

// reportAPIMap api被调用情况
func (a *App) reportAPIMap() error {
	// 清空之前节点
	a.orderKeys = a.orderKeys[:0]
	// 赋值
	for key := range a.apiMap {
		a.orderKeys = append(a.orderKeys, key)
	}
	// 排序
	sort.Sort(a.orderKeys)
	// 如果没有计算节点直接返回
	if a.orderKeys.Len() <= 0 {
		return nil
	}

	inputDate := a.orderKeys[0]
	now := time.Now().Unix()
	if now < inputDate+misc.Conf.Stats.DeferTime {
		return nil
	}

	for apiStr, apiInfo := range a.apiMap[inputDate].APIS {
		for parentName, parentInfo := range apiInfo.Parents {
			gCollector.storage.InsertAPIMapStats(a.name, a.appType, inputDate, apiStr, parentName, parentInfo)
		}
	}

	// 上报打点信息并删除该时间点信息
	delete(a.apiMap, inputDate)

	return nil
}
