package service

import (
	"net"
	"time"

	"github.com/mafanr/vgo/agent/misc"
	"github.com/mafanr/vgo/util"
	"github.com/vmihailenco/msgpack"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/proto/pinpoint/proto"
	"go.uber.org/zap"
)

// Pinpoint  analysis pinpoint
type Pinpoint struct {
	tcpChan chan []byte
	udpChan chan *util.SpanDataModel
}

// NewPinpoint ...
func NewPinpoint() *Pinpoint {
	return &Pinpoint{
		tcpChan: make(chan []byte, 100),
		udpChan: make(chan *util.SpanDataModel, 300),
	}
}

// Start ...
func (pinpoint *Pinpoint) Start() error {

	go pinpoint.AgentInfo()
	go pinpoint.AgentStat()
	go pinpoint.AgentSpan()

	go pinpoint.tcpCollector()
	go pinpoint.udpCollector()

	return nil
}

// Close ...
func (pinpoint *Pinpoint) Close() error {
	return nil
}

// Start ...
func (pinpoint *Pinpoint) tcpCollector() {
	for {
		select {
		case data, ok := <-pinpoint.tcpChan:
			if ok {
				packet := &util.VgoPacket{
					Type:       util.TypeOfPinpoint,
					Version:    util.VersionOf01,
					IsSync:     util.TypeOfSyncNo,
					IsCompress: util.TypeOfCompressNo,
					Payload:    data,
				}
				if err := gAgent.client.WritePacket(packet); err != nil {
					g.L.Warn("sendSpans:gAgent.client.WritePacket", zap.String("error", err.Error()))
				}
			}
		}
	}
}

// AgentInfo ...
func (pinpoint *Pinpoint) AgentInfo() error {
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
		go pinpoint.agentInfo(conn)
	}
}

// AgentStat ...
func (pinpoint *Pinpoint) AgentStat() error {
	addrInfo, _ := net.ResolveUDPAddr("udp", misc.Conf.Pinpoint.StatAddr)
	listener, err := net.ListenUDP("udp", addrInfo)
	if err != nil {
		g.L.Fatal("AgentStat ListenUDP", zap.String("addr", misc.Conf.Pinpoint.StatAddr), zap.String("error", err.Error()))
	}
	data := make([]byte, proto.UDP_MAX_PACKET_SIZE)
	for {
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
		//handleAgentUDP(data[:n])
		spanModel, err := handleAgentUDP(data[:n])
		if err != nil {
			g.L.Warn("AgentStat:handleAgentUDP", zap.String("error", err.Error()))
			return err
		}
		pinpoint.udpChan <- spanModel
		g.L.Debug("AgentStat Recv", zap.String("message", string(data[:n])))
	}
}

// AgentSpan ...
func (pinpoint *Pinpoint) AgentSpan() error {
	addrInfo, _ := net.ResolveUDPAddr("udp", misc.Conf.Pinpoint.SpanAddr)
	listener, err := net.ListenUDP("udp", addrInfo)
	if err != nil {
		g.L.Fatal("AgentSpan ListenUDP", zap.String("addr", misc.Conf.Pinpoint.SpanAddr), zap.String("error", err.Error()))
	}
	data := make([]byte, proto.UDP_MAX_PACKET_SIZE)
	for {
		listener.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := listener.ReadFrom(data)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {

			} else {
				g.L.Warn("AgentSpan ReadFrom", zap.String("addr", misc.Conf.Pinpoint.SpanAddr), zap.String("error", err.Error()))
			}
			continue
		}
		if n == 0 {
			continue
		}

		spanModel, err := handleAgentUDP(data[:n])
		if err != nil {
			g.L.Warn("AgentSpan:handleAgentUDP", zap.String("error", err.Error()))
			return err
		}
		pinpoint.udpChan <- spanModel
		g.L.Debug("AgentSpan Recv", zap.String("message", string(data[:n])))
	}
}

// udpCollector ...
func (pinpoint *Pinpoint) udpCollector() {
	// 定时器
	ticker := time.NewTicker(time.Duration(misc.Conf.Pinpoint.SpanReportInterval) * time.Millisecond)
	pinpointData := util.NewPinpointData()
	pinpointData.Type = util.TypeOfUDPData
	pinpointData.AgentName = gAgent.appName
	pinpointData.AgentID = gAgent.agentID

	packet := &util.VgoPacket{
		Type:       util.TypeOfPinpoint,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncNo,
		IsCompress: util.TypeOfCompressYes,
	}

	for {
		select {
		case spanData, ok := <-pinpoint.udpChan:
			if ok {
				pinpointData.Payload = append(pinpointData.Payload, spanData)
				if len(pinpointData.Payload) >= misc.Conf.Pinpoint.SpanQueueLen {
					payload, err := msgpack.Marshal(pinpointData)
					if err != nil {
						g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
						// 清空缓存
						pinpointData.Payload = pinpointData.Payload[:0]
						continue
					}
					packet.Payload = payload
					// 发送
					if err := gAgent.client.WritePacket(packet); err != nil {
						g.L.Warn("sendSpans:gAgent.client.WritePacket", zap.String("error", err.Error()))
					}
					// 清空缓存
					pinpointData.Payload = pinpointData.Payload[:0]
				}
				// g.L.Debug("udpCollector", zap.Any("packet", packet), zap.Any("data", spanData))
			}
			break
		case <-ticker.C:
			if len(pinpointData.Payload) > 0 {
				payload, err := msgpack.Marshal(pinpointData)
				if err != nil {
					g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
					// 清空缓存
					pinpointData.Payload = pinpointData.Payload[:0]
					continue
				}
				packet.Payload = payload
				// 发送
				if err := gAgent.client.WritePacket(packet); err != nil {
					g.L.Warn("sendSpans:gAgent.client.WritePacket", zap.String("error", err.Error()))
				}
				// 清空缓存
				pinpointData.Payload = pinpointData.Payload[:0]
			}
		}
	}
}

// reportSEND ...
func (pinpoint *Pinpoint) reportSEND(data []byte) error {
	pinpointData := util.NewPinpointData()
	pinpointData.Type = util.TypeOfTCPData
	pinpointData.AgentName = gAgent.appName
	pinpointData.AgentID = gAgent.agentID

	spanData := &util.SpanDataModel{
		Type:  util.TypeOfAgentSEND,
		Spans: data,
	}
	pinpointData.Payload = append(pinpointData.Payload, spanData)
	payload, err := msgpack.Marshal(pinpointData)
	if err != nil {
		g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	gAgent.pinpoint.tcpChan <- payload

	return nil
}

// reportAgentInfo ...
func (pinpoint *Pinpoint) reportAgentInfo(appInfo *util.AgentInfo) error {
	pinpointData := util.NewPinpointData()
	pinpointData.Type = util.TypeOfTCPData
	pinpointData.AgentName = appInfo.AppName
	pinpointData.AgentID = appInfo.AgentID

	appInfoBuf, err := msgpack.Marshal(appInfo)
	if err != nil {
		g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}
	spanData := &util.SpanDataModel{
		Type:  util.TypeOfAgentInfo,
		Spans: appInfoBuf,
	}
	pinpointData.Payload = append(pinpointData.Payload, spanData)
	payload, err := msgpack.Marshal(pinpointData)
	if err != nil {
		g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	gAgent.pinpoint.tcpChan <- payload

	return nil
}
