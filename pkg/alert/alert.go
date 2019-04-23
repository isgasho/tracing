package alert

// Data 告警存储数据
type Data struct {
	Type      int    `msg:"type"`
	AppName   string `msg:"name"`
	AgentID   string `msg:"id"`
	InputDate int64  `msg:"time"`
	Payload   []byte `msg:"payload"`
}

// NewData ...
func NewData() *Data {
	return &Data{}
}

// NewAPIs ...
func NewAPIs() *APIs {
	return &APIs{
		APIS: make(map[string]*API),
	}
}

// APIs ..
type APIs struct {
	APIS map[string]*API `msg:"apis"`
}

// API API信息
type API struct {
	Desc     string `msg:"desc"`
	Count    int    `msg:"count"`
	Errcount int    `msg:"errcount"`
	Duration int32  `msg:"duration"`
}

// SQLs sqls
type SQLs struct {
	SQLs map[int32]*SQL `msg:"sqls"`
}

// NewSQLs ...
func NewSQLs() *SQLs {
	return &SQLs{
		SQLs: make(map[int32]*SQL),
	}
}

// SQL sql计算
type SQL struct {
	ID       int32 `msg:"desc"`
	Count    int   `msg:"count"`
	Errcount int   `msg:"errcount"`
	Duration int32 `msg:"duration"`
}
