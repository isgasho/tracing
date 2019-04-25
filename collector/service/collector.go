package service

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/imdevlab/tracing/pkg/alert"

	"github.com/imdevlab/tracing/pkg/mq"

	"go.uber.org/zap"

	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/collector/storage"
	"github.com/imdevlab/tracing/collector/ticker"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/vmihailenco/msgpack"
)

var logger *zap.Logger

// Collector 采集服务
type Collector struct {
	etcd    *Etcd            // 服务上报
	apps    *AppStore        // app集合
	ticker  *ticker.Tickers  // 定时器
	storage *storage.Storage // 存储
	mq      *mq.Nats         // 消息队列
	pushC   chan *alert.Data // 推送通道
}

var gCollector *Collector

// New new collecotr
func New(l *zap.Logger) *Collector {
	logger = l
	gCollector = &Collector{
		etcd:    newEtcd(),
		apps:    newAppStore(),
		ticker:  ticker.NewTickers(misc.Conf.Ticker.Num, misc.Conf.Ticker.Interval, logger),
		storage: storage.NewStorage(logger),
		mq:      mq.NewNats(logger),
		pushC:   make(chan *alert.Data, 3000),
	}
	return gCollector
}

// Start 启动collector
func (c *Collector) Start() error {

	// 启动存储服务
	if err := c.mq.Start(misc.Conf.MQ.Addrs, misc.Conf.MQ.Topic); err != nil {
		logger.Warn("mq start  error", zap.String("error", err.Error()))
		return err
	}

	// 启动存储服务
	if err := c.storage.Start(); err != nil {
		logger.Warn("storage start  error", zap.String("error", err.Error()))
		return err
	}

	// 存储服务类型
	if err := c.storage.StoreSrvType(); err != nil {
		logger.Warn("store server type", zap.String("error", err.Error()))
		return err
	}

	// 存储服务类型
	if err := c.apps.start(); err != nil {
		logger.Warn("apps start error", zap.String("error", err.Error()))
		return err
	}

	// 初始化上报key
	key, err := reportKey(misc.Conf.Etcd.ReportDir)
	if err != nil {
		logger.Warn("get reportKey error", zap.String("error", err.Error()))
		return err
	}

	// 初始化etcd
	if err := c.etcd.Init(misc.Conf.Etcd.Addrs, key, misc.Conf.Collector.Addr); err != nil {
		logger.Warn("etcd init error", zap.String("error", err.Error()))
		return err
	}

	// 启动etcd服务
	if err := c.etcd.Start(); err != nil {
		logger.Warn("etcd start error", zap.String("error", err.Error()))
		return err
	}

	// 启动tcp服务
	if err := c.startNetwork(); err != nil {
		logger.Warn("start network error", zap.String("error", err.Error()))
		return err
	}

	// 启动tcp服务
	if err := c.pushWork(); err != nil {
		logger.Warn("start push work error", zap.String("error", err.Error()))
		return err
	}

	logger.Info("Collector start ok")
	return nil
}

// Close 关闭collector
func (c *Collector) Close() error {
	close(c.pushC)
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
		logger.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
		return err
	}
	switch cmd.Type {
	case constant.TypeOfPing:
		ping := network.NewPing()
		if err := msgpack.Unmarshal(cmd.Payload, ping); err != nil {
			logger.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
			return err
		}
		// logger.Debug("ping", zap.String("addr", conn.RemoteAddr().String()))
	}
	return nil
}

// start 启动tcp服务
func (c *Collector) startNetwork() error {
	lsocket, err := net.Listen("tcp", misc.Conf.Collector.Addr)
	if err != nil {
		logger.Fatal("Listen", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Collector.Addr))
	}

	go func() {
		for {
			conn, err := lsocket.Accept()
			if err != nil {
				logger.Fatal("Accept", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Collector.Addr))
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
	var appname string
	var agentid string
	initName := false

	defer func() {
		if err := recover(); err != nil {
			logger.Error("tcpClient", zap.Any("msg", err))
			return
		}
	}()

	defer func() {
		if conn != nil {
			conn.Close()
		}
		close(quitC)

		if err := gCollector.storage.UpdateAgentState(appname, agentid, false); err != nil {
			logger.Warn("update agent state Store", zap.String("error", err.Error()))
		}
	}()

	go tcpRead(conn, packetC, quitC)

	for {
		select {
		case packet, ok := <-packetC:
			if !ok {
				logger.Info("quit")
				return
			}
			switch packet.Type {
			case constant.TypeOfCmd:
				if err := cmdPacket(conn, packet); err != nil {
					logger.Warn("cmd packet", zap.String("error", err.Error()))
					return
				}
				break
			case constant.TypeOfPinpoint:
				if err := pinpointPacket(conn, packet, &appname, &agentid, &initName); err != nil {
					logger.Warn("pinpoint packet", zap.String("error", err.Error()))
					return
				}
				break
			case constant.TypeOfSystem:
				// if err := v.dealSystem(packet); err != nil {
				// 	logger.Warn("dealSystem", zap.String("error", err.Error()))
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
				logger.Warn("tcp read error", zap.String("err", err.Error()))
				return
			}
			packetC <- packet
			// 设置超时时间
			conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Collector.Timeout) * time.Second))
		}
	}
}

func (c *Collector) pushWork() error {
	for {
		select {
		case packet, ok := <-c.pushC:
			if ok {
				data, err := msgpack.Marshal(packet)
				if err != nil {
					logger.Warn("msgpack", zap.String("error", err.Error()))
					break
				}
				if err := c.mq.Publish(misc.Conf.MQ.Topic, data); err != nil {
					logger.Warn("publish", zap.Error(err))
				}
			}
			break
		}
	}
	// return nil
}

func (c *Collector) publish(data *alert.Data) {
	c.pushC <- data
}

func getblockIndex(value int) int {
	if 0 <= value && value <= 15 {
		return 0
	} else if 16 <= value && value <= 30 {
		return 1
	} else if 31 <= value && value <= 45 {
		return 2
	} else if 46 <= value && value <= 59 {
		return 3
	}
	return 0
}
