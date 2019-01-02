package service

import (
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

	// queryAgentStat := `SELECT timestamp, stat_info  FROM agent_stats WHERE app_name=? AND  agent_id=? and timestamp>? and timestamp<=?;`
	// for _, agent := range app.Agents {
	// 	iterAgentStat := gAnalyze.appStore.db.Session.Query(queryAgentStat, app.AppName, agent.AgentID, queryStartTime, queryEndTime).Iter()
	// 	var timestamp int64
	// 	var statInfo []byte
	// 	for iterAgentStat.Scan(&timestamp, &statInfo) {
	// 		tStruct := thrift.Deserialize(statInfo)
	// 		switch m := tStruct.(type) {
	// 		case *pinpoint.TAgentStat:
	// 			log.Println("*pinpoint.TAgentStat", m)

	// 			break
	// 		case *pinpoint.TAgentStatBatch:
	// 			log.Println("*pinpoint.TAgentStatBatch", m)
	// 			for _, stat := range m.AgentStats {
	// 				log.Println("stat.Gc ---->>>", stat.Gc)
	// 			}
	// 			break
	// 		default:
	// 			g.L.Warn("unknow type", zap.String("type", fmt.Sprintf("%T", m)))
	// 		}
	// 	}
	// 	iterAgentStat.Close()
	// }

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
