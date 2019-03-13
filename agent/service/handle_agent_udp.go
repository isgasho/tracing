package service

import (
	"fmt"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/util"
	"go.uber.org/zap"
)

func handleAgentUDP(data []byte) (*util.SpanDataModel, error) {
	spanModel := util.NewSpanDataModel()
	tStruct := thrift.Deserialize(data)
	switch m := tStruct.(type) {
	case *trace.TSpan:
		spanModel.Type = util.TypeOfTSpan
		spanModel.Spans = data
		break
	case *trace.TSpanChunk:
		spanModel.Type = util.TypeOfTSpanChunk
		spanModel.Spans = data
		break
	case *pinpoint.TAgentStat:
		spanModel.Type = util.TypeOfTAgentStat
		spanModel.Spans = data
		break
	case *pinpoint.TAgentStatBatch:
		spanModel.Type = util.TypeOfTAgentStatBatch
		spanModel.Spans = data
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("unknow type %t", m)))
		return nil, fmt.Errorf("unknow type %t", m)
	}
	return spanModel, nil
}
