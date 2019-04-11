package service

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/collector/storage"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/vmihailenco/msgpack"
)

// Collector 采集服务
type Collector struct {
	etcd    *Etcd            // 服务上报
	apps    *Apps            // app集合
	tickers *Tickers         // 定时器
	storage *storage.Storage // 存储
}

var gCollector *Collector

// New new collecotr
func New() *Collector {
	gCollector = &Collector{
		etcd:    newEtcd(),
		apps:    newApps(),
		tickers: newTickers(10),
		storage: storage.NewStorage(),
	}
	return gCollector
}

// Start 启动collector
func (c *Collector) Start() error {

	// 启动存储服务
	if err := c.storage.Start(); err != nil {
		g.L.Warn("storage start", zap.String("error", err.Error()))
		return err
	}

	// 初始化上报key
	key, err := reportKey(misc.Conf.Etcd.ReportDir)
	if err != nil {
		g.L.Warn("get reportKey ", zap.String("error", err.Error()))
		return err
	}

	// 初始化etcd
	if err := c.etcd.Init(misc.Conf.Etcd.Addrs, key, misc.Conf.Collector.Addr); err != nil {
		g.L.Warn("etcd init", zap.String("error", err.Error()))
		return err
	}

	// 启动etcd服务
	if err := c.etcd.Start(); err != nil {
		g.L.Warn("etcd start", zap.String("error", err.Error()))
		return err
	}

	// 启动tcp服务
	if err := c.startNetwork(); err != nil {
		g.L.Warn("start network", zap.String("error", err.Error()))
		return err
	}

	g.L.Info("Collector start ok")
	return nil
}

// Close 关闭collector
func (c *Collector) Close() error {
	return nil
}

func reportKey(dir string) (string, error) {
	value, err := collectorName()
	if err != nil {
		return "", err
	}

	dirLen := len(dir)
	if dirLen > 0 && dir[dirLen-1] != '/' {
		return dir + "/" + value, nil
	}
	return dir + value, nil
}

// collectorName etcd 上报key
func collectorName() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%d", host, os.Getpid()), nil
}

func initDir(dir string) string {
	dirLen := len(dir)
	if dirLen > 0 && dir[dirLen-1] != '/' {
		return dir + "/"
	}
	return dir
}

// cmdPacket 处理agent 发送来的cmd报文
func cmdPacket(conn net.Conn, packet *network.TracePack) error {
	cmd := network.NewCMD()
	if err := msgpack.Unmarshal(packet.Payload, cmd); err != nil {
		g.L.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
		return err
	}
	switch cmd.Type {
	case constant.TypeOfPing:
		ping := network.NewPing()
		if err := msgpack.Unmarshal(cmd.Payload, ping); err != nil {
			g.L.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
			return err
		}
		// g.L.Debug("ping", zap.String("addr", conn.RemoteAddr().String()))
	}
	return nil
}

// start 启动tcp服务
func (c *Collector) startNetwork() error {
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
			case constant.TypeOfCmd:
				if err := cmdPacket(conn, packet); err != nil {
					g.L.Warn("cmd packet", zap.String("error", err.Error()))
					return
				}
				break
			case constant.TypeOfPinpoint:
				if err := pinpointPacket(conn, packet); err != nil {
					g.L.Warn("pinpoint packet", zap.String("error", err.Error()))
					return
				}
				break
			case constant.TypeOfSystem:
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
	reader := bufio.NewReaderSize(conn, constant.MaxMessageSize)
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
