package util

type SpanType int32

const (
	SpanType_Entry SpanType = 0
	SpanType_Exit  SpanType = 1
	SpanType_Local SpanType = 2
)

// RefType ...
type RefType int32

const (
	RefType_CrossProcess RefType = 0
	RefType_CrossThread  RefType = 1
)

type SpanLayer int32

const (
	SpanLayer_Unknown      SpanLayer = 0
	SpanLayer_Database     SpanLayer = 1
	SpanLayer_RPCFramework SpanLayer = 2
	SpanLayer_Http         SpanLayer = 3
	SpanLayer_MQ           SpanLayer = 4
	SpanLayer_Cache        SpanLayer = 5
)

// Span ...
type Span struct {
	TraceID         string                `msg:"tid"`
	SpanID          int32                 `msg:"sid"`
	AppID           int32                 `msg:"aid"`
	InstanceID      int32                 `msg:"inid"`
	SpanType        SpanType              `msg:"sty"`
	SpanLayer       SpanLayer             `msg:"sly"`
	Refs            []*SpanRef            `msg:"rfs"`
	StartTime       int64                 `msg:"st"`
	EndTime         int64                 `msg:"et"`
	ParentSpanID    int32                 `msg:"pid"`
	OperationNameID int32                 `msg:"oid"`
	IsError         bool                  `msg:"ie"`
	Tags            []*KeyWithStringValue `msg:"tags"`
	Logs            []*LogMessage         `msg:"logs"`
	// OperationName   string                `msg:"on"`
}

// SpanRef ...
type SpanRef struct {
	TraceID string  `msg:"tid" cql:"trace_id"`
	SpanID  int32   `msg:"sid" cql:"span_id"`
	RefType RefType `msg:"rt"  cql:"ref_type"`
}
