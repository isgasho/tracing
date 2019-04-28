package alert

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
