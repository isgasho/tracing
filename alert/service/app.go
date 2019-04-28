package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/imdevlab/tracing/alert/misc"
	"github.com/imdevlab/tracing/pkg/alert"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/pkg/sql"
	"github.com/imdevlab/tracing/pkg/util"

	"go.uber.org/zap"
)

// Apps apps
type Apps struct {
	sync.RWMutex
	Apps map[string]*App
}

func newApps() *Apps {
	return &Apps{
		Apps: make(map[string]*App),
	}
}

func (a *Apps) start() error {
	if err := a.loadPolicy(); err != nil {
		logger.Warn("load policy error", zap.String("error", err.Error()))
		return err
	}

	go func() {
		time.Sleep(time.Duration(misc.Conf.App.LoadInterval) * time.Second)
		if err := a.loadPolicy(); err != nil {
			logger.Warn("load policy error", zap.String("error", err.Error()))
		}
	}()

	return nil
}

func (a *Apps) loadPolicy() error {
	cql := gAlert.GetCql()
	if cql == nil {
		return fmt.Errorf("unfind cql")
	}
	query := cql.Query(sql.LoadPolicys).Iter()
	defer func() {
		if err := query.Close(); err != nil {
			logger.Warn("close iter error:", zap.Error(err))
		}
	}()

	var name, owner, apiAlertsStr, channel, group, policyID string
	var users []string
	var updateDate int64

	checkTime := time.Now().Unix()
	for query.Scan(&name, &owner, &apiAlertsStr, &channel, &group, &policyID, &updateDate, &users) {
		a.RLock()
		app, ok := a.Apps[name]
		a.RUnlock()
		// 如果已经存在策略并且updatedate不相当，那么删除历史
		if ok {
			if app.policy.UpdateDate == updateDate {
				app.policy.checkTime = checkTime
				continue
			}
			// 定时任务移除
			gAlert.tickers.RemoveTask(app.taskID)
			// 策略被更新，需要删除
			a.remove(name)
		}

		var tmpapiAlerts []*util.ApiAlert
		if err := json.Unmarshal([]byte(apiAlertsStr), &tmpapiAlerts); err != nil {
			logger.Warn("json Unmarshal", zap.String("error", err.Error()))
			continue
		}

		app = newApp()
		app.policy.AppName = name
		app.policy.Owner = owner
		app.policy.Channel = channel
		app.policy.Group = group
		app.policy.ID = policyID
		app.policy.UpdateDate = updateDate
		app.policy.Users = users
		app.policy.checkTime = checkTime

		// 根据alertid加具体载策略,如果policyID为null那么代表该模版的策略被删除，所以不用统计
		if len(policyID) == 0 {
			continue
		}

		alertsQuery := cql.Query(sql.LoadAlert, policyID)
		var tmpAlerts []*util.Alert
		if err := alertsQuery.Scan(&tmpAlerts); err != nil {
			logger.Warn("load alert scan error", zap.String("error", err.Error()), zap.String("sql", sql.LoadAlert))
			continue
		}
		if len(tmpAlerts) == 0 {
			continue
		}

		for _, tmpAlert := range tmpAlerts {
			alert := newAlertInfo()
			alert.Compare = tmpAlert.Compare
			alert.Duration = tmpAlert.Duration
			alert.Value = tmpAlert.Value
			alert.Keys = strings.Split(tmpAlert.Keys, ",")
			alertType, ok := constant.AlertType(tmpAlert.Name)
			if !ok {
				logger.Warn("alertType unfind error", zap.String("name", tmpAlert.Name))
				continue
			}
			alert.Type = alertType
			app.Alerts[alertType] = alert
		}

		// 加载特殊监控
		app.loadAPIAlerts(tmpapiAlerts)

		// app start
		if err := app.start(); err != nil {
			logger.Warn("app start", zap.String("error", err.Error()), zap.String("appName", name))
			continue
		}

		taskID := gAlert.tickers.NewID()
		app.taskID = taskID
		gAlert.tickers.AddTask(taskID, app.tChan)
		// 保存策略
		a.add(app)
	}

	// 对比模版checktime，发现checktime不相等，那么代表该模版已经被删除
	a.checkVersion(checkTime)
	return nil
}

func (a *Apps) checkVersion(checkTime int64) {
	a.Lock()
	for name, app := range a.Apps {
		if app.policy.checkTime != checkTime {
			delete(a.Apps, name)
		}
	}
	a.Unlock()
}

func (a *Apps) remove(name string) {
	a.RLock()
	app, ok := a.Apps[name]
	a.RUnlock()
	if !ok {
		return
	}

	app.close()

	a.Lock()
	delete(a.Apps, name)
	a.Unlock()
}

func (a *Apps) add(app *App) {
	a.Lock()
	a.Apps[app.policy.AppName] = app
	a.Unlock()
}

// App app
type App struct {
	policy       *Policy            // Policy
	SpecialAlert *SpecialAlert      // 特殊监控
	Alerts       map[int]*AlertInfo // 策略模版，通用策略
	tChan        chan bool          // 任务channel
	taskID       int64              // 任务ID
	stopC        chan bool          // stop chan
	apiC         chan *alert.API    // api 数据通道
}

func newApp() *App {
	return &App{
		policy:       newPolicy(),
		SpecialAlert: newSpecialAlert(),
		Alerts:       make(map[int]*AlertInfo),
		tChan:        make(chan bool, 50),
		stopC:        make(chan bool, 1),
		apiC:         make(chan *alert.API, 50),
	}
}

func (a *App) start() error {
	go a.analyze()
	return nil
}

func (a *App) analyze() {
	for {
		select {
		case <-a.stopC:
			return
		case _, ok := <-a.tChan:
			if ok {
				logger.Info("analyze", zap.String("name", a.policy.AppName), zap.Any("SpecialAlert", a.SpecialAlert))
			}
			break
		case api, ok := <-a.apiC:
			if ok {
				logger.Info("api", zap.Any("api", api))
			}
			break
		}
	}
}

func (a *App) close() error {
	close(a.tChan)
	close(a.stopC)
	return nil
}

func (a *App) loadAPIAlerts(tmpapiAlerts []*util.ApiAlert) {
	for _, tmpAPIAlert := range tmpapiAlerts {
		for _, tmpalert := range tmpAPIAlert.Alerts {
			alertType, ok := constant.AlertType(tmpalert.Key)
			if !ok {
				logger.Warn("alertType unfind error", zap.String("name", tmpalert.Key))
				continue
			}
			// 到通用alert列表中去找同类型的alert,然后复用该alert里面的值，找不到直接返回
			universalAlert, ok := a.Alerts[alertType]
			if !ok {
				logger.Warn("alert unfind error", zap.Int("alertType", alertType))
				continue
			}

			// 创建特殊alert，然后赋值（复用universalAlert中的值）
			specialAlert, ok := a.SpecialAlert.API[tmpAPIAlert.Api]
			if !ok {
				specialAlert = newAlertInfo()
				a.SpecialAlert.API[tmpAPIAlert.Api] = specialAlert
			}
			specialAlert.Type = universalAlert.Type
			specialAlert.Compare = universalAlert.Compare
			specialAlert.Duration = universalAlert.Duration
			specialAlert.Value = tmpalert.Value
		}
	}
}

// // Write recv data
// func (a *Analyze) Write(data *alert.Data) error {
// 	switch data.Type {
// 	case constant.ALERT_TYPE_API:
// 		break
// 	case constant.ALERT_TYPE_SQL:
// 		break
// 	}
// 	return nil
// }
