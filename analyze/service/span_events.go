package service

import (
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
)

// SpanEvents ...
type SpanEvents struct {
	spanEvents map[int32]*SpanEvent
}

var gCounterQueryAPI string = `SELECT api_info FROM agent_apis WHERE api_id=? 
		and agent_id=? and app_name=? and start_time=?;`

func (spanEvents *SpanEvents) eventsCounter(events []*trace.TSpanEvent, chunkEvents []*trace.TSpanEvent) error {
	for _, event := range events {
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

		// 是否有异常抛出
		if event.GetExceptionInfo() != nil {
			api.errCount++
		}

		api.averageElapsed = api.elapsed / api.count
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

		api.averageElapsed = api.elapsed / api.count
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
	averageElapsed int
	count          int
	errCount       int
}

// NewSpanEvent ...
func NewSpanEvent() *SpanEvent {
	return &SpanEvent{}
}
