package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/imdevlab/vgo/agent/misc"

	"github.com/imdevlab/vgo/util"

	"github.com/imdevlab/g"
	"github.com/imdevlab/vgo/proto/pinpoint/proto"
	"github.com/imdevlab/vgo/proto/pinpoint/thrift"
	"github.com/imdevlab/vgo/proto/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/vgo/proto/pinpoint/thrift/trace"
	"go.uber.org/zap"
)

// var agentStartTimeTest int64

// func threadDump(streamID int, conn net.Conn) error {
// 	var rePacket proto.Packet

// 	threadDump := command.NewTCommandThreadDump()
// 	threadDump.Type = command.TThreadDumpType(1)
// 	// transfer := command.NewTCommandTransfer()
// 	// transfer.Payload = thrift.SerializeNew(threadDump)

// 	applicationStreamCreate := proto.NewApplicationStreamCreate()
// 	applicationStreamCreate.Payload = thrift.SerializeNew(threadDump) //transfer.Payload //thrift.Serialize(transfer)
// 	applicationStreamCreate.ChannelID = streamID
// 	rePacket = applicationStreamCreate

// 	body, err := rePacket.Encode()
// 	if err != nil {
// 		g.L.Warn("rePacket.Encode", zap.String("error", err.Error()))
// 		return err
// 	}

// 	if _, err := conn.Write(body); err != nil {
// 		g.L.Warn("conn.Write", zap.String("error", err.Error()))
// 		return err
// 	}
// 	return nil
// }

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
	// streamID := int(rand.Int31n(1000))
	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)
	// 		err := threadDump(streamID, conn)
	// 		if err != nil {
	// 			log.Println(err)
	// 			break
	// 		}
	// 		streamID++
	// 		log.Println("发送成功")

	// 		time.Sleep(10 * time.Second)

	// 	}
	// }()

	isRecvOffline := false
	defer func() {
		if !isRecvOffline {
			// sdk客户端断线
			gAgent.agentInfo.IsLive = false
			gAgent.isReportAgentInfo = true
			gAgent.agentInfo.EndTimestamp = time.Now().UnixNano() / 1e6
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

			spanModel, err := handleAgentTCP(applicationSend.GetPayload())
			if err != nil {
				g.L.Warn("AgentStat:handleAgentUDP", zap.String("error", err.Error()))
				return err
			}

			pinpoint.tcpChan <- spanModel
			break

		case proto.APPLICATION_REQUEST:
			// g.L.Debug("agentInfo", zap.String("type", "APPLICATION_REQUEST"))
			applicationRequest := proto.NewApplicationRequest()
			if err := applicationRequest.Decode(conn, reader); err != nil {
				g.L.Warn("applicationRequest.Decode", zap.String("error", err.Error()))
				return err
			}

			spanModel, err := handleAgentTCP(applicationRequest.GetPayload())
			if err != nil {
				g.L.Warn("AgentStat:handleAgentUDP", zap.String("error", err.Error()))
				return err
			}
			pinpoint.tcpChan <- spanModel

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

			// log.Println(applicationStreamCreateFail)
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
			agentInfo := util.NewAgentInfo()
			if err := json.Unmarshal(controlHandShake.GetPayload(), agentInfo); err != nil {
				g.L.Warn("json.Unmarshal", zap.String("error", err.Error()))
				return err
			}

			agentInfo.IsLive = true
			agentInfo.IsContainer = misc.Conf.Agent.IsContainer

			// 运行环境
			agentInfo.OperatingEnv = misc.Conf.Agent.OperatingEnv

			gAgent.agentInfo = agentInfo
			gAgent.isReportAgentInfo = true

			g.L.Info("agentInfo", zap.Any("agentInfo", agentInfo))

			isRePacket = true
			rePacket, err = createResponse(controlHandShake)
			if err != nil {
				g.L.Warn("createResponse", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.CONTROL_CLIENT_CLOSE:

			g.L.Debug("agentInfo", zap.String("type", "CONTROL_CLIENT_CLOSE"))
			// sdk客户端断线
			gAgent.agentInfo.IsLive = false
			gAgent.isReportAgentInfo = true
			gAgent.agentInfo.EndTimestamp = time.Now().UnixNano() / 1e6

			isRecvOffline = true

			controlClientClose := proto.NewControlClientClose()
			if err := controlClientClose.Decode(conn, reader); err != nil {
				g.L.Warn("controlClientClose.Decode", zap.String("error", err.Error()))
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

// handleAgentTCP ...
func handleAgentTCP(message []byte) (*util.SpanDataModel, error) {
	spanModel := util.NewSpanDataModel()
	tStruct := thrift.Deserialize(message)
	switch m := tStruct.(type) {
	case *pinpoint.TAgentInfo:
		spanModel.Type = util.TypeOfAgentInfo
		spanModel.Spans = message
		break
	case *trace.TSqlMetaData:
		spanModel.Type = util.TypeOfSQLMetaData
		spanModel.Spans = message
		break
	case *trace.TApiMetaData:
		spanModel.Type = util.TypeOfAPIMetaData
		spanModel.Spans = message
		break
	case *trace.TStringMetaData:
		spanModel.Type = util.TypeOfStringMetaData
		spanModel.Spans = message
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("unknow type %t", m)))
		return nil, fmt.Errorf("unknow type %t", m)
	}
	return spanModel, nil
}
