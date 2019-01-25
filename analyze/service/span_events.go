package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"go.uber.org/zap"
)

// SpanEvents ...
type SpanEvents struct {
	rpc        string
	spanEvents map[int32]*SpanEvent
}

var gCounterQueryAPI string = `SELECT api_info FROM agent_apis WHERE api_id=? 
		and agent_id=? and app_name=? and input_date=?;`

func (spanEvents *SpanEvents) eventsCounter(rpc string, events []*trace.TSpanEvent, chunkEvents []*trace.TSpanEvent) error {
	if len(spanEvents.rpc) == 0 {
		spanEvents.rpc = rpc
	}
	for _, event := range events {
		api, ok := spanEvents.spanEvents[event.GetApiId()]
		if !ok {
			api = NewSpanEvent()
			spanEvents.spanEvents[event.GetApiId()] = api
		}
		api.count++
		elapsed := int(event.EndElapsed)
		api.elapsed += elapsed
		api.serType = int(event.ServiceType)
		if elapsed > api.maxElapsed {
			api.maxElapsed = api.elapsed
		}

		if api.minElapsed == 0 || api.minElapsed > elapsed {
			api.minElapsed = elapsed
		}

		// 是否有异常抛出
		if event.GetExceptionInfo() != nil {
			api.errCount++
		}

		api.averageElapsed = float64(api.elapsed) / float64(api.count)
	}

	for _, event := range chunkEvents {
		api, ok := spanEvents.spanEvents[event.GetApiId()]
		if !ok {
			api = NewSpanEvent()
			spanEvents.spanEvents[event.GetApiId()] = api
		}
		api.count++
		elapsed := int(event.EndElapsed - event.StartElapsed)
		api.elapsed += elapsed
		api.serType = int(event.ServiceType)
		if elapsed > api.maxElapsed {
			api.maxElapsed = api.elapsed
		}

		if api.minElapsed == 0 || api.minElapsed > elapsed {
			api.minElapsed = elapsed
		}

		// 是否有异常抛出"
		if event.GetExceptionInfo() != nil {
			api.errCount++
		}

		api.averageElapsed = float64(api.elapsed) / float64(api.count)
	}

	return nil
}

var gInserRPCDetailsRecord string = ` INSERT INTO api_details_stats (app_name, api, input_date, api_id, ser_type, elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count) VALUES (?,?,?,?,?,?,?,?,?,?,?);`

// eventRecord ...
func (spanEvents *SpanEvents) eventRecord(app *App, recordTime int64) error {
	for apiID, spanEvent := range spanEvents.spanEvents {
		if err := gAnalyze.cql.Session.Query(gInserRPCDetailsRecord,
			app.AppName,
			spanEvents.rpc,
			recordTime,
			apiID,
			spanEvent.serType,
			spanEvent.elapsed,
			spanEvent.maxElapsed,
			spanEvent.minElapsed,
			spanEvent.averageElapsed,
			spanEvent.count,
			spanEvent.errCount,
		).Exec(); err != nil {
			g.L.Warn("eventRecord error", zap.String("error", err.Error()), zap.String("SQL", gInserRPCDetailsRecord))
		}
	}
	return nil
}

// NewSpanEvents ...
func NewSpanEvents() *SpanEvents {
	return &SpanEvents{
		spanEvents: make(map[int32]*SpanEvent),
	}
}

// SpanEvent ...
type SpanEvent struct {
	serType        int
	elapsed        int
	maxElapsed     int
	minElapsed     int
	averageElapsed float64
	count          int
	errCount       int
}

// NewSpanEvent ...
func NewSpanEvent() *SpanEvent {
	return &SpanEvent{}
}
