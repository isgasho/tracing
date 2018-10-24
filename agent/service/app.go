package service

// AppInfo ...
type AppInfo struct {
	ID   string
	Name string
	Host string
}

// NewAppInfo ...
func NewAppInfo() *AppInfo {
	return &AppInfo{}
}
