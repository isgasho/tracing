package service

import (
	"github.com/imdevlab/g/utils"
)

// Element ...
type Element struct {
	apis       *SpanAPIs
	events     *SpanEvents
	exceptions *SpanExceptions
	stats      *AgentStats
	sqls       *SpanSQLs
}

// NewElement ...
func NewElement() *Element {
	return &Element{
		apis:       NewSpanAPIs(),
		events:     NewSpanEvents(),
		exceptions: NewSpanExceptions(),
		stats:      NewAgentStats(),
		sqls:       NewSpanSQLs(),
	}
}

// GetElements ...
func GetElements(startTime, endTime int64) map[int64]*Element {
	es := make(map[int64]*Element)
	st, _ := utils.MSToTime(startTime)
	min := ((endTime - startTime) / 1000) / 60
	startIndexTime := st.Unix() - int64(st.Second())

	for index := 0; index < int(min); index++ {
		es[startIndexTime+int64(index*60)] = NewElement()
	}

	return es
}
