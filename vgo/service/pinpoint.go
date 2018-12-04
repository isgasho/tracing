package service

import (
	"log"
	"net"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/pinpoint"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

// Pinpoint ...
type Pinpoint struct {
}

// NewPinpoint ...
func NewPinpoint() *Pinpoint {
	return &Pinpoint{}
}

// dealUpload 处理pinpoint上行数据
func (pinpoint *Pinpoint) dealUpload(conn net.Conn, inPacket *util.VgoPacket) error {
	packet := &util.PinpointData{}
	if err := msgpack.Unmarshal(inPacket.Payload, packet); err != nil {
		g.L.Warn("dealUpload:msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}

	switch packet.Type {
	case util.TypeOfTCPData:
		for _, value := range packet.Payload {
			switch value.Type {
			case util.TypeOfAgentInfo:
				appInfo := util.NewAgentInfo()
				if err := msgpack.Unmarshal(value.Spans, appInfo); err != nil {
					g.L.Warn("dealUpload:msgpack.Unmarshal", zap.String("error", err.Error()))
					return err
				}
				break
			case util.TypeOfAgentSEND:
				DealTCPRequestResponse(value.Spans)
				break
			default:
				g.L.Warn("unknow type")
				break
			}
		}
		break
	case util.TypeOfUDPData:
		for _, value := range packet.Payload {
			DealUDPRequestResponse(value.Spans)
		}
		break
	default:
		g.L.Warn("unknow type")
	}
	return nil
}

// DealUDPRequestResponse ...
func DealUDPRequestResponse(data []byte) {
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

// DealTCPRequestResponse ...
func DealTCPRequestResponse(message []byte) error {
	tStruct := thrift.Deserialize(message)
	switch m := tStruct.(type) {
	case *pinpoint.TAgentInfo:
		log.Println("pinpoint.TAgentInfo")
		break
	case *trace.TSqlMetaData:
		log.Println("pinpoint.TSqlMetaData")
		break
	case *trace.TApiMetaData:
		log.Println("pinpoint.TApiMetaData")
		break
	case *trace.TStringMetaData:
		log.Println("pinpoint.TStringMetaData")
		break
	default:
		g.L.Warn("unknown type", zap.Any("data", m))
	}
	return nil
}
