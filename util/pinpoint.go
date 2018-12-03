package util

// PinpointData pinpoint data
type PinpointData struct {
	Type      uint16           `msg:"type"` // tcp or udp
	AgentName string           `msg:"agentName"`
	AgentID   string           `msg:"agentID"`
	SpanTime  int64            `msg:"spanTime"`
	Payload   []*SpanDataModel `msg:"payload"`
}

func NewPinpointData() *PinpointData {
	return &PinpointData{
		Payload: make([]*SpanDataModel, 0),
	}
}

// type DataType : SpanV2 SpanChunk AgentStat AgentStatBatch
// SpanDataModel data
type SpanDataModel struct {
	Type  uint16 `msg:"type"`
	Spans []byte `msg:"spans"`
}
