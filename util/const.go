package util

// 第一级类型
const (
	TypeOfSkywalking byte = 1 // 	Skywalking 监控数据 uint16(iota + 1)
	TypeOfPinpoint   byte = 2 // 	Pinpoint 日志数据
	TypeOfCmd        byte = 3 // 	指令包 数据
	TypeOfLog        byte = 4 // 	日志数据
)

// 指令报文类型
const (
	TypeOfPing uint16 = 100 // 	Skywalking 监控数据 uint16(iota + 1) // 	Skywalking 监控数据
)

// 监控报文类型SKYWalking
const (
	TypeOfAppRegister             uint16 = 1 // 注册app id
	TypeOfAppRegisterInstance     uint16 = 2 // 注册app实例ID
	TypeOfSerNameDiscoveryService uint16 = 3 // 注册Api
	TypeOfNewworkAddrRegister     uint16 = 4 // 注册Addr
	TypeOfJVMMetrics              uint16 = 5 // JVM信息
)

// 其他控制类型
const (
	// MaxMessageSize max message size
	MaxMessageSize    int  = 16 * 1024
	TypeOfCompressYes byte = 1 // 数据压缩
	TypeOfCompressNo  byte = 2 // 数据不压缩
	TypeOfSyncYes     byte = 1 // 同步
	TypeOfSyncNo      byte = 2 // 非同步
)

// 版本类型
const (
	VersionOf01 byte = 1
)
