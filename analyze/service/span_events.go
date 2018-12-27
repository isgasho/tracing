package service

import (
	"log"

	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
)

// SpanEvents ...
type SpanEvents struct {
	spanEvents map[string]*SpanEvent
}

var gCounterQueryAPI string = `SELECT api_info FROM agent_apis WHERE api_id=? 
		and agent_id=? and app_name=? and start_time=?;`

func (spanEvents *SpanEvents) eventsCounter(traceID string, spanID int64, agentID string, agentTtartTime int64, events []*trace.TSpanEvent) error {

	for _, event := range events {
		log.Println(event.GetApiId())
		log.Println(event.GetEndElapsed())
		log.Println(event.GetServiceType())
		log.Println(traceID, spanID, agentID, event.GetApiId(), agentTtartTime)

	}

	return nil
}

// NewSpanEvents ...
func NewSpanEvents() *SpanEvents {
	return &SpanEvents{
		spanEvents: make(map[string]*SpanEvent),
	}
}

// SpanEvent ...
type SpanEvent struct {
}

// NewSpanEvent ...
func NewSpanEvent() *SpanEvent {
	return &SpanEvent{}
}
