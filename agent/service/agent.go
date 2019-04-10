package service

import (
	"sync/atomic"
	"time"

	"github.com/shaocongcong/tracing/pkg/proto/network"
	"github.com/shaocongcong/tracing/pkg/proto/ttype"
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
	}
	return gAgent
}

// Start 启动
func (a *Agent) Start() error {

	// 获取Appname
	if err := a.getAppname(); err != nil {
		g.L.Warn("get app name", zap.String("error", err.Error()))
		return err
	}

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

// getAppname 获取本机App名
func (a *Agent) getAppname() error {
	a.appName = "test"
	a.agentID = "test-1"
	return nil
}

// @TODO 上报agent info机制需要优化，在各种断线&collector服务切换的情况都需要考虑
func reportAgentInfo() {
	for {
		time.Sleep(3 * time.Second)
		if !gAgent.isReportInfo {
			infoPack := network.NewSpansPacket()
			infoPack.Type = ttype.TypeOfTCPData
			infoPack.AppName = gAgent.appName
			infoPack.AgentID = gAgent.agentID
			infoBuf, err := msgpack.Marshal(infoPack)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				continue
			}
			spanData := &network.Spans{
				Spans: infoBuf,
			}

			if gAgent.isLive == false {
				spanData.Type = ttype.TypeOfAgentOffline
			} else {
				spanData.Type = ttype.TypeOfRegister
			}

			infoPack.Payload = append(infoPack.Payload, spanData)
			payload, err := msgpack.Marshal(infoPack)
			if err != nil {
				g.L.Warn("msgpack Marshal", zap.String("error", err.Error()))
				continue
			}

			id := gAgent.getSyncID()
			packet := &network.TracePack{
				Type:       ttype.TypeOfPinpoint,
				IsSync:     ttype.TypeOfSyncYes,
				IsCompress: ttype.TypeOfCompressNo,
				ID:         id,
				Payload:    payload,
			}

			if err := gAgent.collector.write(packet); err != nil {
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
