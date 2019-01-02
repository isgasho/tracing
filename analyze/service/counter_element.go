package service

import (
	"github.com/mafanr/g/utils"
)

// Element ...
type Element struct {
	urls       *SpanURLs
	events     *SpanEvents
	exceptions *SpanExceptions
	stats      *AgentStats
	jvm        *JVM
	api        *API
	FUNC       *FUNC
}

// NewElement ...
func NewElement() *Element {
	return &Element{
		urls:       NewSpanURLs(),
		events:     NewSpanEvents(),
		exceptions: NewSpanExceptions(),
		stats:      NewAgentStats(),
		// api:  NewAPI(),
		// FUNC: NewFUNC(),
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

// JVM ...
type JVM struct {
}

// NewJVM ...
func NewJVM() *JVM {
	return &JVM{}
}

// API ...
type API struct {
}

// NewAPI ...
func NewAPI() *API {
	return &API{}
}

// FUNC ...
type FUNC struct {
}

// NewFUNC ...
func NewFUNC() *FUNC {
	return &FUNC{}
}
