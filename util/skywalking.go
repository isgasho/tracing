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
	Key   string `msg:"k" cql:"key"`
	Value string `msg:"v" cql:"value"`
}

// API ...
type API struct {
	AppID    int32  `msg:"aid"`
	SerID    int32  `msg:"sid"`
	SerName  string `msg:"sn"`
	SpanType int32  `msg:"st"`
}

// SerNameDiscoveryServices ...
type SerNameDiscoveryServices struct {
	SerNames []*API `msg:"as"`
}

// RegisterAddrs ...
type RegisterAddrs struct {
	AppID int32                  `msg:"aid"`
	Addrs []*KeyWithIntegerValue `msg:"Addr"`
}
