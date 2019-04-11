package service

import (
	"fmt"
	"net"
	"os"

	"go.uber.org/zap"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/collector/storage"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/imdevlab/tracing/pkg/ttype"
	"github.com/vmihailenco/msgpack"
)

// Collector 采集服务
type Collector struct {
	etcd      *Etcd            // 服务上报
	apps      *Apps            // app集合
	tickers   *Tickers         // 定时器
	storage   *storage.Storage // 存储
	tcpServer *tcpServer
}

var gCollector *Collector

// New new collecotr
func New() *Collector {
	gCollector = &Collector{
		etcd:      newEtcd(),
		apps:      newApps(),
		tickers:   newTickers(10),
		storage:   storage.NewStorage(),
		tcpServer: newtcpServer(),
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
	if err := c.tcpServer.start(); err != nil {
		g.L.Warn("tcp server start", zap.String("error", err.Error()))
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
	case ttype.TypeOfPing:
		ping := network.NewPing()
		if err := msgpack.Unmarshal(cmd.Payload, ping); err != nil {
			g.L.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
			return err
		}
		// g.L.Debug("ping", zap.String("addr", conn.RemoteAddr().String()))
	}
	return nil
}
