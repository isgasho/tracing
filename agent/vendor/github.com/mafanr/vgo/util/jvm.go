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
	// AppID      int32  `msg:"i"`
	AppName    string `msg:"n"`
	InstanceID int32  `msg:"i"`
	Time       int64  `msg:"t"`
	JVMs       []*JVM `msg:"jvms"`
}

// JVM jvm 信息
type JVM struct {
	Time       int64         `msg:"t"`
	CPU        *CPU          `msg:"c"`
	Memory     []*Memory     `msg:"m"`
	MemoryPool []*MemoryPool `msg:"mp"`
	Gc         []*GC         `msg:"gc"`
}

// CPU ...
type CPU struct {
	UsagePercent float64 `msg:"up"`
}

// MemoryPool ...
type MemoryPool struct {
	Type     PoolType `msg:"t"`
	Init     int64    `msg:"i"`
	Max      int64    `msg:"m"`
	Used     int64    `msg:"u"`
	Commited int64    `msg:"c"`
}

// Memory ...
type Memory struct {
	IsHeap    bool  `msg:"ih"`
	Init      int64 `msg:"i"`
	Max       int64 `msg:"m"`
	Used      int64 `msg:"u"`
	Committed int64 `msg:"c"`
}

// GC ...
type GC struct {
	Phrase GCPhrase `msg:"p"`
	Count  int64    `msg:"c"`
	Time   int64    `msg:"t"`
}
