package service

import (
	"encoding/json"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/pinpoint/proto"
	"go.uber.org/zap"
)

func createResponse(in proto.Packet) (proto.Packet, error) {
	packType := in.GetPacketType()
	switch packType {
	case proto.CONTROL_HANDSHAKE:
		resultMap := make(map[string]interface{})
		resultMap[proto.CODE] = proto.HANDSHAKE_DUPLEX_COMMUNICATION.Code
		resultMap[proto.SUB_CODE] = proto.HANDSHAKE_DUPLEX_COMMUNICATION.SubCode
		payload, err := json.Marshal(resultMap)
		if err != nil {
			g.L.Warn("json.Marshal", zap.String("error", err.Error()))
			return nil, err
		}
		controlHandShakeResponse := proto.NewControlHandShakeResponse()
		controlHandShakeResponse.Payload = payload
		controlHandShakeResponse.RequestID = in.GetRequestID()
		return controlHandShakeResponse, nil

	default:
		break
	}

	return nil, nil
}
