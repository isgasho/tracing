package service

import (
	"strings"
	"sync"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/vgo/analyze/misc"
	"go.uber.org/zap"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	cql  *g.Cassandra
	Apps map[string]*App
}

func (appStore *AppStore) storeApp(app *App) {
	appStore.Lock()
	appStore.Apps[app.AppName] = app
	appStore.Unlock()
}

func (appStore *AppStore) getApp(appName string) (*App, bool) {
	appStore.RLock()
	app, isExist := appStore.Apps[appName]
	appStore.RUnlock()
	return app, isExist
}

// NewAppStore ...
func NewAppStore(cql *g.Cassandra) *AppStore {
	return &AppStore{
		cql:  cql,
		Apps: make(map[string]*App),
	}
}

// Start ...
func (appStore *AppStore) Start() error {
	go appStore.LoadAppAndCounter()
	return nil
}

// Close ...
func (appStore *AppStore) Close() error {
	return nil
}

// LoadAppAndCounter ...
func (appStore *AppStore) LoadAppAndCounter() {
	ticker := time.NewTicker(time.Duration(misc.Conf.Analyze.LoadAppInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			// 定时加载app，然后计算每一个APP的数据
			if err := appStore.loadApp(); err != nil {
				g.L.Warn("loadApp", zap.String("error", err.Error()))
				break
			}
			// 计算模块
			if err := gAnalyze.stats.Counter(); err != nil {
				g.L.Warn("Counter", zap.String("error", err.Error()))
				break
			}
			break
		}
	}
}

// loadApp ...
func (appStore *AppStore) loadApp() error {

	query := `SELECT app_name, last_count_time FROM apps; `
	iterApp := appStore.cql.Session.Query(query).Iter()

	defer func() {
		if err := iterApp.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}()

	var appName string
	var lastCountTime int64

	for iterApp.Scan(&appName, &lastCountTime) {
		key, err := gAnalyze.hash.Get(appName)
		if err != nil {
			continue
		}

		// 集群模式只做hash出来属于自己节点的APP
		if !strings.EqualFold(key, gAnalyze.clusterName) {
			continue
		}
		app, isExist := appStore.getApp(appName)
		if !isExist {
			app = NewApp(appName)
			appStore.storeApp(app)
		}

		// 从 agent_stat 中取最早的启动时间记录
		if lastCountTime == 0 {

			iterStartTime := appStore.cql.Session.Query(misc.QueryStartTime, app.AppName).Iter()
			iterStartTime.Scan(&lastCountTime)

			if err := iterStartTime.Close(); err != nil {
				g.L.Warn("close iter error:", zap.Error(err))
			}

			newMin, _ := ModMs2Min(lastCountTime)
			lastCountTime = newMin * 1000
		}
		app.lastCountTime = lastCountTime

		// load agents

		agentsIter := appStore.cql.Session.Query(misc.QueryAgents, appName).Iter()

		var agentID string
		var isLive bool
		for agentsIter.Scan(&agentID, &isLive) {
			if isLive {
				agent, isExist := app.getAgent(agentID)
				if !isExist {
					agent = NewAgent(agentID)
					app.storeAgent(agent)
				}
			} else {
				app.delAgent(agentID)
			}
		}
		if err := agentsIter.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}

	return nil
}
