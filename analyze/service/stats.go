package service

import (
	"sync"
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
	"go.uber.org/zap"
)

// Stats 离线计算
type Stats struct {
}

// NewStats ...
func NewStats() *Stats {
	return &Stats{}
}

// Start ...
func (s *Stats) Start() error {
	g.L.Info("Start Stats")

	return nil
}

// Close ...
func (s *Stats) Close() error {
	g.L.Info("Close Stats")

	return nil
}

func (s *Stats) counter(app *App, wg *sync.WaitGroup) {
	defer wg.Done()
	// 如果最后一次计算点为0，那么放弃本次计算
	if app.lastCountTime == 0 {
		return
	}

	var queryStartTime int64
	var queryEndTime int64

	queryStartTime = app.lastCountTime
	queryEndTime = app.lastCountTime + misc.Conf.Stats.Range*1000

	// 结束时间要比当前时间少三分钟，这样可以确保数据准确性
	if (queryEndTime + 180*1000) >= time.Now().UnixNano()/1e6 {
		return
	}

	es := GetElements(queryStartTime, queryEndTime)
	queryTraceID := `SELECT trace_id, span_id FROM app_operation_index WHERE app_name=? and start_time>? and start_time<=?;`
	iterTraceID := gAnalyze.appStore.db.Session.Query(queryTraceID, app.AppName, queryStartTime, queryEndTime).Iter()
	defer iterTraceID.Close()

	var traceID string
	var spanID int64
	for iterTraceID.Scan(&traceID, &spanID) {
		spanCounter(traceID, spanID, es)
	}

	statsCounter(app, queryStartTime, queryEndTime, es)

	// @TODO记录计算结果
	for recordTime, e := range es {
		spanCounterRecord(app, recordTime, e)
	}

	// @TODO
	// 记录计算时间到表
	if err := gAnalyze.db.Session.Query(gUpdateLastCounterTime, queryEndTime, app.AppName).Exec(); err != nil {
		g.L.Warn("update Last Counter Time error", zap.String("error", err.Error()), zap.String("SQL", gUpdateLastCounterTime))
		return
	}
}

// Counter ...
func (s *Stats) Counter() error {
	var wg sync.WaitGroup
	// 这里appStore没有用锁的原因是因为Counter和loadApp函数是串联调用的
	// gAnalyze.appStore.RLock()
	for _, app := range gAnalyze.appStore.Apps {
		wg.Add(1)
		// 每个应用一个携程去计算
		// 只有等所有应用计算完毕才会进行下一轮计算
		go gAnalyze.stats.counter(app, &wg)
	}
	// gAnalyze.appStore.RUnlock()
	wg.Wait()
	return nil
}
