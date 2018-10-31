package util

// LogMessage log
type LogMessage struct {
	Time int64                 `msg:"t" cql:"ts"`
	Data []*KeyWithStringValue `msg:"d" cql:"fields"`
}
