package service

import (
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/network"
	"github.com/vmihailenco/msgpack"

	"github.com/imdevlab/g"
	"go.uber.org/zap"
)

// Agent ...
type Agent struct {
	appName      string             // 服务名
	agentID      string             // 服务agent ID
	etcd         *Etcd              // 服务发现
	collector    *Collector         // 监控指标上报
	pinpoint     *Pinpoint          // pinpoint采集服务
	isLive       bool               // app是否存活
	isReportInfo bool               // 是否已经上报agentInfo
	syncID       uint32             // 同步请求ID
	syncCall     *SyncCall          // 同步请求
	agentInfo    *network.AgentInfo // 监控上报的agent info原信息
}

var gAgent *Agent

// New new agent
func New() *Agent {
	gAgent = &Agent{
		etcd:      newEtcd(),
		collector: newCollector(),
		pinpoint:  newPinpoint(),
		syncCall:  NewSyncCall(),
		agentInfo: network.NewAgentInfo(),
	}
	return gAgent
}

// Start 启动
func (a *Agent) Start() error {

	// etcd 初始化
	if err := a.etcd.Init(); err != nil {
		g.L.Warn("etcd init", zap.String("error", err.Error()))
		return err
	}

	// 启动服务发现
	if err := a.etcd.Start(); err != nil {
		g.L.Warn("etcd start", zap.String("error", err.Error()))
		return err
	}

	// 监控采集服务启动
	if err := a.pinpoint.Start(); err != nil {
		g.L.Warn("pinpoint start", zap.String("error", err.Error()))
		return err
	}

	// agent 信息上报服务
	go reportAgentInfo()

	g.L.Info("Agent start ok")

	return nil
}

// Close 关闭
func (a *Agent) Close() error {
	return nil
}

func getApplicationName() error {
	hostname, err := os.Hostname()
	if err != nil {
		g.L.Warn("get host name error", zap.Error(err))
		return err
	}

	names := strings.Split(hostname, "-")
	if len(names) == 1 {
		gAgent.appName = hostname
	} else if len(names) == 3 {
		gAgent.appName = names[1]
	} else if len(names) == 4 {
		gAgent.appName = names[1] + names[2]
	}
	return nil
}

func getAgentID() error {
	hostname, err := os.Hostname()
	if err != nil {
		g.L.Warn("get host name error", zap.Error(err))
		return err
	}
	names := strings.Split(hostname, "-")
	if len(names) == 1 {
		gAgent.agentID = hostname
	} else if len(names) == 3 {
		var id string
		if strings.EqualFold(names[2], "vip") {
			id = "v"
		} else if strings.EqualFold(names[2], "yf") {
			id = "y"
		} else {
			id = names[2]
		}
		gAgent.agentID = names[1] + id
	} else if len(names) == 4 {
		var id string
		if strings.EqualFold(names[3], "vip") {
			id = "v"
		} else if strings.EqualFold(names[3], "yf") {
			id = "y"
		} else {
			id = names[3]
		}
		gAgent.agentID = names[1] + names[2] + id
	}
	return nil
}

// getAppname 获取本机App名
func (a *Agent) getAppname() error {

	getApplicationName()
	getAgentID()

	g.L.Info("init", zap.String("appName", a.appName))
	g.L.Info("init", zap.String("agentID", a.agentID))

	return nil
}

// reportAgentInfo 上报agent 信息
func reportAgentInfo() {
	for {
		time.Sleep(1 * time.Second)
		if !gAgent.isLive {
			continue
		}
		break
	}
	for {
		time.Sleep(3 * time.Second)
		if !gAgent.isReportInfo {
			spanPackets := network.NewSpansPacket()
			spanPackets.Type = constant.TypeOfTCPData
			spanPackets.AppName = gAgent.appName
			spanPackets.AgentID = gAgent.agentID

			agentInfo, err := msgpack.Marshal(gAgent.agentInfo)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				continue
			}
			spans := &network.Spans{
				Spans: agentInfo,
			}
			if gAgent.isLive == false {
				spans.Type = constant.TypeOfAgentOffline
			} else {
				spans.Type = constant.TypeOfRegister
			}

			spanPackets.Payload = append(spanPackets.Payload, spans)
			payload, err := msgpack.Marshal(spanPackets)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				continue
			}

			id := gAgent.getSyncID()
			tracePacket := &network.TracePack{
				Type:       constant.TypeOfPinpoint,
				IsSync:     constant.TypeOfSyncYes,
				IsCompress: constant.TypeOfCompressNo,
				ID:         id,
				Payload:    payload,
			}

			if err := gAgent.collector.write(tracePacket); err != nil {
				g.L.Warn("write info", zap.String("error", err.Error()))
				continue
			}

			// 创建chan
			if _, ok := gAgent.syncCall.newChan(id, 10); !ok {
				g.L.Warn("syncCall newChan", zap.String("error", "创建sync chan失败"))
				continue
			}

			// 阻塞同步等待，并关闭chan
			if _, err := gAgent.syncCall.syncRead(id, 10, true); err != nil {
				g.L.Warn("syncRead", zap.String("error", err.Error()))
				continue
			}

			gAgent.isReportInfo = true
		}
	}
}

// getSyncID ...
func (a *Agent) getSyncID() uint32 {
	return atomic.AddUint32(&a.syncID, 1)
}
