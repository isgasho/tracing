package util

type GCPhrase int32

const (
	GCPhrase_NEW GCPhrase = 0
	GCPhrase_OLD GCPhrase = 1
)

type PoolType int32

const (
	PoolType_CODE_CACHE_USAGE PoolType = 0
	PoolType_NEWGEN_USAGE     PoolType = 1
	PoolType_OLDGEN_USAGE     PoolType = 2
	PoolType_SURVIVOR_USAGE   PoolType = 3
	PoolType_PERMGEN_USAGE    PoolType = 4
	PoolType_METASPACE_USAGE  PoolType = 5
)

// JVMS jvm批量数据
type JVMS struct {
	AppName    string `msg:"n" cql:"app_name"`
	InstanceID int32  `msg:"i" cql:"instance_id"`
	Time       int64  `msg:"t" cql:"report_time"`
	JVMs       []*JVM `msg:"jvms" cql:"jvms"`
}

// JVM jvm 信息
type JVM struct {
	Time       int64         `msg:"t"  cql:"time"`
	CPU        *CPU          `msg:"c"  cql:"cpu"`
	Memory     []*Memory     `msg:"m"  cql:"memory"`
	MemoryPool []*MemoryPool `msg:"mp" cql:"memory_pool"`
	Gc         []*GC         `msg:"gc" cql:"gc"`
}

// CPU ...
type CPU struct {
	UsagePercent float64 `msg:"up" cql:"usage_percent"`
}

// MemoryPool ...
type MemoryPool struct {
	Type     PoolType `msg:"t" cql:"type"`
	Init     int64    `msg:"i" cql:"init"`
	Max      int64    `msg:"m" cql:"max"`
	Used     int64    `msg:"u" cql:"used"`
	Commited int64    `msg:"c" cql:"commited"`
}

// Memory ...
type Memory struct {
	IsHeap    bool  `msg:"ih" cql:"is_heap"`
	Init      int64 `msg:"i"  cql:"init"`
	Max       int64 `msg:"m"  cql:"max"`
	Used      int64 `msg:"u"  cql:"used"`
	Committed int64 `msg:"c"  cql:"committed"`
}

// GC ...
type GC struct {
	Phrase GCPhrase `msg:"p" cql:"phrase"`
	Count  int64    `msg:"c" cql:"count"`
	Time   int64    `msg:"t" cql:"time"`
}
