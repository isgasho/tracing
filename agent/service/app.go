package service

// AppInfo ...
type AppInfo struct {
	Code int    // appCode
	ID   string // agentID
	Name string // app name
}

// NewAppInfo ...
func NewAppInfo() *AppInfo {
	return &AppInfo{}
}
