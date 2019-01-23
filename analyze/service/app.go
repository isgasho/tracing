package service

import "sync"

// App ...
type App struct {
	sync.RWMutex
	AppName       string
	lastCountTime int64
	Agents        map[string]*Agent
	URLs          map[string]struct{}
}

// NewApp ...
func NewApp(appName string) *App {
	return &App{
		AppName: appName,
		Agents:  make(map[string]*Agent),
		URLs:    make(map[string]struct{}),
	}
}

func (app *App) getAgent(agentID string) (*Agent, bool) {
	app.RLock()
	agent, isExist := app.Agents[agentID]
	app.RUnlock()
	return agent, isExist
}

func (app *App) storeAgent(agent *Agent) {
	app.Lock()
	app.Agents[agent.AgentID] = agent
	app.Unlock()
}

func (app *App) getURL(url string) (struct{}, bool) {
	app.RLock()
	v, ok := app.URLs[url]
	app.RUnlock()
	return v, ok
}

func (app *App) storeURL(url string) {
	app.Lock()
	app.URLs[url] = struct{}{}
	app.Unlock()
}
func (app *App) delAgent(agentID string) {
	app.Lock()
	delete(app.Agents, agentID)
	app.Unlock()
}
