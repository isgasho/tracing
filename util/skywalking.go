package util

// SkywalkingPacket ...
type SkywalkingPacket struct {
	Type    uint16 `msg:"type"`
	Payload []byte `msg:"payload"`
}

// KeyWithIntegerValue ...
type KeyWithIntegerValue struct {
	Key   string `msg:"k"`
	Value int32  `msg:"v"`
}

// KeyWithStringValue ...
type KeyWithStringValue struct {
	Key   string `msg:"k"`
	Value string `msg:"v"`
}

// API ...
type API struct {
	AppID    int32  `msg:"aid"`
	SerID    int32  `msg:"sid"`
	SerName  string `msg:"sn"`
	SpanType int32  `msg:"st"`
}

// API ...
type SerNameDiscoveryServices struct {
	SerNames []*API `msg:"as"`
}
