package service

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/mafanr/vgo/agent/misc"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"

	"go.uber.org/zap"
)

// Agent ...
type Agent struct {
	isGetID       bool
	appName       string // 应用名
	appID         int32  //	应用ID
	appInstanceID int32  // 应用实例ID
	agentUUID     string
	syncCall      *SyncCall
	client        *TCPClient
	skyWalk       *SkyWalking
	syncID        uint32
	agentInfo     *util.AgentInfo
	quitC         chan bool
	uploadC       chan *util.VgoPacket
	downloadC     chan *util.VgoPacket
}

var gAgent *Agent

// New ...
func New() *Agent {
	gAgent = &Agent{
		isGetID:   false,
		syncCall:  NewSyncCall(),
		client:    NewTCPClient(),
		skyWalk:   NewSkyWalking(),
		agentInfo: util.NewAgentInfo(),
		quitC:     make(chan bool, 1),
		uploadC:   make(chan *util.VgoPacket, 100),
		downloadC: make(chan *util.VgoPacket, 100),
	}
	return gAgent
}

// getSyncID ...
func (a *Agent) getSyncID() uint32 {
	return atomic.AddUint32(&a.syncID, 1)
}

// initAppName ...
func (a *Agent) initAppName() error {
	// 使用环境变量
	if misc.Conf.Agent.UseEnv {
		if misc.Conf.Agent.ENV == "" {
			g.L.Fatal("initAppName:.", zap.Error(fmt.Errorf("env is nil")))
		}
		name := os.Getenv(misc.Conf.Agent.ENV)
		if name == "" {
			g.L.Fatal("initAppName:os.Getenv", zap.Error(fmt.Errorf("get env is nil")), zap.String("env", misc.Conf.Agent.ENV))
		}
		a.appName = name
	} else {
		// 从配置文件获取
		if len(misc.Conf.Agent.AppName) > 0 {
			a.appName = misc.Conf.Agent.AppName
		} else {
			// 从主机名获取
			_, agentName := getAgentIdAndName()
			a.appName = agentName
		}
	}

	g.L.Info("initAppName", zap.String("AppName", a.appName))

	return nil
}

// reloadID ...
func (a *Agent) reloadID(appID, instanceID int32) {
	if !a.isGetID {
		a.appID = appID
		a.appInstanceID = instanceID
	}
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

	// 启动本地接收采集信息端口
	if err := a.skyWalk.Start(); err != nil {
		g.L.Fatal("Start:a.skyWalk.Start", zap.Error(err))
	}

	return nil
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

func getAgentIdAndName() (agentId, agentName string) {
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
