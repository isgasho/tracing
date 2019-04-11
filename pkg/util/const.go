package util

// 第一级类型
const (
	TypeOfSkywalking byte = 1 // 	Skywalking 监控数据 uint16(iota + 1)
	TypeOfPinpoint   byte = 2 // 	Pinpoint 日志数据
	TypeOfCmd        byte = 3 // 	指令包 数据
	TypeOfLog        byte = 4 // 	日志数据
	TypeOfSystem     byte = 5 // 	系统数据
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
	TypeOfTraceSegment            uint16 = 6 // trace信息
)

// 监控报文类型Pinpoint
const (
	TypeOfTCPData         uint16 = 1 // tcp
	TypeOfUDPData         uint16 = 2 // udp
	TypeOfRegister        uint16 = 3 // 注册app id
	TypeOfAgentInfo       uint16 = 4
	TypeOfSQLMetaData     uint16 = 5
	TypeOfAPIMetaData     uint16 = 6
	TypeOfStringMetaData  uint16 = 7
	TypeOfTSpan           uint16 = 8
	TypeOfTSpanChunk      uint16 = 9
	TypeOfTAgentStat      uint16 = 10
	TypeOfTAgentStatBatch uint16 = 11
	TypeOfAgentOffline    uint16 = 12 // Agent 下线
)

// 系统监控数据类型
const (
	TypeOfCPU        uint16 = 1 // cpu
	TypeOfSystemload uint16 = 2 // Systemload
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

// 运行环境
const (
	TypeOfEnvJAVA int32 = 1
	TypeOfEnvGO   int32 = 2
)

// 版本类型
const (
	VersionOf01    byte = 1
	VersionVERSION      = 0
)

const (
	// SpanInsert string = `INSERT INTO traces (trace_id, trace_segment_id, span_id, app_id, instance_id, span_type, span_layer, start_time, end_time, parent_span_id, operation_id, is_error, refs, tags, logs) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	// JVMInsert  string = `INSERT INTO jvms (app_name, instance_id, report_time, jvms) VALUES (?,?,?,?)`
	// AgentInsert     string = `INSERT INTO agents (app_name, agent_id, ser_type, socket_id, host_name, ip, pid, version, start_time, is_live, is_container, end_time, operating_env) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);`
	AgentInsert     string = "INSERT into `agent`(`app_name`, `agent_id`, `ser_type`, `host_name`, `ip`, `pid`, `version`, `start_time`, `end_time`, `is_live`, `is_container`, `operating_env`) values ( '%s', '%s', %d, '%s', '%s', %d, '%s', %d, %d, %t, %t, %d)"
	AgentUpdate     string = "UPDATE `agent` SET ser_type=%d, host_name='%s', ip='%s', pid=%d, version='%s', start_time=%d, end_time=%d, is_live=%t, is_container=%t, operating_env=%d  WHERE app_name='%s' and agent_id='%s';"
	AgentOffLine    string = "UPDATE `agent` SET is_live=%t , end_time=%d WHERE app_name='%s' and agent_id='%s';"
	AgentInfoInsert string = "UPDATE `agent` SET agent_info='%q' WHERE app_name='%s' and agent_id='%s';"
)