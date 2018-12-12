package util

// PinpointData pinpoint data
type PinpointData struct {
	Type     uint16           `msg:"type"` // tcp or udp
	AppName  string           `msg:"appName"`
	AgentID  string           `msg:"agentID"`
	SpanTime int64            `msg:"spanTime"`
	Payload  []*SpanDataModel `msg:"payload"`
}

// NewPinpointData ...
func NewPinpointData() *PinpointData {
	return &PinpointData{
		Payload: make([]*SpanDataModel, 0),
	}
}

// NewSpanDataModel ...
func NewSpanDataModel() *SpanDataModel {
	return &SpanDataModel{}
}

// type DataType : SpanV2 SpanChunk AgentStat AgentStatBatch
// SpanDataModel data
type SpanDataModel struct {
	Type  uint16 `msg:"type"`
	Spans []byte `msg:"spans"`
}
