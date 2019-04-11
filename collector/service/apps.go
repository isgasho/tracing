package service

import (
	"fmt"
	"sync"

	"github.com/imdevlab/g"
	"go.uber.org/zap"

	"github.com/imdevlab/tracing/pkg/pinpoint/thrift/trace"
)

// Apps 所有app服务信息
type Apps struct {
	sync.RWMutex
	apps map[string]*App // app集合
}

// isExist app是否存在
func (a *Apps) isAppExist(name string) bool {
	a.RLock()
	_, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		return false
	}
	return true
}

func (a *Apps) storeAgent(name string, id string, startTime int64) {
	a.RLock()
	app, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		app = newApp(name)
		a.Lock()
		a.apps[name] = app
		a.Unlock()
	}
	app.storeAgent(id, startTime)
}

// isExist agent是否存在
func (a *Apps) isAgentExist(name, agentid string) bool {

	a.RLock()
	app, ok := a.apps[name]
	a.RUnlock()
	if !ok {
		return false
	}

	// app中是否存在
	return app.isExist(agentid)
}

func newApps() *Apps {
	return &Apps{
		apps: make(map[string]*App),
	}
}

func (a *Apps) getApp(appName string) (*App, bool) {
	a.RLock()
	app, ok := a.apps[appName]
	a.RUnlock()
	return app, ok
}

// routerSapn 路由span
func (a *Apps) routerSapn(appName, agentID string, span *trace.TSpan) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, span.StartTime)

		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
	}

	// 接收span
	if err := app.recvSpan(appName, agentID, span); err != nil {
		g.L.Warn("recv span", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}

// routerSapnChunk 路由sapnChunk
func (a *Apps) routersapnChunk(appName, agentID string, spanChunk *trace.TSpanChunk) error {
	app, ok := a.getApp(appName)
	if !ok {
		// 缓存App
		a.storeAgent(appName, agentID, spanChunk.AgentStartTime)

		// 新App在重新找一次
		app, ok = a.getApp(appName)
		if !ok {
			return fmt.Errorf("unfind app, app name is %s", appName)
		}
	}

	// 接收spanChunk
	if err := app.recvSpanChunk(appName, agentID, spanChunk); err != nil {
		g.L.Warn("recv spanChunk", zap.String("appName", appName), zap.String("agentID", agentID), zap.String("error", err.Error()))
		return err
	}

	return nil
}
