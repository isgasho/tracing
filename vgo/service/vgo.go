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

	"github.com/mafanr/vgo/vgo/web"
	"go.uber.org/zap"
)

// Vgo ...
type Vgo struct {
	storage  *Storage  // 存储
	pinpoint *Pinpoint // 处理pinpoint 数据
	web      *web.Web
	appStore *AppStore
}

var gVgo *Vgo

// New ...
func New() *Vgo {
	gVgo = &Vgo{
		storage:  NewStorage(),
		pinpoint: NewPinpoint(),
		web:      web.New(),
		appStore: NewAppStore(),
	}
	return gVgo
}

// Start ...
func (v *Vgo) Start() error {

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

	// start web ser
	if err := v.web.Start(); err != nil {
		g.L.Warn("init:v.web.Start", zap.String("error", err.Error()))
		return err
	}

	// init service
	v.acceptAgent()

	return nil
}

func (v *Vgo) initMysql() error {
	return nil
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

// LoadApps 加载数据库中的所有app
func (v *Vgo) LoadApps() error {
	//// 加载所有appCode
	//apps := make([]*util.App, 0)
	//if err := g.DB.Select(&apps, "select * from app"); err != nil {
	//	g.L.Fatal("LoadApps:g.DB.Select", zap.Error(err))
	//}
	//
	//for _, app := range apps {
	//	v.apps.Store(app.AppID, app)
	//	// 维护AppName和AppID映射关系
	//	v.appN2c.Store(app.Name, app.AppID)
	//}
	//
	//v.apps.Range(func(key, value interface{}) bool {
	//	g.L.Debug("LoadApps ---- 应用", zap.Any("appID", key), zap.Any("app", value))
	//	return true
	//})
	//
	//v.appN2c.Range(func(key, value interface{}) bool {
	//	g.L.Debug("LoadApps 应用 ID", zap.Any("appID", key), zap.Any("app", value))
	//	return true
	//})
	return nil
}

// LoadAgents 加载数据库中的所有agent
func (v *Vgo) LoadAgents() error {
	//// 加载所有appCode
	//agents := make([]*util.AgentInfo, 0)
	//if err := g.DB.Select(&agents, "select * from agent"); err != nil {
	//	g.L.Fatal("LoadAgents:g.DB.Select", zap.Error(err))
	//}
	//
	//for _, agent := range agents {
	//	app, ok := v.apps.Load(agent.AppID)
	//	if ok {
	//		app.(*util.App).Agents.Store(agent.InstanceID, agent)
	//	}
	//}

	return nil
}

// GetAppID 获取AppID
func (v *Vgo) GetAppID(name string) (int32, error) {
	// // 内存查找
	// id, ok := v.appN2c.Load(name)
	// if ok {
	// 	return id.(int32), nil
	// }

	// // 数据库中查找
	// query := fmt.Sprintf("SELECT id FROM app WHERE `app`.`name`='%s';", name)
	// rows, err := g.DB.Query(query)
	// if err != nil {
	// 	g.L.Warn("GetAppID:g.DB.Exec", zap.Error(err), zap.String("sql", query))
	// 	return 0, err
	// }

	// defer rows.Close()
	// isFind := false
	// var appID int32
	// for rows.Next() {
	// 	rows.Scan(&appID)
	// 	isFind = true
	// 	break
	// }

	// if isFind {
	// 	app := util.NewApp()
	// 	app.Name = name
	// 	app.AppID = int32(appID)
	// 	// 缓存到内存中
	// 	v.apps.Store(int32(appID), app)
	// 	v.apps.Store(app.Name, int32(appID))
	// 	return int32(appID), nil
	// }

	// // 如果不存在插入
	// query = fmt.Sprintf("insert into `app` (`name`) values ('%s')", name)
	// result, err := g.DB.Exec(query)
	// if err != nil {
	// 	g.L.Warn("GetAppID:g.DB.Exec", zap.Error(err), zap.String("sql", query))
	// 	return 0, err
	// }

	// newID, err := result.LastInsertId()
	// if err != nil {
	// 	g.L.Warn("LoadAppCode:result.LastInsertId", zap.Error(err))
	// 	return 0, err
	// }

	// app := util.NewApp()
	// app.Name = name
	// app.AppID = int32(newID)
	// // 缓存到内存中
	// v.apps.Store(int32(newID), app)
	// v.apps.Store(app.Name, int32(newID))
	// return int32(newID), nil
	return 0, nil

}

// loadApp 通过Appid到数库中加载app
func (v *Vgo) loadApp(appid int32) (*util.App, error) {
	// oApp, ok := v.apps.Load(appid)
	// if ok {
	// 	return oApp.(*util.App), nil
	// }
	// query := fmt.Sprintf("select name from  `app` where id='%d'", appid)
	// rows, err := g.DB.Query(query)
	// if err != nil {
	// 	g.L.Warn("loadApp:g.DB.Query", zap.Error(err), zap.Int32("appid", appid))
	// 	return nil, err
	// }
	// // 防止泄漏
	// defer rows.Close()

	// var name string
	// isFind := false
	// for rows.Next() {
	// 	rows.Scan(&name)
	// 	isFind = true
	// 	break
	// }
	// // 数据库中可能不存在
	// if !isFind {
	// 	return nil, fmt.Errorf("unfind app, appid is %d", appid)
	// }
	// app := &util.App{
	// 	Name:  name,
	// 	AppID: appid,
	// }

	// // 缓存到内存
	// v.apps.Store(app.AppID, app)
	// v.appN2c.Store(app.Name, app.AppID)

	// return app, nil

	return nil, nil
}

// GetInstanceID 获取App实例ID
func (v *Vgo) GetInstanceID(agent *util.AgentInfo) (int32, error) {

	//app, err := v.loadApp(agent.AppID)
	//if err != nil {
	//	g.L.Warn("GetInstanceID:v.loadApp", zap.Error(err), zap.Any("agent", agent))
	//	return 0, err
	//}
	//
	////  agent.AgentUUID
	//UID := util.String2Uint32(agent.AgentUUID)
	//_, ok := app.Agents.Load(int32(UID))
	//if !ok {
	//	agent.InstanceID = int32(UID)
	//	app.Agents.Store(int32(UID), agent)
	//	v.storeAgent(agent)
	//}
	//
	//return int32(UID), nil

	return 0, nil
}

// storeAgent ...
func (v *Vgo) storeAgent(agent *util.AgentInfo) error {
	//query := fmt.Sprintf("insert into agent (instance_id, agent_uuid, app_id, app_name, os_name, ipv4s, register_time, process_id, host_name) values ('%d','%s','%d','%s','%s','%s','%d','%d','%s')",
	//	agent.InstanceID, agent.AgentUUID, agent.AppID, agent.AppName, agent.OsName, agent.Ipv4S, agent.RegisterTime, agent.ProcessID, agent.HostName)
	//_, err := g.DB.Exec(query)
	//if err != nil {
	//	g.L.Warn("loadInstanceID:g.DB.Exec", zap.String("query", query), zap.Error(err))
	//	return err
	//}
	return nil
}

//
//// dealSkywalking skywlking报文处理
//func (v *Vgo) dealSkywalking(conn net.Conn, packet *util.VgoPacket) error {
//	skypacker := &util.SkywalkingPacket{}
//	if err := msgpack.Unmarshal(packet.Payload, skypacker); err != nil {
//		g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//		return err
//	}
//	switch skypacker.Type {
//	// 应用注册
//	case util.TypeOfAppRegister:
//		appRegister := &util.KeyWithStringValue{}
//		if err := msgpack.Unmarshal(skypacker.Payload, appRegister); err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		id, err := v.GetAppID(appRegister.Value)
//		if err != nil {
//			g.L.Warn("dealSkywalking:v.apps.LoadAppCode", zap.String("name", appRegister.Value), zap.String("error", err.Error()))
//			return err
//		}
//
//		repPack := &util.KeyWithIntegerValue{
//			Key:   "id",
//			Value: id,
//		}
//
//		mbuf, err := msgpack.Marshal(repPack)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", appRegister.Value), zap.String("error", err.Error()))
//			return err
//		}
//		skypacker.Payload = mbuf
//
//		payload, err := msgpack.Marshal(skypacker)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", appRegister.Value), zap.String("error", err.Error()))
//			return err
//		}
//
//		packet.Payload = payload
//		if _, err := conn.Write(packet.Encode()); err != nil {
//			g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
//			return err
//		}
//		break
//		// 应用实例注册
//	case util.TypeOfAppRegisterInstance:
//		agentInfo := &util.AgentInfo{}
//		if err := msgpack.Unmarshal(skypacker.Payload, agentInfo); err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		id, err := v.GetInstanceID(agentInfo)
//		if err != nil {
//			g.L.Warn("dealSkywalking:v.apps.LoadAppCode", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
//			return err
//		}
//
//		appRegisterIns := &util.KeyWithIntegerValue{
//			Key:   "id",
//			Value: id,
//		}
//
//		mbuf, err := msgpack.Marshal(appRegisterIns)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
//			return err
//		}
//		skypacker.Payload = mbuf
//
//		payload, err := msgpack.Marshal(skypacker)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("name", agentInfo.AppName), zap.String("error", err.Error()))
//			return err
//		}
//
//		packet.Payload = payload
//
//		if _, err := conn.Write(packet.Encode()); err != nil {
//			g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
//			return err
//		}
//		break
//		// 应用服务名/注册
//	case util.TypeOfSerNameDiscoveryService:
//		repPacket := &util.SerNameDiscoveryServices{}
//		if err := msgpack.Unmarshal(skypacker.Payload, repPacket); err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		for index, serName := range repPacket.SerNames {
//			app, err := v.loadApp(serName.AppID)
//			if err != nil {
//				continue
//			}
//			id, err := v.loadAPI(serName, app)
//			if err != nil {
//				continue
//			}
//			repPacket.SerNames[index].SerID = id
//		}
//
//		mbuf, err := msgpack.Marshal(repPacket)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		skypacker.Payload = mbuf
//		payload, err := msgpack.Marshal(skypacker)
//		if err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Marshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		packet.Payload = payload
//
//		if _, err := conn.Write(packet.Encode()); err != nil {
//			g.L.Warn("dealSkywalking:conn.Write", zap.String("error", err.Error()))
//			return err
//		}
//		break
//		// 注册Addr
//	case util.TypeOfNewworkAddrRegister:
//		break
//		// jvm 数据
//	case util.TypeOfJVMMetrics:
//		repPacket := &util.JVMS{}
//		if err := msgpack.Unmarshal(skypacker.Payload, repPacket); err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		v.storage.jvmC <- repPacket
//		break
//
//		// trace 数据
//	case util.TypeOfTraceSegment:
//		var spans []*util.Span
//		if err := msgpack.Unmarshal(skypacker.Payload, &spans); err != nil {
//			g.L.Warn("dealSkywalking:msgpack.Unmarshal", zap.String("error", err.Error()))
//			return err
//		}
//
//		v.storage.spansC <- spans
//		break
//	}
//	return nil
//}
