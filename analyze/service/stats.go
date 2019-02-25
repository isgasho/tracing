package service

import (
	"sync"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/vgo/analyze/misc"
	"go.uber.org/zap"
)

// Stats ...
type Stats interface {
	Start() error
	Close() error
	Counter() error
}

// stats 离线计算
type stats struct {
}

// newStats ...
func newStats() *stats {
	return &stats{}
}

// Start ...
func (s *stats) Start() error {
	g.L.Info("Start stats")

	return nil
}

// Close ...
func (s *stats) Close() error {
	g.L.Info("Close stats")

	return nil
}

func (s *stats) counter(app *App, wg *sync.WaitGroup) {
	defer wg.Done()
	// 如果最后一次计算点为0，那么放弃本次计算
	// if app.lastCountTime == 0 || len(app.Agents) == 0 {
	if app.lastCountTime == 0 {
		// log.Println("如果最后一次计算点为0，那么放弃本次计算")
		return
	}

	var queryStartTime int64
	var queryEndTime int64

	queryStartTime = app.lastCountTime
	queryEndTime = app.lastCountTime + misc.Conf.Stats.Range*1000

	// 结束时间要比当前时间少三分钟，这样可以确保数据准确性
	if (queryEndTime + 3*60*1000) >= time.Now().UnixNano()/1e6 {
		// log.Println("上次计算时间间隔太短,等待")
		return
	}

	//	计算出每个分钟点，并生成map
	es := GetElements(queryStartTime, queryEndTime)

	// 根据时间范围查出所有符合范围的traceID
	iterTraceID := gAnalyze.appStore.cql.Session.Query(misc.QueryTraceID, app.AppName, queryStartTime, queryEndTime).Iter()

	defer func() {
		if err := iterTraceID.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}()

	var traceID string
	var spanID int64

	// 根据traceID 查出span并进行计算
	for iterTraceID.Scan(&traceID, &spanID) {
		spanCounter(traceID, spanID, es)
	}

	// 统计jvm信息
	statsCounter(app, queryStartTime, queryEndTime, es)

	// 记录计算结果
	for recordTime, e := range es {
		spanCounterRecord(app, recordTime, e)
	}

	// 将本次计算时间记录到表中
	query := gAnalyze.cql.Session.Query(misc.UpdateLastCounterTime, queryEndTime, app.AppName)
	if err := query.Exec(); err != nil {
		g.L.Warn("update Last Counter Time error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
		return
	}
	// 缓存到内存
	app.lastCountTime = queryEndTime
}

// Counter ...
func (s *stats) Counter() error {
	var wg sync.WaitGroup
	// 这里appStore没有用锁的原因是因为Counter和loadApp函数是串联调用的
	// gAnalyze.appStore.RLock()
	for _, app := range gAnalyze.appStore.Apps {
		wg.Add(1)
		// 每个应用一个携程去计算
		// 只有等所有应用计算完毕才会进行下一轮计算
		go s.counter(app, &wg)
	}
	// gAnalyze.appStore.RUnlock()
	wg.Wait()
	return nil
}
