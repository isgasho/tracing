package service

import (
	"bufio"
	"net"
	"time"

	"github.com/vmihailenco/msgpack"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/tracing/misc"
	"github.com/imdevlab/tracing/util"
	"go.uber.org/zap"
)

// Vgo ...
type Vgo struct {
	storage      *Storage  // 存储
	pinpoint     *Pinpoint // 处理pinpoint 数据
	appStore     *AppStore
	srvDiscovery SrvDiscovery
}

var gVgo *Vgo

// New ...
func New() *Vgo {
	gVgo = &Vgo{
		storage:      NewStorage(),
		pinpoint:     NewPinpoint(),
		appStore:     NewAppStore(),
		srvDiscovery: newEtcd(),
	}
	return gVgo
}

// Start ...
func (v *Vgo) Start() error {

	if err := v.srvDiscovery.Start(); err != nil {
		g.L.Fatal("etcd start", zap.String("error", err.Error()))
		return err
	}

	if err := v.storage.Init(); err != nil {
		g.L.Fatal("storage init", zap.String("error", err.Error()))
		return err
	}

	if err := v.storage.Start(); err != nil {
		g.L.Fatal("storage start", zap.String("error", err.Error()))
		return err
	}

	if err := v.init(); err != nil {
		g.L.Fatal("init", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (v *Vgo) init() error {
	// init mysql
	if err := v.initMysql(); err != nil {
		g.L.Warn("init mysql", zap.String("error", err.Error()))
		return err
	}

	v.initServiceType()

	// init service
	v.acceptAgent()

	return nil
}

func (v *Vgo) initServiceType() error {
	if err := v.storage.storeServiceType(); err != nil {
		g.L.Warn("initServiceType err", zap.Error(err))
		return err
	}
	return nil
}

func (v *Vgo) initMysql() error {
	return nil
	// init sql
	// g.InitMysql(misc.Conf.Mysql.Acc, misc.Conf.Mysql.Pw, misc.Conf.Mysql.Addr, misc.Conf.Mysql.Port, misc.Conf.Mysql.Database)
	// return nil
}

func (v *Vgo) acceptAgent() error {
	ln, err := net.Listen("tcp", misc.Conf.Vgo.ListenAddr)
	if err != nil {
		g.L.Fatal("Listen", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Vgo.ListenAddr))
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				g.L.Fatal("Accept", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Vgo.ListenAddr))
			}
			conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Vgo.AgentTimeout) * time.Second))
			go v.agentWork(conn)
		}
	}()

	return nil
}

func (v *Vgo) agentWork(conn net.Conn) {
	quitC := make(chan bool, 1)
	packetC := make(chan *util.TracingPacket, 100)

	defer func() {
		if err := recover(); err != nil {
			g.L.Error("agentWork", zap.Any("msg", err))
			return
		}
	}()

	defer func() {
		close(quitC)
		close(packetC)
		conn.Close()
	}()

	go v.agentRead(conn, packetC, quitC)

	for {
		select {
		case <-quitC:
			g.L.Debug("Quit")
			return
		case packet, ok := <-packetC:
			if ok {
				switch packet.Type {
				case util.TypeOfCmd:
					if err := v.dealCmd(conn, packet); err != nil {
						g.L.Warn("dealCmd", zap.String("error", err.Error()))
						return
					}
					break
				case util.TypeOfPinpoint:
					if err := v.pinpoint.dealUpload(conn, packet); err != nil {
						g.L.Warn("dealUpload", zap.String("error", err.Error()))
						return
					}
					break
				case util.TypeOfSystem:
					if err := v.dealSystem(packet); err != nil {
						g.L.Warn("dealSystem", zap.String("error", err.Error()))
						return
					}
					break
				}
			}
		}
	}
}

func (v *Vgo) dealCmd(conn net.Conn, packet *util.TracingPacket) error {
	cmd := util.NewCMD()
	if err := msgpack.Unmarshal(packet.Payload, cmd); err != nil {
		g.L.Warn("dealCmd:msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}
	switch cmd.Type {
	case util.TypeOfPing:
		ping := util.NewPing()
		if err := msgpack.Unmarshal(cmd.Payload, ping); err != nil {
			g.L.Warn("dealCmd:msgpack.Unmarshal", zap.String("error", err.Error()))
			return err
		}
		g.L.Debug("dealCmd:ping", zap.String("addr", conn.RemoteAddr().String()))
	}
	return nil
}

func (v *Vgo) dealSystem(packet *util.TracingPacket) error {
	metric := util.NewMetricData()
	if err := msgpack.Unmarshal(packet.Payload, metric); err != nil {
		g.L.Warn("msgpack Unmarshal", zap.String("error", err.Error()))
		return err
	}
	v.storage.metricsChan <- metric
	return nil
}

func (v *Vgo) agentRead(conn net.Conn, packetC chan *util.TracingPacket, quitC chan bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	defer func() {
		quitC <- true
	}()
	reader := bufio.NewReaderSize(conn, util.MaxMessageSize)
	for {
		packet := util.NewTracingPacket()
		if err := packet.Decode(reader); err != nil {
			g.L.Warn("agentRead:msg.Decode", zap.String("err", err.Error()))
			return
		}
		packetC <- packet
		// 设置超时时间
		conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Vgo.AgentTimeout) * time.Second))
	}
}

// Close ...
func (v *Vgo) Close() error {

	// 关闭存储
	if err := v.storage.Close(); err != nil {
		g.L.Warn("Close:v.storage.Close", zap.String("error", err.Error()))
	}

	return nil
}
