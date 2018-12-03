package service

import (
	"github.com/mafanr/vgo/agent/misc"
	"net"
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/pinpoint/proto"
	"go.uber.org/zap"
)

// Pinpoint  analysis pinpoint
type Pinpoint struct {
	//infoAddr string // tcp addr for info
	//statAddr string // udp addr for stat
	//spanAddr string // udp addr for span
}

// NewPinpoint ...
func NewPinpoint() *Pinpoint {
	return &Pinpoint{}
}

// Start ...
func (pinpoint *Pinpoint) Start() error {

	go pinpoint.AgentInfo()
	go pinpoint.AgentStat()
	go pinpoint.AgentSpan()

	return nil
}

// Close ...
func (pinpoint *Pinpoint) Close() error {
	return nil
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
		handleAgentUDP(data[:n])
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
	for {
		data := make([]byte, proto.UDP_MAX_PACKET_SIZE)
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
		handleAgentUDP(data[:n])
		g.L.Debug("AgentSpan Recv", zap.String("message", string(data[:n])))
	}
}
