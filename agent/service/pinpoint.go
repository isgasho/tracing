package service

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/agent/misc"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/imdevlab/tracing/pkg/pinpoint/proto"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

// Pinpoint p数据采集
type Pinpoint struct {
	tcpChan chan *network.Spans // tcp报文接收管道
	udpChan chan *network.Spans // udp报文接收管道
}

func newPinpoint() *Pinpoint {
	return &Pinpoint{
		tcpChan: make(chan *network.Spans, 100),
		udpChan: make(chan *network.Spans, 300),
	}
}

// Start 启动pinpoint采集服务
func (p *Pinpoint) Start() error {
	// 接收pinpoint三种数据
	go p.AgentInfo()
	go p.AgentStat()
	go p.AgentSpan()

	// pinpoint tcp、udp信息采集&上报
	go p.tcpCollector()
	go p.udpCollector()
	return nil
}

// Close 关闭pinpoint采集服务
func (p *Pinpoint) Close() error {
	return nil
}

// AgentInfo ...
func (p *Pinpoint) AgentInfo() error {
	l, err := net.Listen("tcp", misc.Conf.Pinpoint.InfoAddr)
	if err != nil {
		g.L.Fatal("AgentInfo", zap.Error(err))
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			g.L.Fatal("AgentInfo", zap.String("addr", misc.Conf.Pinpoint.InfoAddr), zap.Error(err))
			return err
		}
		go p.agentInfo(conn)
	}
}

// AgentStat ...
func (p *Pinpoint) AgentStat() error {
	addrInfo, _ := net.ResolveUDPAddr("udp", misc.Conf.Pinpoint.StatAddr)
	listener, err := net.ListenUDP("udp", addrInfo)
	if err != nil {
		g.L.Fatal("AgentStat ListenUDP", zap.String("addr", misc.Conf.Pinpoint.StatAddr), zap.String("error", err.Error()))
	}

	for {
		data := make([]byte, proto.UDP_MAX_PACKET_SIZE)
		listener.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := listener.ReadFrom(data)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {

			} else {
				g.L.Warn("AgentStat ReadFrom", zap.String("addr", misc.Conf.Pinpoint.StatAddr), zap.String("error", err.Error()))
			}
			continue
		}
		if n == 0 {
			continue
		}

		spans, err := udpRead(data[:n])
		if err != nil {
			g.L.Warn("udpRead", zap.String("error", err.Error()))
			return err
		}
		p.udpChan <- spans
	}
}

// AgentSpan ...
func (p *Pinpoint) AgentSpan() error {
	addrInfo, _ := net.ResolveUDPAddr("udp", misc.Conf.Pinpoint.SpanAddr)
	listener, err := net.ListenUDP("udp", addrInfo)
	if err != nil {
		g.L.Fatal("listen udp", zap.String("addr", misc.Conf.Pinpoint.SpanAddr), zap.String("error", err.Error()))
	}

	for {
		data := make([]byte, proto.UDP_MAX_PACKET_SIZE)
		listener.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := listener.ReadFrom(data)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {

			} else {
				g.L.Warn("readfrom", zap.String("addr", misc.Conf.Pinpoint.SpanAddr), zap.String("error", err.Error()))
			}
			continue
		}
		if n == 0 {
			continue
		}

		spans, err := udpRead(data[:n])
		if err != nil {
			g.L.Warn("udpRead", zap.String("error", err.Error()))
			return err
		}
		p.udpChan <- spans
	}
}

// tcpCollector ...
func (p *Pinpoint) tcpCollector() {

	spanPack := network.NewSpansPacket()
	spanPack.Type = constant.TypeOfTCPData
	spanPack.AppName = gAgent.appName
	spanPack.AgentID = gAgent.agentID

	tracePack := &network.TracePack{
		Type:       constant.TypeOfPinpoint,
		IsSync:     constant.TypeOfSyncNo,
		IsCompress: constant.TypeOfCompressYes,
	}

	for {
		select {
		case span, ok := <-p.tcpChan:
			if !ok {
				break
			}
			spanPack.Payload = append(spanPack.Payload, span)
			payload, err := msgpack.Marshal(spanPack)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				// 清空缓存
				spanPack.Payload = spanPack.Payload[:0]
				break
			}
			tracePack.Payload = payload
			// 发送
			if err := gAgent.collector.write(tracePack); err != nil {
				g.L.Warn("write", zap.String("error", err.Error()))
			}
			// 清空缓存
			spanPack.Payload = spanPack.Payload[:0]
			break
		}
	}
}

