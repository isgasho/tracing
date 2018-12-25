package service

import (
	"log"
	"sync"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
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

	queryTraceID := `SELECT trace_id FROM app_operation_index WHERE app_name=? and start_time>=? and start_time<=?;`
	iterTraceID := gAnalyze.appStore.db.Session.Query(queryTraceID, app.AppName, queryStartTime, queryEndTime).Iter()
	var traceID string
	for iterTraceID.Scan(&traceID) {
		log.Println("------------------>>>>", traceID)
	}
	iterTraceID.Close()
}

// func HourTimestamp() int64 {
// 	now := time.Now()
// 	timestamp := now.Unix() - int64(now.Second()) - int64((60 * now.Minute()))
// 	fmt.Println(timestamp, time.Unix(timestamp, 0), now.Unix())
// 	return timestamp
// }

// type Counter struct {
// }

// type

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
