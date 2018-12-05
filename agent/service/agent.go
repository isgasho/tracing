package service

import (
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/vmihailenco/msgpack"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"

	"go.uber.org/zap"
)

// Agent ...
type Agent struct {
	// appName           string // 应用名
	// agentID           string //	应用ID
	syncCall          *SyncCall
	client            *TCPClient
	syncID            uint32
	agentInfo         *util.AgentInfo
	quitC             chan bool
	uploadC           chan *util.VgoPacket
	downloadC         chan *util.VgoPacket
	pinpoint          *Pinpoint
	isReportAgentInfo bool
}

var gAgent *Agent

// New ...
func New() *Agent {
	gAgent = &Agent{
		syncCall:  NewSyncCall(),
		client:    NewTCPClient(),
		agentInfo: util.NewAgentInfo(),
		quitC:     make(chan bool, 1),
		uploadC:   make(chan *util.VgoPacket, 100),
		downloadC: make(chan *util.VgoPacket, 100),
		pinpoint:  NewPinpoint(),
	}
	return gAgent
}

// getSyncID ...
func (a *Agent) getSyncID() uint32 {
	return atomic.AddUint32(&a.syncID, 1)
}

// initAppName ...
func (a *Agent) initAppName() error {

	return nil
	// 使用环境变量
	// if misc.Conf.Agent.UseEnv {
	// 	if misc.Conf.Agent.ENV == "" {
	// 		g.L.Fatal("initAppName:.", zap.Error(fmt.Errorf("env is nil")))
	// 	}
	// 	name := os.Getenv(misc.Conf.Agent.ENV)
	// 	if name == "" {
	// 		g.L.Fatal("initAppName:os.Getenv", zap.Error(fmt.Errorf("get env is nil")), zap.String("env", misc.Conf.Agent.ENV))
	// 	}
	// 	a.appName = name
	// } else {
	// 	// 从配置文件获取
	// 	if len(misc.Conf.Agent.AppName) > 0 {
	// 		a.appName = misc.Conf.Agent.AppName
	// 	} else {
	// 		// 从主机名获取
	// 		_, agentName := getAgentIDAndName()
	// 		a.appName = agentName
	// 	}
	// }

	// g.L.Info("initAppName", zap.String("AppName", a.appName))

	return nil
}

// setAgentInfo ...
func (a *Agent) setAgentInfo(agentID string) {
	// a.agentID = agentID
}

// Start ...
func (a *Agent) Start() error {
	// 获取App name
	if err := a.initAppName(); err != nil {
		g.L.Fatal("Start:a.initAppInfo", zap.Error(err))
	}

	// 启动upload
	go a.upload()

	// 初始化处理下行命令等
	go a.download()

	// 初始化tcp client
	go a.client.Init()

	// 上报agent信息
	go a.reportAgentInfo()

	// start pinpoint
	if err := a.pinpoint.Start(); err != nil {
		g.L.Fatal("Start:a.pinpoint.Start", zap.Error(err))
	}

	return nil
}

func (a *Agent) reportAgentInfo() {
	for {
		time.Sleep(3 * time.Second)
		if a.isReportAgentInfo {
			pinpointData := util.NewPinpointData()
			pinpointData.Type = util.TypeOfTCPData
			pinpointData.AgentName = a.agentInfo.AppName
			pinpointData.AgentID = a.agentInfo.AgentID
			agentInfoBuf, err := msgpack.Marshal(a.agentInfo)
			if err != nil {
				g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
				continue
			}

			spanData := &util.SpanDataModel{
				Type:  util.TypeOfRegister,
				Spans: agentInfoBuf,
			}
			pinpointData.Payload = append(pinpointData.Payload, spanData)
			payload, err := msgpack.Marshal(pinpointData)
			if err != nil {
				g.L.Warn("agentInfo:msgpack.Marshal", zap.String("error", err.Error()))
				continue
			}

			// 获取ID
			id := gAgent.getSyncID()
			packet := &util.VgoPacket{
				Type:       util.TypeOfPinpoint,
				Version:    util.VersionOf01,
				IsSync:     util.TypeOfSyncYes,
				IsCompress: util.TypeOfCompressNo,
				ID:         id,
				Payload:    payload,
			}

			if err := gAgent.client.WritePacket(packet); err != nil {
				g.L.Warn("ApplicationCodeRegister:gAgent.client.WritePacket", zap.String("error", err.Error()))
				continue
			}

			// 创建chan
			if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
				g.L.Warn("ApplicationCodeRegister:gAgent.syncCall.newChan", zap.String("error", "创建sync chan失败"))
				continue
			}

			// 阻塞同步等待，并关闭chan
			if _, err := gAgent.syncCall.syncRead(id, 10, true); err != nil {
				g.L.Warn("ApplicationCodeRegister:gAgent.syncCall.syncRead", zap.String("error", err.Error()))
				continue
			}
			// 上报成功无须上报
			a.isReportAgentInfo = false
		}
	}
}

func (a *Agent) write(data *util.VgoPacket) {

}

// Close ...
func (a *Agent) Close() error {
	return nil
}

func (a *Agent) upload() {
	defer func() {
		if err := recover(); err != nil {
			g.L.Warn("report:.", zap.Stack("server"), zap.Any("err", err))
		}
	}()

	for {
		select {
		case p, ok := <-a.uploadC:
			if ok {
				if err := a.client.WritePacket(p); err != nil {
					g.L.Warn("report:client.WritePacket", zap.String("error", err.Error()))
				}
			}
			break
		}
	}
}

func (a *Agent) download() {
	for {
		select {
		case p, ok := <-a.downloadC:
			if ok {
				g.L.Info("cmd", zap.Any("msg", p))
			}
		case <-a.quitC:
			return
		}
	}
}

// getAgentIDAndName ...
func getAgentIDAndName() (agentId, agentName string) {
	host, err := os.Hostname()
	if err != nil {
		log.Fatalln("[FATAL] get hostname error: ", err)
	}
	hostS := strings.Split(host, "-")
	if len(hostS) == 1 {
		return host, host
	} else if len(hostS) == 3 {
		var id string
		if strings.ToLower(hostS[2]) == "vip" {
			id = "v"
		} else if strings.ToLower(hostS[2]) == "yf" {
			id = "y"
		} else {
			id = hostS[2]
		}
		return hostS[1] + id, hostS[1]
	} else if len(hostS) == 4 {
		var id string
		if strings.ToLower(hostS[3]) == "vip" {
			id = "v"
		} else if strings.ToLower(hostS[3]) == "yf" {
			id = "y"
		} else {
			id = hostS[3]
		}
		return hostS[1] + hostS[2] + id, hostS[1] + hostS[2]
	}
	return "", ""
}

// GetHostName get host name
func GetHostName() (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return host, err
}
