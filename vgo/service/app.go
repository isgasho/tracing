package service

import (
	"fmt"
	"sync"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"go.uber.org/zap"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	Apps map[string]*util.App
}

// NewAppStore ...
func NewAppStore() *AppStore {
	return &AppStore{
		Apps: make(map[string]*util.App),
	}
}

// LoadApps 加载数据库中的所有app
func (as *AppStore) LoadApps() error {
	// 加载所有appCode
	apps := make([]*util.App, 0)
	if err := g.DB.Select(&apps, "select * from app"); err != nil {
		g.L.Fatal("LoadApps:g.DB.Select", zap.Error(err))
	}

	for _, app := range apps {
		as.Apps[app.Name] = app
	}

	g.L.Debug("LoadApps", zap.Any("apps", as.Apps))

	agents := make([]*util.AgentInfo, 0)
	if err := g.DB.Select(&agents, "select * from agent"); err != nil {
		g.L.Fatal("LoadApps:g.DB.Select", zap.Error(err))
	}

	for _, agent := range agents {
		app, ok := as.Apps[agent.AppName]
		if !ok {
			continue
		}
		if app.Agents == nil {
			app.Agents = make(map[int32]*util.AgentInfo)
		}
		app.Agents[agent.ID] = agent
	}

	return nil
}

// LoadAgentID 获取agentID
func (as *AppStore) LoadAgentID(agentInfo *util.AgentInfo) (int32, error) {
	as.RLock()
	app, ok := as.Apps[agentInfo.AppName]
	as.RUnlock()
	if !ok {
		return 0, fmt.Errorf("unfind app, app name is %s", agentInfo.AppName)
	}
	app.RLock()
	for _, agent := range app.Agents {
		if agent.AgentUUID == agentInfo.AgentUUID {
			app.RUnlock()
			return agent.ID, nil
		}
	}

	app.RUnlock()
	query := fmt.Sprintf("insert into agent (agent_uuid, app_code, app_name, os_name, ipv4s, register_time, process_id, host_name) values ('%s','%d','%s','%s','%s','%d','%d','%s')",
		agentInfo.AgentUUID, agentInfo.AppCode, agentInfo.AppName, agentInfo.OsName, agentInfo.Ipv4S, agentInfo.RegisterTime, agentInfo.ProcessID, agentInfo.HostName)
	// 如果不存在插入
	result, err := g.DB.Exec(query)
	if err != nil {
		g.L.Warn("LoadAgentID:g.DB.Exec", zap.String("query", query), zap.Error(err))
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		g.L.Warn("LoadAgentID:result.LastInsertId", zap.String("query", query), zap.Error(err))
		return 0, err
	}

	// 缓存
	agentInfo.ID = int32(id)
	app.Lock()
	if app.Agents == nil {
		app.Agents = make(map[int32]*util.AgentInfo)
	}
	app.Agents[int32(id)] = agentInfo
	app.Unlock()

	return int32(id), nil
}

// LoadAppCode  通过appname获取code
func (as *AppStore) LoadAppCode(name string) (int32, error) {
	as.RLock()
	app, ok := as.Apps[name]
	as.RUnlock()
	if ok {
		return app.Code, nil
	}

	// 如果不存在插入
	result, err := g.DB.Exec(fmt.Sprintf("insert into `app` (`name`) values ('%s')", name))
	if err != nil {
		g.L.Warn("LoadAppCode:g.DB.Exec", zap.Error(err))
		return 0, err
	}

	code, err := result.LastInsertId()
	if err != nil {
		g.L.Warn("LoadAppCode:result.LastInsertId", zap.Error(err))
		return 0, err
	}
	app = &util.App{
		Name: name,
		Code: int32(code),
	}

	// 缓存到内存中
	as.Lock()
	as.Apps[name] = app
	as.Unlock()

	return int32(code), nil
}
