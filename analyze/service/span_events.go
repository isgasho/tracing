package service

import (
	"log"

	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
)

// SpanEvents ...
type SpanEvents struct {
	spanEvents map[int32]*SpanEvent
}

var gCounterQueryAPI string = `SELECT api_info FROM agent_apis WHERE api_id=? 
		and agent_id=? and app_name=? and start_time=?;`

func (spanEvents *SpanEvents) eventsCounter(events []*trace.TSpanEvent) error {
	for _, event := range events {
		api, ok := spanEvents.spanEvents[event.GetApiId()]
		if !ok {
			api = NewSpanEvent()
			spanEvents.spanEvents[event.GetApiId()] = api
		}
		// event.SpanId
		log.Println("event.StartElapsed", event.StartElapsed)
		log.Println("event.EndElapsed", event.EndElapsed)
		api.count++
		// api.elapsed += event.EndElapsed
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
