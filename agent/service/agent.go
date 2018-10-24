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
	quitC     chan bool
	uploadC   chan *util.VgoPacket
	downloadC chan *util.VgoPacket
	syncCall  *SyncCall
	client    *TCPClient
	skyWalk   *SkyWalking
	id        uint32
	appInfo   *AppInfo
}

var gAgent *Agent

// New ...
func New() *Agent {
	gAgent = &Agent{
		quitC:     make(chan bool, 1),
		uploadC:   make(chan *util.VgoPacket, 100),
		downloadC: make(chan *util.VgoPacket, 100),
		syncCall:  NewSyncCall(),
		client:    NewTCPClient(),
		skyWalk:   NewSkyWalking(),
		appInfo:   NewAppInfo(),
	}
	return gAgent
}

// ID ...
func (a *Agent) ID() uint32 {
	return atomic.AddUint32(&a.id, 1)
}

// initAppInfo ...
func (a *Agent) initAppInfo() error {
	// 使用环境变量
	if misc.Conf.Agent.UseEnv {
		if misc.Conf.Agent.ENV == "" {
			g.L.Fatal("initAppInfo:.", zap.Error(fmt.Errorf("env is nil")))
		}
		name := os.Getenv(misc.Conf.Agent.ENV)
		if name == "" {
			g.L.Fatal("initAppInfo:os.Getenv", zap.Error(fmt.Errorf("get env is nil")), zap.String("env", misc.Conf.Agent.ENV))
		}
		a.appInfo.Name = name
	} else {
		// 从配置文件获取
		if len(misc.Conf.Agent.AppName) > 0 {
			a.appInfo.Name = misc.Conf.Agent.AppName
		} else {
			// 从主机名获取
			_, agentName := getAgentIdAndName()
			a.appInfo.Name = agentName
		}
	}

	g.L.Info("initAppInfo", zap.String("AppName", a.appInfo.Name))

	return nil
}

// Start ...
func (a *Agent) Start() error {

	//	获取App信息
	if err := a.initAppInfo(); err != nil {
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
