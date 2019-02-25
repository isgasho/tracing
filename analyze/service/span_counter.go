package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/imdevlab/g"
	"go.uber.org/zap"

	"github.com/imdevlab/vgo/analyze/misc"
	"github.com/imdevlab/vgo/proto/pinpoint/thrift/trace"
)

// spanCounter ...
func spanCounter(traceID string, spanID int64, es map[int64]*Element) error {

	iterTrace := gAnalyze.appStore.cql.Session.Query(misc.CounterQuerySpan, traceID, spanID).Iter()
	var startTime int64
	var rpc string
	var elapsed int
	var serviceType int
	var parentAppName string
	var parentAppType int
	var spanEventList []byte
	var isErr int
	var agentID string
	var appName string

	var chunkEvents []*trace.TSpanEvent

	// 查询分片的span信息
	{
		var spanChunkEventList []byte
		iterChunkEvents := gAnalyze.appStore.cql.Session.Query(misc.ChunkEventsIterTrace, traceID, spanID).Iter()

		iterChunkEvents.Scan(&spanChunkEventList)

		if err := iterChunkEvents.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}

		if len(spanChunkEventList) == 0 {
			goto DoSpan
		}
		err := json.Unmarshal(spanChunkEventList, &chunkEvents)
		if err != nil {
			g.L.Warn("json.Unmarshal error", zap.String("error", err.Error()))
		}

	}

DoSpan:
	for iterTrace.Scan(&appName, &startTime, &rpc, &elapsed, &serviceType, &parentAppName, &parentAppType, &spanEventList, &isErr, &agentID) {
		index, _ := ModMs2Min(startTime)
		var spanEvents []*trace.TSpanEvent
		err := json.Unmarshal(spanEventList, &spanEvents)
		if err != nil {
			g.L.Warn("json.Unmarshal error", zap.String("error", err.Error()))
			continue
		}

		// 查询缓存并记录API信息到数据
		if app, ok := gAnalyze.appStore.getApp(appName); ok {
			if _, ok := app.getAPI(rpc); !ok {
				query := gAnalyze.cql.Session.Query(misc.InsertAPIs, appName, rpc)
				if err := query.Exec(); err != nil {
					g.L.Warn("json.Unmarshal error", zap.String("error", err.Error()), zap.String("query", query.String()))
				}
				app.storeAPI(rpc)
			}
		}

		// 对index时间点到数据进行计算
		if e, ok := es[index]; ok {
			// API
			e.apis.apiCounter(rpc, elapsed, isErr)
			// 内部事件method_id相关计算
			e.events.eventsCounter(rpc, spanEvents, chunkEvents)
			// SQL统计
			e.sqls.sqlCounter(spanEvents, chunkEvents)
			// 异常统计
			e.exceptions.exceptionCounter(spanEvents, chunkEvents)
		}
	}

	if err := iterTrace.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}
	return nil
}

// ModMs2Min 取整
func ModMs2Min(ms int64) (int64, error) {
	if ms == 0 {
		return 0, fmt.Errorf("ms is 0")
	}

	nsec := ms * 1e6
	t := time.Unix(0, nsec)

	return t.Unix() - int64(t.Second()), nil
}

func spanCounterRecord(app *App, inputDate int64, e *Element) error {
	e.apis.apiRecord(app, inputDate)
	e.events.eventRecord(app, inputDate)
	e.sqls.sqlRecord(app, inputDate)
	e.stats.statRecord(app, inputDate)
	e.exceptions.exceptionRecord(app, inputDate)
	return nil
}