// udpCollector ...
func (p *Pinpoint) udpCollector() {
	// 定时器
	ticker := time.NewTicker(time.Duration(misc.Conf.Pinpoint.SpanReportInterval) * time.Millisecond)
	defer ticker.Stop()

	spanPack := network.NewSpansPacket()
	spanPack.Type = constant.TypeOfUDPData
	spanPack.AppName = gAgent.appName
	spanPack.AgentID = gAgent.agentID

	tracePack := &network.TracePack{
		Type:       constant.TypeOfPinpoint,
		IsSync:     constant.TypeOfSyncNo,
		IsCompress: constant.TypeOfCompressYes,
	}

	for {
		select {
		case spanData, ok := <-p.udpChan:
			if ok {
				spanPack.Payload = append(spanPack.Payload, spanData)
				if len(spanPack.Payload) >= misc.Conf.Pinpoint.SpanQueueLen {
					payload, err := msgpack.Marshal(spanPack)
					if err != nil {
						g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
						// 清空缓存
						spanPack.Payload = spanPack.Payload[:0]
						continue
					}
					tracePack.Payload = payload
					// 发送
					if err := gAgent.collector.write(tracePack); err != nil {
						g.L.Warn("write", zap.String("error", err.Error()))
					}
					// 清空缓存
					spanPack.Payload = spanPack.Payload[:0]
				}
			}
			break
		case <-ticker.C:
			if len(spanPack.Payload) > 0 {
				payload, err := msgpack.Marshal(spanPack)
				if err != nil {
					g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
					// 清空缓存
					spanPack.Payload = spanPack.Payload[:0]
					continue
				}
				tracePack.Payload = payload
				// 发送
				if err := gAgent.collector.write(tracePack); err != nil {
					g.L.Warn("write", zap.String("error", err.Error()))
				}
				// 清空缓存
				spanPack.Payload = spanPack.Payload[:0]
			}
		}
	}
}

