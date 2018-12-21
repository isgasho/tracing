package service

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
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
	go appStore.LoadAppAndCounter()
	return nil
}

// Close ...
func (appStore *AppStore) Close() error {
	return nil
}

// LoadAppAndCounter ...
func (appStore *AppStore) LoadAppAndCounter() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := appStore.loadAppAndCounter(); err != nil {
				g.L.Warn("loadApp", zap.String("error", err.Error()))
			}
			break
		}
	}
}

// loadApp ...
func (appStore *AppStore) loadAppAndCounter() error {
	query := fmt.Sprintf("SELECT app_name, agent_id, start_time, is_live, last_point_time FROM agents ; ")
	iter := appStore.db.Session.Query(query).Iter()

	defer iter.Close()

	var wg sync.WaitGroup

	var appName, agentID string
	var startTime int64
	var isLive bool
	var lastPointTime int64
	for iter.Scan(&appName, &agentID, &startTime, &isLive, &lastPointTime) {
		key, err := gAnalyze.hash.Get(appName)
		if err != nil {
			continue
		}

		if !strings.EqualFold(key, misc.Conf.Cluster.Name) {
			continue
		}

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

	for _, app := range appStore.Apps {
		wg.Add(1)
		go gAnalyze.stats.counter(app, &wg)
	}

	wg.Wait()

	return nil
}
