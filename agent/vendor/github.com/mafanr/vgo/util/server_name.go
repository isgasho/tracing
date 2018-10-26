package util

// SerNameInfo ...
type SerNameInfo struct {
	SerID    int32  `db:"id" json:"id" msg:"id"`
	AppCode  int32  `db:"app_code" json:"app_code" msg:"app_code"`
	SerName  string `db:"server_name" json:"server_name" msg:"server_name"`
	SpanType int32  `db:"span_type" json:"span_type" msg:"span_type"`
}
