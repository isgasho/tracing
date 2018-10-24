package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/mafanr/g"
	"go.uber.org/zap"
)

// AppStore ...
type AppStore struct {
	sync.RWMutex
	Apps map[string]*App
}

// NewAppStore ...
func NewAppStore() *AppStore {
	return &AppStore{
		Apps: make(map[string]*App),
	}
}

// LoadApps 加载数据库中的所有app
func (as *AppStore) LoadApps() error {
	// 加载所有appCode
	apps := make([]*App, 0)
	err := g.DB.Select(&apps, "select * from app")
	if err != nil {
		g.L.Fatal("LoadApps:g.DB.Select", zap.Error(err))
	}

	for _, app := range apps {
		as.Apps[app.Name] = app
	}
	return nil
}

// LoadAppCode  通过appname获取code
func (as *AppStore) LoadAppCode(name string) (int, error) {
	as.RLock()
	app, ok := as.Apps[name]
	as.RUnlock()
	if ok {
		return app.Code, nil
	}
	apps := make([]*App, 0)
	err := g.DB.Select(&apps, fmt.Sprintf("select * from app where name = '%s'", name))
	if err != nil {
		g.L.Warn("LoadAppCode:g.DB.Select", zap.Error(err))
		return 0, err
	}

	log.Println("1111111111", "apps ", apps)

	if len(apps) > 0 {
		// 缓存到数据库中
		as.Lock()
		as.Apps[name] = apps[0]
		as.Unlock()
		log.Println("1111111111", apps[0])
		return apps[0].Code, nil
	}

	// 如果没插入
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
	app = &App{
		Name: name,
		Code: int(code),
	}

	// 缓存到数据库中
	as.Lock()
	as.Apps[name] = app
	as.Unlock()

	log.Println("222121233412321212", app)
	return int(code), nil
}

// App ...
type App struct {
	Code int    `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}

// NewApp ...
func NewApp() *App {
	return &App{}
}
