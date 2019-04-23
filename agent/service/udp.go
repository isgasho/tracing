package service

import (
	"fmt"

	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/network"

	"github.com/imdevlab/tracing/pkg/pinpoint/thrift"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"go.uber.org/zap"
)

func udpRead(data []byte) (*network.Spans, error) {
	spans := network.NewSpans()
	tStruct := thrift.Deserialize(data)
	switch m := tStruct.(type) {
	case *trace.TSpan:
		spans.Type = constant.TypeOfTSpan
		spans.Spans = data
		break
	case *trace.TSpanChunk:
		spans.Type = constant.TypeOfTSpanChunk
		spans.Spans = data
		break
	case *pinpoint.TAgentStat:
		spans.Type = constant.TypeOfTAgentStat
		spans.Spans = data
		break
	case *pinpoint.TAgentStatBatch:
		spans.Type = constant.TypeOfTAgentStatBatch
		spans.Spans = data
		break
	default:
		logger.Warn("unknown type", zap.String("type", fmt.Sprintf("unknow type %t", m)))
		return nil, fmt.Errorf("unknow type %t", m)
	}
	return spans, nil
}
