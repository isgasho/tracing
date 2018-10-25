package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mafanr/g"
	"github.com/mafanr/vgo/util"
	"go.uber.org/zap"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	Apps     map[string]*util.App
	AppCodes map[int32]string
	slock    sync.RWMutex
	SerNames map[int32]*Apis
}

// Apis ...
type Apis struct {
	Apis map[int32]*util.SerNameInfo
}

// NewApis ...
func NewApis() *Apis {
	return &Apis{
		Apis: make(map[int32]*util.SerNameInfo),
	}
}

// NewAppStore ...
func NewAppStore() *AppStore {
	return &AppStore{
		Apps:     make(map[string]*util.App),
		AppCodes: make(map[int32]string),
		SerNames: make(map[int32]*Apis),
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
		as.AppCodes[app.Code] = app.Name
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
	as.AppCodes[app.Code] = app.Name
	as.Unlock()

	return int32(code), nil
}

// LoadAppName  通过Code获取Name
func (as *AppStore) LoadAppName(code int32) (string, bool) {
	as.RLock()
	defer as.RUnlock()
	name, ok := as.AppCodes[code]
	return name, ok
}

// // LoadSerCode ...
// func (as *AppStore) LoadSerCode() error {
// 	// 如果不存在插入
// 	log.Println("测试 测试")
// 	// insert into `vgo`.`server_name` ( `id`, ) values ( '3', '3', '3', '3')
// 	result, err := g.DB.Exec(fmt.Sprintf("insert into `server_name` (`server_name`, `span_type`, `app_code`) values ('2' , '2' '2')"))

// 	log.Println("测试 测试", err)
// 	if err != nil {
// 		g.L.Warn("LoadAppCode:g.DB.Exec", zap.Error(err), zap.Any("result", result))
// 		return err
// 	}

// 	code, err := result.LastInsertId()
// 	if err != nil {
// 		g.L.Warn("LoadAppCode:result.LastInsertId", zap.Error(err))
// 		return err
// 	}
// 	log.Println("测试 测试 xxx")
// 	log.Println("xxxxxxxxxxxxxxxxxxx", code)
// 	return nil
// }

// LoadSerCode ...
func (as *AppStore) LoadSerCode() error {
	// 加载所有server name code
	serNames := make([]*util.SerNameInfo, 0)
	if err := g.DB.Select(&serNames, "select * from server_name"); err != nil {
		g.L.Fatal("LoadSerCode:g.DB.Select", zap.Error(err))
		return err
	}

	for _, sInfo := range serNames {
		api, ok := as.SerNames[sInfo.AppCode]
		if !ok {
			api = NewApis()
			as.SerNames[sInfo.AppCode] = api
		}
		api.Apis[sInfo.SerID] = sInfo
	}

	g.L.Debug("LoadSerCode", zap.Any("apps", as.SerNames))
	return nil
}

// GetSerCode  通过Code&name获取code
func (as *AppStore) GetSerCode(serInfo *util.SerNameInfo) (int32, error) {
	var id int64
	as.slock.RLock()
	apis, ok := as.SerNames[serInfo.AppCode]
	as.slock.RUnlock()
	if ok {
		for _, api := range apis.Apis {
			if strings.EqualFold(api.SerName, serInfo.SerName) {
				return api.SerID, nil
			}
		}
	}

	result, err := g.DB.Exec(fmt.Sprintf("insert into `server_name` (`server_name`, `span_type`, `app_code`) values ('%s' , '%d' '%d')", serInfo.SerName,
		serInfo.SpanType, serInfo.AppCode))
	if err != nil {
		query := fmt.Sprintf("select id from server_name where app_code='%d' and server_name='%s'", serInfo.AppCode, serInfo.SerName)
		rows, err := g.DB.Query(query)
		if err != nil {
			g.L.Warn("LoadSerCode:g.DB.Query", zap.Error(err), zap.String("query", query))
			return 0, err
		}
		defer rows.Close()

		rows.Next()
		rows.Scan(&id)

		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		g.L.Warn("LoadAppCode:result.LastInsertId", zap.Error(err))
		return 0, err
	}

	serInfo.SerID = int32(id)
	// as.SerNames[serInfo.AppCode] = serInfo

	return 0, nil
}
