package metric

// APIMap api被调用情况
type APIMap struct {
	APIS map[int32]*API
}

// NewAPIMap ...
func NewAPIMap() *APIMap {
	return &APIMap{
		APIS: make(map[int32]*API),
	}
}

// API 所有调用者信息
type API struct {
	Parents map[string]*APIMapInfo
}

// NewAPI ...
func NewAPI() *API {
	return &API{
		Parents: make(map[string]*APIMapInfo),
	}
}

// NewAPIMapInfo ...
func NewAPIMapInfo() *APIMapInfo {
	return &APIMapInfo{}
}

// APIMapInfo 调用信息
type APIMapInfo struct {
	Type           int16
	AccessCount    int
	AccessErrCount int
	AccessDuration int32
}
