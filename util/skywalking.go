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

// AppRegister ...
type AppRegister struct {
	Name string `msg:"n"`
	Code int32  `msg:"c"`
}
