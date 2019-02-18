package util

// Metric ...
type Metric struct {
	Name     string                 `msg:"n"`
	Tags     map[string]string      `msg:"ts"`
	Fields   map[string]interface{} `msg:"f"`
	Time     int64                  `msg:"t"`
	Interval int                    `msg:"i"`
}
