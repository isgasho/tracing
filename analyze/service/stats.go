package service

import (
	"sync"
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
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

	es := GetElements(queryStartTime, queryEndTime)
	queryTraceID := `SELECT trace_id, span_id FROM app_operation_index WHERE app_name=? and input_date>? and input_date<=?;`
	iterTraceID := gAnalyze.appStore.cql.Session.Query(queryTraceID, app.AppName, queryStartTime, queryEndTime).Iter()

	defer func() {
		if err := iterTraceID.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}()
	// SELECT trace_id, span_id FROM app_operation_index WHERE app_name='helm' and input_date>1548514140000 and input_date<=1548514200000;
	var traceID string
	var spanID int64

	// log.Println("查询")
	for iterTraceID.Scan(&traceID, &spanID) {
		// log.Println("查询到T让测ID", app.AppName, traceID, spanID)
		spanCounter(traceID, spanID, es)
	}
	// log.Println("查询 2 ", app.AppName, queryStartTime, queryEndTime)
	statsCounter(app, queryStartTime, queryEndTime, es)

	// @TODO记录计算结果
	for recordTime, e := range es {
		spanCounterRecord(app, recordTime, e)
	}

	// @TODO
	// 记录计算时间到表
	if err := gAnalyze.cql.Session.Query(gUpdateLastCounterTime, queryEndTime, app.AppName).Exec(); err != nil {
		g.L.Warn("update Last Counter Time error", zap.String("error", err.Error()), zap.String("SQL", gUpdateLastCounterTime))
		return
	}
	app.lastCountTime = queryEndTime
	// log.Println("插入时间")
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
