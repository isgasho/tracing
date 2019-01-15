package service

import (
	"net"
	"sync"

	"github.com/mafanr/vgo/util"
	"github.com/mafanr/vgo/vgo/misc"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	Apps map[string]*App
}

// App app
type App struct {
	sync.RWMutex
	AppName string
	Agents  map[string]*Agent
}

// StoreAgent ...
func (appStore *AppStore) StoreAgent(agentInfo *util.AgentInfo, conn net.Conn) {
	appStore.Lock()
	app, ok := appStore.Apps[agentInfo.AppName]
	if !ok {
		app = NewApp()
		app.AppName = agentInfo.AppName
		appStore.Apps[agentInfo.AppName] = app
	}
	appStore.Unlock()

	app.Lock()
	agent, ok := app.Agents[agentInfo.AgentID]
	if !ok {
		agent = NewAgent()
		agent.AgentID = agentInfo.AgentID
		app.Agents[agentInfo.AgentID] = agent
	}
	agent.conn = conn
	app.Unlock()

	return
}

// RemoveAgent ...
func (appStore *AppStore) RemoveAgent(agentInfo *util.AgentInfo) {
	appStore.RLock()
	app, ok := appStore.Apps[agentInfo.AppName]
	appStore.RUnlock()
	if !ok {
		return
	}

	app.Lock()
	delete(app.Agents, agentInfo.AgentID)
	app.Unlock()
}

// NewApp ...
func NewApp() *App {
	return &App{
		Agents: make(map[string]*Agent),
	}
}

// Agent ....
type Agent struct {
	sync.RWMutex
	AgentID string
	conn    net.Conn
	Apis    map[int32]*API
}

// NewAgent ...
func NewAgent() *Agent {
	return &Agent{
		Apis: make(map[int32]*API),
	}
}

// API ...
type API struct {
	ID     int32
	APIStr string
}

// NewAPI ...
func NewAPI() *API {
	return &API{}
}

// NewAppStore ...
func NewAppStore() *AppStore {
	return &AppStore{
		Apps: make(map[string]*App),
	}
}

func (appStore *AppStore) checkApp(appName string) bool {
	appStore.RLock()
	_, ok := appStore.Apps[appName]
	appStore.RUnlock()
	if !ok {
		iter := gVgo.storage.cql.Query(misc.CheckApp, appName).Iter()
		var count int
		iter.Scan(&count)
		iter.Close()
		if count == 0 {
			return false
		}

		appStore.Lock()
		app := NewApp()
		app.AppName = appName
		appStore.Apps[appName] = app
		appStore.Unlock()
		return true
	}
	return ok
}

func (appStore *AppStore) checkAndSaveAgent(appName, agentID string) bool {
	appStore.RLock()
	app, ok := appStore.Apps[appName]
	appStore.RUnlock()
	if !ok {
		appStore.Lock()
		app = NewApp()
		app.AppName = appName

		// 缓存Agent
		agent := NewAgent()
		app.Agents[agentID] = agent

		// 缓存App
		appStore.Apps[appName] = app
		appStore.Unlock()

		return false
	}

	// 缓存Agent
	app.RLock()
	agent, ok := app.Agents[agentID]
	app.RUnlock()
	if !ok {
		agent = NewAgent()
		app.Lock()
		app.Agents[agentID] = agent
		app.Unlock()

		return false
	}
	return true
}

func (appStore *AppStore) checkAndSaveAPIID(appName, agentID string, apiID int32) bool {

	appStore.RLock()
	app, ok := appStore.Apps[appName]
	appStore.RUnlock()
	if !ok {
		return false
	}

	app.RLock()
	agent, ok := app.Agents[agentID]
	app.RUnlock()
	if !ok {
		return false
	}

	agent.RLock()
	api, ok := agent.Apis[apiID]
	agent.RUnlock()
	if !ok {
		api = NewAPI()
		api.ID = apiID
		agent.Lock()
		agent.Apis[apiID] = api
		agent.Unlock()
		return false
	}
	return true
}
