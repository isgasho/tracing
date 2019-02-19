package util

// MetricData ...
type MetricData struct {
	AppName string    `msg:"an"`
	AgentID string    `msg:"aid"`
	Time    int64     `msg:"t"`
	Payload []*Metric `msg:"p"`
}

// NewMetricData ...
func NewMetricData() *MetricData {
	return &MetricData{
		Payload: make([]*Metric, 0),
	}
}

// Metric ...
type Metric struct {
	Name     string                 `msg:"n"  json:"name"`
	Tags     map[string]string      `msg:"ts" json:"tags"`
	Fields   map[string]interface{} `msg:"f"  json:"fields"`
	Time     int64                  `msg:"t"  json:"time"`
	Interval int                    `msg:"i"  json:"interval"`
}

// NewMetric ...
func NewMetric() *Metric {
	return &Metric{}
}
