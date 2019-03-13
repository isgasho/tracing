package service

import (
	"bufio"
	"io"
	"net"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/agent/misc"
	"github.com/imdevlab/tracing/util"
	"github.com/vmihailenco/msgpack"
	"go.uber.org/zap"
)

// TCPClient tcp client
type TCPClient struct {
	conn net.Conn
}

// NewTCPClient ...
func NewTCPClient() *TCPClient {
	return &TCPClient{}
}

// Init ...
func (t *TCPClient) Init() error {
	var err error

	isRestart := true
	quitC := make(chan bool, 1)
	// 定时器
	keepLiveTc := time.NewTicker(time.Duration(misc.Conf.Agent.KeepLiveInterval) * time.Second)

	defer func() {
		if err := recover(); err != nil {
			g.L.Warn("Init:.", zap.Stack("server"), zap.Any("err", err))
		}
		// 是否重启
		if isRestart {
			t.Init()
		}
	}()

	defer func() {
		close(quitC)
		t.conn.Close()
		keepLiveTc.Stop()
	}()

	// connect tracing
	for {
		t.conn, err = net.Dial("tcp", misc.Conf.Agent.TracingAddr)
		if err != nil {
			g.L.Warn("Init:net.Dial", zap.String("err", err.Error()), zap.String("addr", misc.Conf.Agent.TracingAddr))
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// 启动心跳
	go func() {
		for {
			select {
			case <-keepLiveTc.C:
				if err := t.KeepLive(); err != nil {
					g.L.Warn("Init:t.KeepLive", zap.String("error", err.Error()))
				}
				break
			case <-quitC:
				return
			}
		}
	}()
	reader := bufio.NewReaderSize(t.conn, util.MaxMessageSize)
	for {
		packet, err := t.ReadPacket(reader)
		if err != nil {
			g.L.Warn("Init:t.ReadPacket", zap.Error(err))
			return err
		}

		// g.L.Debug("Init:t.ReadPacket", zap.Any("packet", packet))

		// 发给上层处理
		switch packet.IsSync {
		case util.TypeOfSyncYes:
			if err := gAgent.syncCall.syncWrite(packet.ID, packet); err != nil {
				g.L.Warn("Init:gAgent.syncCall.syncWrite", zap.Error(err))
			}
			break
		default:
			gAgent.downloadC <- packet
			break
		}

	}
}

// KeepLive ...
func (t *TCPClient) KeepLive() error {

	ping := util.NewPing()
	b, err := msgpack.Marshal(ping)
	if err != nil {
		g.L.Warn("KeepLive:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	cmd := util.NewCMD()
	cmd.Type = util.TypeOfPing
	cmd.Payload = b

	buf, err := msgpack.Marshal(cmd)
	if err != nil {
		g.L.Warn("KeepLive:msgpack.Marshal", zap.String("error", err.Error()))
		return err
	}

	// packet := util.NewTracingPacket(util.TypeOfCmd, util.VersionOf01, util.TypeOfSyncNo, util.TypeOfCompressNo, 0, buf)
	packet := &util.TracingPacket{
		Type:       util.TypeOfCmd,
		Version:    util.VersionOf01,
		IsSync:     util.TypeOfSyncNo,
		IsCompress: util.TypeOfCompressNo,
		ID:         0,
		Payload:    buf,
	}

	if err := t.WritePacket(packet); err != nil {
		g.L.Warn("KeepLive:t.WritePacket", zap.String("error", err.Error()))
		return err
	}

	return nil
}

// ReadPacket ...
func (t *TCPClient) ReadPacket(rdr io.Reader) (*util.TracingPacket, error) {
	packet := &util.TracingPacket{}
	if err := packet.Decode(rdr); err != nil {
		g.L.Warn("ReadPacket:packet.Decode", zap.String("error", err.Error()))
		return nil, err
	}
	return packet, nil
}

// WritePacket ...
func (t *TCPClient) WritePacket(packet *util.TracingPacket) error {
	body := packet.Encode()
	if t.conn != nil {
		_, err := t.conn.Write(body)
		if err != nil {
			g.L.Warn("WritePacket:t.conn.Write", zap.String("error", err.Error()))
			return err
		}
	}
	return nil
}

// Close ....
func (t *TCPClient) Close() error {
	if t.conn != nil {
		t.conn.Close()
	}
	return nil
}
