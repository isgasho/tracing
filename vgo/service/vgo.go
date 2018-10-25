package service

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/vmihailenco/msgpack"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"github.com/mafanr/vgo/vgo/misc"
	"github.com/mafanr/vgo/vgo/stats"
	"go.uber.org/zap"
)

// Vgo ...
type Vgo struct {
	stats   *stats.Stats
	storage *Storage
	apps    *AppStore
}

// New ...
func New() *Vgo {
	return &Vgo{
		stats:   stats.New(),
		storage: NewStorage(),
		apps:    NewAppStore(),
	}
}

// Start ...
func (v *Vgo) Start() error {

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

	// load apps
	if err := v.apps.LoadApps(); err != nil {
		g.L.Warn("init:apps.LoadApps", zap.String("error", err.Error()))
		return err
	}

	// load server name code
	if err := v.apps.LoadSerCode(); err != nil {
		g.L.Warn("init:apps.LoadSerCode", zap.String("error", err.Error()))
		return err
	}

	// start web ser

	// start stats
	if err := v.stats.Start(); err != nil {
		g.L.Warn("init:v.stats.Start", zap.String("error", err.Error()))
		return err
	}

	// init service
	v.acceptAgent()

	return nil
}

func (v *Vgo) initMysql() error {
	// init sql
	g.InitMysql(misc.Conf.Mysql.Acc, misc.Conf.Mysql.Pw, misc.Conf.Mysql.Addr, misc.Conf.Mysql.Port, misc.Conf.Mysql.Database)
	return nil
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
			g.L.Info("Quit")
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
				case util.TypeOfSkywalking:
					if err := v.dealSkywalking(conn, packet); err != nil {
						g.L.Warn("agentWork:v.dealSkywalking", zap.String("error", err.Error()))
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

// dealSkywalking skywlking报文处理
func (v *Vgo) dealSkywalking(conn net.Conn, packet *util.VgoPacket) error {
	skypacker := &util.SkywalkingPacket{}
	if err := msgpack.Unmarshal(packet.Payload, skypacker); err != nil {
		g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
		return err
	}
	switch skypacker.Type {
	case util.TypeOfAppRegister:
		appRegister := &util.AppRegister{}
		if err := msgpack.Unmarshal(skypacker.Payload, appRegister); err != nil {
			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
			return err
		}

		code, err := v.apps.LoadAppCode(appRegister.Name)
		if err != nil {
			g.L.Warn("dealSkywalking:v.apps.LoadAppCode", zap.String("name", appRegister.Name), zap.String("error", err.Error()))
			return err
		}

		appRegister.Code = code

		mbuf, err := msgpack.Marshal(appRegister)
		if err != nil {
			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", appRegister.Name), zap.String("error", err.Error()))
			return err
		}
		skypacker.Payload = mbuf

		payload, err := msgpack.Marshal(skypacker)
		if err != nil {
			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", appRegister.Name), zap.String("error", err.Error()))
			return err
		}

		packet.Payload = payload
		if _, err := conn.Write(packet.Encode()); err != nil {
			g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
			return err
		}
		break
	case util.TypeOfAppRegisterInstance:
		agentInfo := &util.AgentInfo{}
		if err := msgpack.Unmarshal(skypacker.Payload, agentInfo); err != nil {
			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
			return err
		}

		id, err := v.apps.LoadAgentID(agentInfo)
		if err != nil {
			g.L.Warn("dealSkywalking:v.apps.LoadAppCode", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
			return err
		}

		appRegisterIns := &util.KeyWithIntegerValue{
			Key:   "id",
			Value: id,
		}

		mbuf, err := msgpack.Marshal(appRegisterIns)
		if err != nil {
			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
			return err
		}
		skypacker.Payload = mbuf

		payload, err := msgpack.Marshal(skypacker)
		if err != nil {
			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
			return err
		}

		packet.Payload = payload

		if _, err := conn.Write(packet.Encode()); err != nil {
			g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
			return err
		}
		break
	case util.TypeOfSerNameDiscoveryService:
		serNames := &util.SerNameDiscoveryServices{}
		if err := msgpack.Unmarshal(skypacker.Payload, serNames); err != nil {
			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
			return err
		}

		appName, ok := v.apps.LoadAppName(serNames.AppCode)
		if !ok {
			g.L.Warn("dealSkywalking:v.apps.LoadAppName", zap.String("error", "unfind app name"), zap.Int32("appCode", serNames.AppCode))
			return fmt.Errorf("unfind app name, app code is %d", serNames.AppCode)
		}

		log.Println(appName)
		// id, err := v.apps.LoadAgentID(agentInfo)
		// if err != nil {
		// 	g.L.Warn("dealSkywalking:v.apps.LoadAppCode", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
		// 	return err
		// }

		// appRegisterIns := &util.KeyWithIntegerValue{
		// 	Key:   "id",
		// 	Value: id,
		// }

		// mbuf, err := msgpack.Marshal(appRegisterIns)
		// if err != nil {
		// 	g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
		// 	return err
		// }
		// skypacker.Payload = mbuf

		// payload, err := msgpack.Marshal(skypacker)
		// if err != nil {
		// 	g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
		// 	return err
		// }

		// packet.Payload = payload

		// if _, err := conn.Write(packet.Encode()); err != nil {
		// 	g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
		// 	return err
		// }
		break
	}
	return nil
}
