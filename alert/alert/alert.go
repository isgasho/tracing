package alert

// Alert 告警服务
type Alert struct {
}

// New new alert
func New() *Alert {
	return &Alert{}
}

// Start start server
func (a *Alert) Start() error {
	return nil
}

// Close stop server
func (a *Alert) Close() error {
	return nil
}
