package service

import (
	"bufio"
	"net"
	"time"

	"github.com/vmihailenco/msgpack"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"github.com/mafanr/vgo/vgo/misc"
	"go.uber.org/zap"
)

// Vgo ...
type Vgo struct {
	storage  *Storage  // 存储
	pinpoint *Pinpoint // 处理pinpoint 数据
	appStore *AppStore
	etcd     *Etcd
}

var gVgo *Vgo

// New ...
func New() *Vgo {
	gVgo = &Vgo{
		storage:  NewStorage(),
		pinpoint: NewPinpoint(),
		appStore: NewAppStore(),
		etcd:     NewEtcd(),
	}
	return gVgo
}

// Start ...
func (v *Vgo) Start() error {

	if err := v.etcd.Start(); err != nil {
		g.L.Fatal("Start:etcd.Start", zap.String("error", err.Error()))
		return err
	}

	if err := v.storage.Init(); err != nil {
		g.L.Fatal("Start:storage.Start", zap.String("error", err.Error()))
		return err
	}

	if err := v.storage.Start(); err != nil {
		g.L.Fatal("Start:storage.Start", zap.String("error", err.Error()))
		return err
	}

	if err := v.init(); err != nil {
		g.L.Fatal("Start:v.init", zap.String("error", err.Error()))
		return err
	}

	return nil
}

func (v *Vgo) init() error {
	// init mysql
	if err := v.initMysql(); err != nil {
		g.L.Warn("init:v.initMysql", zap.String("error", err.Error()))
		return err
	}

	// init service
	v.acceptAgent()

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
		g.L.Fatal("acceptAgent:net.Listen", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Vgo.ListenAddr))
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				g.L.Fatal("acceptAgent:ln.Accept", zap.String("msg", err.Error()), zap.String("addr", misc.Conf.Vgo.ListenAddr))
			}
			conn.SetReadDeadline(time.Now().Add(time.Duration(misc.Conf.Vgo.AgentTimeout) * time.Second))
			go v.agentWork(conn)
		}
	}()

	return nil
}

func (v *Vgo) agentWork(conn net.Conn) {
	quitC := make(chan bool, 1)
	packetC := make(chan *util.VgoPacket, 100)

	defer func() {
		if err := recover(); err != nil {
			g.L.Error("agentWork:.", zap.Any("msg", err))
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
						g.L.Warn("agentWork:v.dealCmd", zap.String("error", err.Error()))
						return
					}
					break
				case util.TypeOfPinpoint:
					if err := v.pinpoint.dealUpload(conn, packet); err != nil {
						g.L.Warn("agentWork:v.pinpoint.dealUpload", zap.String("error", err.Error()))
						return
					}
					break
				}
			}
		}
	}
}

func (v *Vgo) dealCmd(conn net.Conn, packet *util.VgoPacket) error {
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

func (v *Vgo) agentRead(conn net.Conn, packetC chan *util.VgoPacket, quitC chan bool) {
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
		packet := util.NewVgoPacket()
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
