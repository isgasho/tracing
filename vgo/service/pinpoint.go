package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/thrift"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/pinpoint"
	"github.com/mafanr/vgo/proto/pinpoint/thrift/trace"
	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
	"log"
	"net"
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
				msgpack.Unmarshal(value.Spans, appInfo)
				log.Print("获取到AgentInfo", appInfo)
				break
			case util.TypeOfAgentSEND:
				DealTCPRequestResponse(value.Spans)
				break
			}
		}
		break
	case util.TypeOfUDPData:

		break
	}
	return nil
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

//if err := v.dealSkywalking(conn, packet); err != nil {
//	g.L.Warn("agentWork:v.dealSkywalking", zap.String("error", err.Error()))
//	return
//}
