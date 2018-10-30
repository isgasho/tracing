package util

type RefType int32

const (
	RefType_CrossProcess RefType = 0
	RefType_CrossThread  RefType = 1
)

// Span ...
type Span struct {
	TraceID         string     `msg:"tid"`
	SpanID          int32      `msg:"sid"`
	Refs            []*SpanRef `msg:"rfs"`
	StartTime       int64      `msg:"st"`
	EndTime         int64      `msg:"et"`
	ParentSpanID    int32      `msg:"pid"`
	OperationNameID int32      `msg:"oid"`
	OperationName   string     `msg:"on"`
}

// SpanRef ...
type SpanRef struct {
	TraceID string  `msg:"tid"`
	SpanID  string  `msg:"sd"`
	RefType RefType `msg:"rt"`
}
