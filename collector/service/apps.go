package service

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gocql/gocql"
	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/sql"
)

// Apps 所有app服务信息
type Apps struct {
	sync.RWMutex
	apps  map[string]*App // app集合
	ips   map[string]string
	hosts map[string]string
}

func (a *Apps) start() error {
	cql := gCollector.storage.GetCql()
	if err := a.loadApps(cql); err != nil {
		logger.Warn("loadApps", zap.String("error", err.Error()))
		return err
	}

	go func() {
		for {
			time.Sleep(time.Duration(misc.Conf.Apps.LoadInterval) * time.Second)
			cql := gCollector.storage.GetCql()
			if err := a.loadApps(cql); err != nil {
				logger.Warn("loadApps", zap.String("error", err.Error()))
			}
		}
	}()
	return nil
}

func (a *Apps) loadApps(cql *gocql.Session) error {
	if cql == nil {
		return fmt.Errorf("get cql failed")
	}
	appsIter := cql.Query(sql.LoadApps).Iter()
	defer func() {
		if err := appsIter.Close(); err != nil {
			logger.Warn("close apps iter error:", zap.Error(err))
		}
	}()

	var appName string
	for appsIter.Scan(&appName) {
		var appType int32
		var agentID, ip string
		var startTime int64
		var isLive bool
		var hostName string

		agentsIter := cql.Query(sql.LoadAgents, appName).Iter()
		for agentsIter.Scan(&appType, &agentID, &startTime, &ip, &isLive, &hostName) {
			a.storeAgent(appName, agentID, appType, startTime, isLive, hostName, ip)
			a.storeIPandHost(appName, ip, hostName)
		}

		if err := agentsIter.Close(); err != nil {
			logger.Warn("close apps iter error:", zap.Error(err))
		}
	}

	return nil
}

// isExist app是否存在
func (a *Apps) storeIPandHost(appName, ip, host string) bool {
	// a.Lock()
	a.ips[ip] = appName
	a.hosts[ip] = appName
	// a.Unlock()
	return true
}

func (a *Apps) getNameByIP(ip string) (string, bool) {
	name, ok := a.ips[ip]
	return name, ok
}

func (a *Apps) getNameByHost(host string) (string, bool) {
	name, ok := a.hosts[host]
	return name, ok
}

func (a *Apps) storeAgent(appName, agentID string, appType int32, startTime int64, isLive bool, hostName, ip string) {
	a.RLock()
	app, ok := a.apps[appName]
	a.RUnlock()
	if !ok {
		app = newApp(appName)
		a.Lock()
		a.apps[appName] = app
		a.Unlock()

		app.start()
	}
	app.appType = appType
	app.storeAgent(agentID, isLive)
	// app.storeAgent(id, startTime)
}

func newApps() *Apps {
	return &Apps{
		apps:  make(map[string]*App),
		ips:   make(map[string]string),
		hosts: make(map[string]string),
	}
}

func (a *Apps) getApp(appName string) (*App, bool) {
	a.RLock()
	app, ok := a.apps[appName]
	a.RUnlock()
	return app, ok
}

// routerSapn 路由span
func (a *Apps) routerStatBatch(appName, agentID string, stats *pinpoint.TAgentStatBatch) error {
	app, ok := a.getApp(appName)
	if !ok {
		return fmt.Errorf("unfind app, app name is %s", appName)
	}

	// 接收 stat
	for _, stat := range stats.AgentStats {
		if err := app.recvAgentStat(appName, agentID, stat); err != nil {
			logger.Warn("recv agent stat", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
			return err
		}
	}

	return nil
}

// routerSapn 路由span
func (a *Apps) routerStat(appName, agentID string, stat *pinpoint.TAgentStat) error {
	app, ok := a.getApp(appName)
	if !ok {
		return fmt.Errorf("unfind app, app name is %s", appName)
	}

	// 接收 stat
	if err := app.recvAgentStat(appName, agentID, stat); err != nil {
		logger.Warn("recv agent stat", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// routerSapn 路由span
func (a *Apps) routerSapn(appName, agentID string, span *trace.TSpan) error {
	app, ok := a.getApp(appName)
	if !ok {
		return fmt.Errorf("unfind app, app name is %s", appName)
	}

	// 接收span
	if err := app.recvSpan(appName, agentID, span); err != nil {
		logger.Warn("recv span", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// routerSapnChunk 路由sapnChunk
func (a *Apps) routersapnChunk(appName, agentID string, spanChunk *trace.TSpanChunk) error {
	app, ok := a.getApp(appName)
	if !ok {
		return fmt.Errorf("unfind app, app name is %s", appName)
	}

	// 接收spanChunk
	if err := app.recvSpanChunk(appName, agentID, spanChunk); err != nil {
		logger.Warn("recv spanChunk", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}
