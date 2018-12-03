package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/pinpoint/thrift"
	"github.com/mafanr/vgo/pinpoint/thrift/pinpoint"
	"github.com/mafanr/vgo/pinpoint/thrift/trace"
	"go.uber.org/zap"
)

func handleAgentUDP(data []byte) {
	tStruct := thrift.Deserialize(data)
	switch m := tStruct.(type) {
	case *trace.TSpan:
		g.L.Debug("udp", zap.String("type", "TSpan"))
		break
	case *trace.TSpanChunk:
		g.L.Debug("udp", zap.String("type", "TSpanChunk"))
		break
	case *pinpoint.TAgentStat:
		g.L.Debug("udp", zap.String("type", "TAgentStat"))
		break
	case *pinpoint.TAgentStatBatch:

		g.L.Debug("udp", zap.String("type", "TAgentStatBatch"))
		break
	default:
		g.L.Warn("unknow type", zap.Any("data", m))
	}
}
