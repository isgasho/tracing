package util

// MetricData ...
type MetricData struct {
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
	Name     string                 `msg:"n"`
	Tags     map[string]string      `msg:"ts"`
	Fields   map[string]interface{} `msg:"f"`
	Time     int64                  `msg:"t"`
	Interval int                    `msg:"i"`
}

// NewMetric ...
func NewMetric() *Metric {
	return &Metric{}
}
