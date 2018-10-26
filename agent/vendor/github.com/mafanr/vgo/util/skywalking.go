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

// SerNameDiscoveryService ...
type SerNameDiscoveryService struct {
	SerName  string `msg:"sn"`
	SerID    int32  `msg:"sid"`
	SpanType int32  `msg:"st"`
}

// SerNameDiscoveryServices ...
type SerNameDiscoveryServices struct {
	AppCode int32                      `msg:"ac"`
	Sers    []*SerNameDiscoveryService `msg:"as"`
}

// // AppRegister ...
// type AppRegister struct {
// 	Name string `msg:"n"`
// 	Code int32  `msg:"c"`
// }