func (p *Pinpoint) agentInfo(conn net.Conn) error {
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

	isRecvOffline := false
	defer func() {
		if !isRecvOffline {
			// sdk客户端断线
			// gAgent.agentInfo.IsLive = false
			// gAgent.isReportAgentInfo = true
			// gAgent.agentInfo.EndTimestamp = time.Now().UnixNano() / 1e6
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

			spans, err := handletcp(applicationSend.GetPayload())
			if err != nil {
				g.L.Warn("handletcp", zap.String("error", err.Error()))
				return err
			}

			p.tcpChan <- spans
			break

		case proto.APPLICATION_REQUEST:
			// g.L.Debug("agentInfo", zap.String("type", "APPLICATION_REQUEST"))
			applicationRequest := proto.NewApplicationRequest()
			if err := applicationRequest.Decode(conn, reader); err != nil {
				g.L.Warn("decode", zap.String("error", err.Error()))
				return err
			}

			spans, err := handletcp(applicationRequest.GetPayload())
			if err != nil {
				g.L.Warn("handle tcp", zap.String("error", err.Error()))
				return err
			}
			p.tcpChan <- spans

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
				g.L.Warn("response Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE"))
			applicationStreamCreate := proto.NewApplicationStreamCreate()
			if err := applicationStreamCreate.Decode(conn, reader); err != nil {
				g.L.Warn("stream create Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CLOSE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CLOSE"))
			applicationStreamClose := proto.NewApplicationStreamClose()
			if err := applicationStreamClose.Decode(conn, reader); err != nil {
				g.L.Warn("stream close Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE_SUCCESS:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE_SUCCESS"))
			applicationStreamCreateSuccess := proto.NewApplicationStreamCreateSuccess()
			if err := applicationStreamCreateSuccess.Decode(conn, reader); err != nil {
				g.L.Warn("stream create success Decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_CREATE_FAIL:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_CREATE_FAIL"))
			applicationStreamCreateFail := proto.NewApplicationStreamCreateFail()
			if err := applicationStreamCreateFail.Decode(conn, reader); err != nil {
				g.L.Warn("stream create fail Decode", zap.String("error", err.Error()))
				return err
			}

			break

		case proto.APPLICATION_STREAM_RESPONSE:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_RESPONSE"))
			applicationStreamResponse := proto.NewApplicationStreamResponse()
			if err := applicationStreamResponse.Decode(conn, reader); err != nil {
				g.L.Warn("stream response decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_PING:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_PING"))
			applicationStreamPing := proto.NewApplicationStreamPing()
			if err := applicationStreamPing.Decode(conn, reader); err != nil {
				g.L.Warn("stream ping decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.APPLICATION_STREAM_PONG:
			g.L.Debug("agentInfo", zap.String("type", "APPLICATION_STREAM_PONG"))
			applicationStreamPong := proto.NewApplicationStreamPong()
			if err := applicationStreamPong.Decode(conn, reader); err != nil {
				g.L.Warn("stream pong decode", zap.String("error", err.Error()))
				return err
			}
			break

		case proto.CONTROL_HANDSHAKE:
			// 获取并保存Agent信息
			g.L.Debug("agentInfo", zap.String("type", "CONTROL_HANDSHAKE"))
			controlHandShake := proto.NewControlHandShake()
			if err := controlHandShake.Decode(conn, reader); err != nil {
				g.L.Warn("control hand shake decode", zap.String("error", err.Error()))
				return err
			}

			agentInfo := network.NewAgentInfo()
			if err := json.Unmarshal(controlHandShake.GetPayload(), agentInfo); err != nil {
				g.L.Warn("json unmarshal", zap.String("error", err.Error()))
				return err
			}

			agentInfo.IsLive = true
			gAgent.isLive = true

			// 保存App信息
			gAgent.appName = agentInfo.AppName
			gAgent.agentID = agentInfo.AgentID

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
			gAgent.isLive = false
			// gAgent.isReportAgentInfo = true
			// gAgent.agentInfo.EndTimestamp = time.Now().UnixNano() / 1e6
			// isRecvOffline = true

			// controlClientClose := proto.NewControlClientClose()
			// if err := controlClientClose.Decode(conn, reader); err != nil {
			// 	g.L.Warn("controlClientClose.Decode", zap.String("error", err.Error()))
			// }
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
				g.L.Warn("control ping payload decode", zap.String("error", err.Error()))
				return err
			}
			break

		default:
			g.L.Warn("unaware packet Type", zap.Int16("packetType", packetType))
		}

		if isRePacket {
			body, err := rePacket.Encode()
			if err != nil {
				g.L.Warn("rePacket encode", zap.String("error", err.Error()))
				return err
			}

			if _, err := conn.Write(body); err != nil {
				g.L.Warn("write", zap.String("error", err.Error()))
				return err
			}
		}
	}
}

// handletcp ...
func handletcp(message []byte) (*network.Spans, error) {
	spans := network.NewSpans()
	tStruct := thrift.Deserialize(message)
	switch m := tStruct.(type) {
	case *pinpoint.TAgentInfo:
		spans.Type = constant.TypeOfAgentInfo
		spans.Spans = message
		break
	case *trace.TSqlMetaData:
		spans.Type = constant.TypeOfSQLMetaData
		spans.Spans = message
		break
	case *trace.TApiMetaData:
		spans.Type = constant.TypeOfAPIMetaData
		spans.Spans = message
		break
	case *trace.TStringMetaData:
		spans.Type = constant.TypeOfStringMetaData
		spans.Spans = message
		break
	default:
		g.L.Warn("unknown type", zap.String("type", fmt.Sprintf("unknow type [%T]", m)), zap.Any("data", tStruct))
		return nil, fmt.Errorf("unknow type %t", m)
	}
	return spans, nil
}

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
