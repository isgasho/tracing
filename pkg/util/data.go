package util

/* 公共数据结构 */

type Alert struct {
	Name     string  `json:"name" cql:"name"`
	Type     string  `json:"type" cql:"type"`
	Label    string  `json:"label" cql:"label"`
	Compare  int     `json:"compare" cql:"compare"`
	Unit     string  `json:"unit" cql:"unit"`
	Duration int     `json:"duration" cql:"duration"`
	Value    float64 `json:"value" cql:"value"`
}
