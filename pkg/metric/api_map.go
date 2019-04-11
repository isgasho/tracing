package metric

// APICallStats api被调用情况
type APICallStats struct {
	APIS map[int32]*API
}

// NewAPICallStats ...
func NewAPICallStats() *APICallStats {
	return &APICallStats{
		APIS: make(map[int32]*API),
	}
}

// API 所有调用者信息
type API struct {
	Parents map[string]*ParentInfo
}

// NewAPI ...
func NewAPI() *API {
	return &API{
		Parents: make(map[string]*ParentInfo),
	}
}
