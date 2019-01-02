package service

import "sync"

// App ...
type App struct {
	sync.RWMutex
	AppName       string
	lastCountTime int64
	Agents        map[string]*Agent
}

// NewApp ...
func NewApp(appName string) *App {
	return &App{
		AppName: appName,
		Agents:  make(map[string]*Agent),
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

func (app *App) delAgent(agentID string) {
	app.Lock()
	delete(app.Agents, agentID)
	app.Unlock()
}
