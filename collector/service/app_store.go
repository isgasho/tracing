package service

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/imdevlab/tracing/collector/misc"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/pinpoint"
	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
	"github.com/imdevlab/tracing/pkg/sql"
)

// AppStore 所有app服务信息
type AppStore struct {
	sync.RWMutex
	apps map[string]*App // app集合
}

func (a *AppStore) start() error {
	if err := a.loadApps(); err != nil {
		logger.Warn("loadApps", zap.String("error", err.Error()))
		return err
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if err := a.loadApps(); err != nil {
				logger.Warn("loadApps", zap.String("error", err.Error()))
			}
		}
	}()
	return nil
}

func (a *AppStore) loadApps() error {
	cql := gCollector.storage.GetCql()
	if cql == nil {
		return fmt.Errorf("get cql failed")
	}

	appsIter := cql.Query(sql.LoadApps).Iter()
	defer func() {
		if err := appsIter.Close(); err != nil {
			logger.Warn("close apps iter error:", zap.Error(err))
		}
	}()

	// @TODO 这里未来要做hash过滤，不属于该collector节点App信息不需要保存，以节省资源
	var appName string
	for appsIter.Scan(&appName) {
		var appType int32
		var agentID, ip string
		var startTime int64

		agentsIter := cql.Query(sql.LoadAgents, appName).Iter()
		for agentsIter.Scan(&appType, &agentID, &startTime, &ip) {
			gCollector.apps.storeAgent(appName, agentID, startTime, appType)
			misc.AddrStore.Add(appName, ip)
		}
		if err := agentsIter.Close(); err != nil {
			logger.Warn("close apps iter error:", zap.Error(err))
		}
	}

	return nil
}

// isExist app是否存在
func (a *AppStore) isAppExist(name string) bool {
	a.RLock()
	_, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		return false
	}
	return true
}

func (a *AppStore) storeAgent(name string, id string, startTime int64, appType int32) {
	a.RLock()
	app, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		app = newApp(name, appType)
		a.Lock()
		a.apps[name] = app
		a.Unlock()
	}
	app.appType = appType
	app.storeAgent(id, startTime)
}

// isExist agent是否存在
func (a *AppStore) isAgentExist(name, agentid string) bool {
	a.RLock()
	app, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		return false
	}
	// app中是否存在
	return app.isExist(agentid)
}

func newAppStore() *AppStore {
	return &AppStore{
		apps: make(map[string]*App),
	}
}

func (a *AppStore) getApp(appName string) (*App, bool) {
	a.RLock()
	app, ok := a.apps[appName]
	a.RUnlock()
	return app, ok
}

// routerSapn 路由span
func (a *AppStore) routerStatBatch(appName, agentID string, stats *pinpoint.TAgentStatBatch) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, stats.GetStartTimestamp(), 0)

		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
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
func (a *AppStore) routerStat(appName, agentID string, stat *pinpoint.TAgentStat) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, stat.GetStartTimestamp(), 0)

		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
	}

	// 接收 stat
	if err := app.recvAgentStat(appName, agentID, stat); err != nil {
		logger.Warn("recv agent stat", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// routerSapn 路由span
func (a *AppStore) routerSapn(appName, agentID string, span *trace.TSpan) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, span.StartTime, int32(span.GetServiceType()))

		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
	}

	// 接收span
	if err := app.recvSpan(appName, agentID, span); err != nil {
		logger.Warn("recv span", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// routerSapnChunk 路由sapnChunk
func (a *AppStore) routersapnChunk(appName, agentID string, spanChunk *trace.TSpanChunk) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, spanChunk.AgentStartTime, int32(spanChunk.GetServiceType()))
		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
	}

	// 接收spanChunk
	if err := app.recvSpanChunk(appName, agentID, spanChunk); err != nil {
		logger.Warn("recv spanChunk", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}
