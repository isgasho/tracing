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
