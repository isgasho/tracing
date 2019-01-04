package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mafanr/g"
	"go.uber.org/zap"

	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
)

var gCounterQuerySpan string = `SELECT agent_start_time, start_time, rpc, elapsed,  service_type, parent_app_name,
	parent_app_type, span_event_list, err, agent_id
	FROM traces WHERE trace_id=? AND span_id=?;`

var gChunkEventsIterTrace string = `SELECT span_event_list FROM traces_chunk WHERE trace_id=? AND  span_id=?;`

// spanCounter ...
func spanCounter(traceID string, spanID int64, es map[int64]*Element) error {

	iterTrace := gAnalyze.appStore.db.Session.Query(gCounterQuerySpan, traceID, spanID).Iter()

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

	{
		var spanChunkEventList []byte
		iterChunkEvents := gAnalyze.appStore.db.Session.Query(gChunkEventsIterTrace, traceID, spanID).Iter()

		iterChunkEvents.Scan(&spanChunkEventList)
		iterChunkEvents.Close()

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

		if e, ok := es[index]; ok {
			e.urls.urlCounter(rpc, elapsed, isErr)
			e.events.eventsCounter(spanEvents, chunkEvents)
			e.sqls.sqlCounter(spanEvents, chunkEvents)
			// e.exceptions.exceptionCounter(urlStr string, elapsed int, isError int)
		}
	}

	// for _, v := range es {
	// 	if v.urls != nil {
	// 		for uk, uv := range v.urls.urls {
	// 			log.Println(uk, uv)
	// 		}
	// 	}
	// 	// log.Println(v.urls)
	// }

	iterTrace.Close()

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
