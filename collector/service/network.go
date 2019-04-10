package service

import (
	"bufio"
	"log"
	"net"
	"time"

	"github.com/shaocongcong/tracing/pkg/proto/network"
	"github.com/shaocongcong/tracing/pkg/proto/ttype"

	"github.com/imdevlab/g"
	"github.com/shaocongcong/tracing/collector/misc"
	"go.uber.org/zap"
)

// tcpServer tcp服务端
type tcpServer struct {
}

func newtcpServer() *tcpServer {
	return &tcpServer{}
}

// start 启动tcp服务
func (t *tcpServer) start() error {
	lsocket, err := net.Listen("tcp", misc.Conf.Collector.Addr)
	if err != nil {
		g.L.Fatal("Listen", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Collector.Addr))
	}

	go func() {
		for {
			conn, err := lsocket.Accept()
			if err != nil {
				g.L.Fatal("Accept", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Collector.Addr))
			}
			conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Collector.Timeout) * time.Second))
			go tcpClient(conn)
		}
	}()
	return nil
}

// close 关闭tcp服务
func (t *tcpServer) close() error {
	return nil
}

func tcpClient(conn net.Conn) {
	quitC := make(chan bool, 1)
	packetC := make(chan *network.TracePack, 100)

	defer func() {
		if err := recover(); err != nil {
			g.L.Error("tcpClient", zap.Any("msg", err))
			return
		}
	}()

	defer func() {
		if conn != nil {
			conn.Close()
		}
		close(quitC)
	}()

	go tcpRead(conn, packetC, quitC)

	for {
		select {
		case packet, ok := <-packetC:
			if !ok {
				return
			}
			switch packet.Type {
			case ttype.TypeOfCmd:
				if err := cmdPacket(conn, packet); err != nil {
					g.L.Warn("cmd packet", zap.String("error", err.Error()))
					return
				}
				break
			case ttype.TypeOfPinpoint:
				if err := pinpointPacket(conn, packet); err != nil {
					g.L.Warn("pinpoint packet", zap.String("error", err.Error()))
					return
				}
				break
			case ttype.TypeOfSystem:
				// if err := v.dealSystem(packet); err != nil {
				// 	g.L.Warn("dealSystem", zap.String("error", err.Error()))
				// 	return
				// }
				log.Println("TypeOfSystem")
				break
			}
		}
	}

}

func tcpRead(conn net.Conn, packetC chan *network.TracePack, quitC chan bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	defer func() {
		close(packetC)
	}()
	reader := bufio.NewReaderSize(conn, ttype.MaxMessageSize)
	for {

		select {
		case <-quitC:
			break
		default:
			packet := network.NewTracePack()
			if err := packet.Decode(reader); err != nil {
				g.L.Warn("agentRead:msg.Decode", zap.String("err", err.Error()))
				return
			}
			packetC <- packet
			// 设置超时时间
			conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Collector.Timeout) * time.Second))
		}
	}
}
