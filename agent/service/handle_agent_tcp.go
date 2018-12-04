package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"

	"github.com/mafanr/vgo/util"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/proto"
	"github.com/mafanr/vgo/proto/pinpoint/thrift"
	"go.uber.org/zap"
)

func (pinpoint *Pinpoint) agentInfo(conn net.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	reader := bufio.NewReaderSize(conn, proto.TCP_MAX_PACKET_SIZE)
	buf := make([]byte, 2)
	for {
		// response packet
		var rePacket proto.Packet
		var err error
		// response flag
		isRePacket := false
		if _, err := io.ReadFull(reader, buf); err != nil {
			g.L.Error("io.ReadFull", zap.Error(err))
			return err
		}
		// read packet type
		packetType := int16(binary.BigEndian.Uint16(buf[0:2]))
		switch packetType {
		case proto.APPLICATION_SEND:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_SEND"))
			applicationSend := proto.NewApplicationSend()
			if err := applicationSend.Decode(conn, reader); err != nil {
				g.L.Warn("applicationSend.Decode", zap.String("error", err.Error()))
				return err
			}

			if err := pinpoint.reportSEND(applicationSend.GetPayload()); err != nil {
				g.L.Warn("pinpoint.reportSEND", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_REQUEST:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_REQUEST"))
			applicationRequest := proto.NewApplicationRequest()
			if err := applicationRequest.Decode(conn, reader); err != nil {
				g.L.Warn("applicationRequest.Decode", zap.String("error", err.Error()))
				return err
			}

			if err := pinpoint.reportSEND(applicationRequest.GetPayload()); err != nil {
				g.L.Warn("pinpoint.reportSEND", zap.String("error", err.Error()))
				return err
			}

			tResult := proto.DealRequestResponse(applicationRequest)
			response := proto.NewApplicationResponse()
			response.RequestID = applicationRequest.GetRequestID()
			response.Payload = thrift.Serialize(tResult)
			isRePacket = true
			rePacket = response
			break

		case proto.APPLICATION_RESPONSE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_RESPONSE"))
			applicationResponse := proto.NewApplicationResponse()
			if err := applicationResponse.Decode(conn, reader); err != nil {
				g.L.Warn("applicationResponse.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE"))
			applicationStreamCreate := proto.NewApplicationStreamCreate()
			if err := applicationStreamCreate.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamCreate.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CLOSE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CLOSE"))
			applicationStreamClose := proto.NewApplicationStreamClose()
			if err := applicationStreamClose.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamClose.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE_SUCCESS:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE_SUCCESS"))
			applicationStreamCreateSuccess := proto.NewApplicationStreamCreateSuccess()
			if err := applicationStreamCreateSuccess.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamCreateSuccess.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE_FAIL:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE_FAIL"))
			applicationStreamCreateFail := proto.NewApplicationStreamCreateFail()
			if err := applicationStreamCreateFail.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamCreateFail.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_RESPONSE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_RESPONSE"))
			applicationStreamResponse := proto.NewApplicationStreamResponse()
			if err := applicationStreamResponse.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamResponse.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_PING:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_PING"))
			applicationStreamPing := proto.NewApplicationStreamPing()
			if err := applicationStreamPing.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamPing.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_PONG:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_PONG"))
			applicationStreamPong := proto.NewApplicationStreamPong()
			if err := applicationStreamPong.Decode(conn, reader); err != nil {
				g.L.Warn("applicationStreamPong.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.CONTROL_HANDSHAKE:
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_HANDSHAKE"))
			controlHandShake := proto.NewControlHandShake()
			if err := controlHandShake.Decode(conn, reader); err != nil {
				g.L.Warn("controlHandShake.Decode", zap.String("error", err.Error()))
				return err
			}
			appInfo := util.NewAgentInfo()
			if err := json.Unmarshal(controlHandShake.GetPayload(), appInfo); err != nil {
				g.L.Warn("json.Unmarshal", zap.String("error", err.Error()))
				return err
			}

			appInfo.AppName = gAgent.appName
			gAgent.setAgentInfo(appInfo.AgentID)

			if err := pinpoint.reportAgentInfo(appInfo); err != nil {
				g.L.Warn("pinpoint.reportAgentInfo", zap.String("error", err.Error()))
				return err
			}

			g.L.Info("agentInfo", zap.Any("body", appInfo))

			isRePacket = true
			rePacket, err = createResponse(controlHandShake)
			if err != nil {
				g.L.Warn("createResponse", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.CONTROL_CLIENT_CLOSE:
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_CLIENT_CLOSE"))
			controlClientClose := proto.NewControlClientClose()
			if err := controlClientClose.Decode(conn, reader); err != nil {
				g.L.Warn("controlClientClose.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.CONTROL_PING:
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_PING"))
			controlPing := proto.NewControlPing()
			if err := controlPing.Decode(conn, reader); err != nil {
				g.L.Warn("controlPing.Decode", zap.String("error", err.Error()))
				return err
			}
			isRePacket = true
			controlPong := proto.NewControlPong()
			rePacket = controlPong
			break

		case proto.CONTROL_PING_SIMPLE:
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_PING_SIMPLE"))
			isRePacket = true
			controlPong := proto.NewControlPong()
			rePacket = controlPong
			break

		case proto.CONTROL_PING_PAYLOAD:
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_PING"))
			controlPingPayload := proto.NewControlPing()
			if err := controlPingPayload.Decode(conn, reader); err != nil {
				g.L.Warn("controlPingPayload.Decode", zap.String("error", err.Error()))
				return err
			}
			break

		default:
			g.L.Warn("unaware packet Type", zap.Int16("packetType", packetType))
		}

		if isRePacket {
			body, err := rePacket.Encode()
			if err != nil {
				g.L.Warn("rePacket.Encode", zap.String("error", err.Error()))
				return err
			}

			if _, err := conn.Write(body); err != nil {
				g.L.Warn("conn.Write", zap.String("error", err.Error()))
				return err
			}
		}
	}
}
