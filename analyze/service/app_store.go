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
	ticker := time.NewTicker(time.Duration(misc.Conf.Analyze.LoadAppInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := appStore.loadApp(); err != nil {
				g.L.Warn("loadApp", zap.String("error", err.Error()))
				break
			}

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

	query := fmt.Sprintf("SELECT app_name, last_count_time FROM apps; ")
	iterApp := appStore.db.Session.Query(query).Iter()
	defer iterApp.Close()

	var appName string
	var lastCountTime int64

	for iterApp.Scan(&appName, &lastCountTime) {
		key, err := gAnalyze.hash.Get(appName)
		if err != nil {
			continue
		}

		// 集群模式只做hash出来属于自己节点的APP
		if !strings.EqualFold(key, misc.Conf.Cluster.Name) {
			continue
		}

		app, isExist := appStore.getApp(appName)
		if !isExist {
			app = NewApp(appName)
			appStore.storeApp(app)
		}

		// 从 agent_stat 中取最早的启动时间记录
		if lastCountTime == 0 {
			queryStartTime := `SELECT start_time  FROM agent_stats  WHERE app_name=? LIMIT 1;`
			iterStartTime := appStore.db.Session.Query(queryStartTime, app.AppName).Iter()
			iterStartTime.Scan(&lastCountTime)
			iterStartTime.Close()
		}
		app.lastCountTime = lastCountTime
	}

	return nil
}
