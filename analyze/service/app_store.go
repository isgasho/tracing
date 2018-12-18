package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/mafanr/g"
	"go.uber.org/zap"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	db   *g.Cassandra
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
func NewAppStore(db *g.Cassandra) *AppStore {
	return &AppStore{
		db:   db,
		Apps: make(map[string]*App),
	}
}

// Start ...
func (appStore *AppStore) Start() error {
	go appStore.LoadApp()
	return nil
}

// Close ...
func (appStore *AppStore) Close() error {
	return nil
}

// LoadApp ...
func (appStore *AppStore) LoadApp() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := appStore.loadApp(); err != nil {
				g.L.Warn("loadApp", zap.String("error", err.Error()))
			}
			break
		}
	}
}

// loadApp ...
func (appStore *AppStore) loadApp() error {
	query := fmt.Sprintf("SELECT app_name, agent_id, start_time, is_live, last_point_time FROM agents ; ")
	iter := appStore.db.Session.Query(query).Iter()

	defer iter.Close()

	var appName, agentID string
	var startTime int64
	var isLive bool
	var lastPointTime int64
	for iter.Scan(&appName, &agentID, &startTime, &isLive, &lastPointTime) {
		app, isExist := appStore.getApp(appName)
		if !isExist {
			app = NewApp(appName)
			appStore.storeApp(app)
		}
		agent, isExist := app.getAgent(agentID)
		if !isExist {
			agent = NewAgent(agentID)
			app.storeAgent(agent)
		}

		agent.startTime = startTime
		agent.isLive = isLive
		agent.lastPointTime = lastPointTime
	}

	return nil
}
