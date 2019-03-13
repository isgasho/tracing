package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/util"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func (web *Web) createAppAlert(c echo.Context) error {
	appName := c.FormValue("app_name")
	policy := c.FormValue("policy")
	channel := c.FormValue("channel")
	usersS := c.FormValue("users")

	if appName == "" || policy == "" || channel == "" || usersS == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	var users []string
	err := json.Unmarshal([]byte(usersS), &users)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 查询目标appAlert是否已经存在
	n := web.Cql.Query(`SELECT name FROM  app_alert WHERE name=?`, appName).Iter().NumRows()
	if n > 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: AppAlertNameExistC,
			Message: AppAlertNameExistE,
		})
	}

	// 插入
	q := web.Cql.Query(`INSERT INTO  app_alert (name,owner,policy,channel,users,update_date) VALUES (?,?,?,?,?,?)`, appName, li.ID, policy, channel, users, time.Now().Unix())
	err = q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})

}

func (web *Web) editAppAlert(c echo.Context) error {
	appName := c.FormValue("app_name")
	policy := c.FormValue("policy")
	channel := c.FormValue("channel")
	usersS := c.FormValue("users")

	if appName == "" || policy == "" || channel == "" || usersS == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	var users []string
	err := json.Unmarshal([]byte(usersS), &users)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 查询目标appAlert是否已经存在
	var owner string
	err = web.Cql.Query(`SELECT owner FROM  app_alert WHERE name=?`, appName).Scan(&owner)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	// 必须是owner本人或者管理员才能编辑
	if owner != li.ID && li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// 插入
	q := web.Cql.Query(`UPDATE app_alert SET policy=?,channel=?,users=?,update_date=? WHERE name=? and owner=? IF EXISTS`, policy, channel, users, time.Now().Unix(), appName, owner)
	err = q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})

}

func (web *Web) deleteAppAlert(c echo.Context) error {
	name := c.FormValue("name")

	if name == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 查询目标appAlert是否已经存在
	var owner string
	err := web.Cql.Query(`SELECT owner FROM  app_alert WHERE name=?`, name).Scan(&owner)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	// 必须是owner本人才能编辑
	if owner != li.ID {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// 删除
	q := web.Cql.Query(`DELETE FROM  app_alert WHERE name=? and owner=?`,
		name, li.ID)
	err = q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

type AppAlert struct {
	Name       string   `json:"name"`
	OwnerID    string   `json:"owner_id"`
	OwnerName  string   `json:"owner_name"`
	Policy     string   `json:"policy"`
	PolicyName string   `json:"policy_name"`
	Channel    string   `json:"channel"`
	Users      []string `json:"users"`
	UserNames  []string `json:"user_names"`
}

func (web *Web) alertsAppList(c echo.Context) error {
	tp := c.FormValue("type")
	// 获取当前用户
	li := web.getLoginInfo(c)

	apps := make([]*AppAlert, 0)
	polies := make(map[string]string)
	var userNames []string
	var name, owner, channel, policy string
	var users []string

	var iter *gocql.Iter
	switch tp {
	case "1": // 查看全部应用告警
		iter = web.Cql.Query(`SELECT name,owner,policy,channel,users FROM app_alert`).Iter()
	case "2": // 用户自己创建的
		iter = web.Cql.Query(`SELECT name,owner,policy,channel,users FROM app_alert WHERE owner=?`, li.ID).Iter()
	case "3": // 用户设定的应用列表
		_, appNames := web.userAppSetting(li.ID)
		for _, an := range appNames {
			q := web.Cql.Query(`SELECT name,owner,policy,channel,users FROM app_alert WHERE name=?`, an)
			err := q.Scan(&name, &owner, &policy, &channel, &users)
			if err != nil {
				g.L.Warn("query database error:", zap.Error(err))
				continue
			}
			var on string
			ownerNameR, ok := web.usersMap.Load(owner)
			if ok {
				on = ownerNameR.(*User).Name
			}
			for _, u := range users {
				var un string
				unr, ok := web.usersMap.Load(u)
				if ok {
					un = unr.(*User).Name
				}
				userNames = append(userNames, un)
			}
			apps = append(apps, &AppAlert{name, owner, on, policy, "", channel, users, userNames})
			polies[policy] = ""
		}
	}

	if tp == "1" || tp == "2" {
		for iter.Scan(&name, &owner, &policy, &channel, &users) {
			var on string
			ownerNameR, ok := web.usersMap.Load(owner)
			if ok {
				on = ownerNameR.(*User).Name
			}
			for _, u := range users {
				var un string
				unr, ok := web.usersMap.Load(u)
				if ok {
					un = unr.(*User).Name
				}
				userNames = append(userNames, un)
			}
			apps = append(apps, &AppAlert{name, owner, on, policy, "", channel, users, userNames})
			polies[policy] = ""
		}

		if err := iter.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}
	}

	for pid := range polies {
		var pname string
		web.Cql.Query(`SELECT name FROM alerts_policy WHERE id = ?`, pid).Scan(&pname)
		for _, app := range apps {
			if app.Policy == pid {
				app.PolicyName = pname
			}
		}
	}

	// 查询policy id对应的名称
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   apps,
	})
}

type Policy struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	OwnerID   string        `json:"owner_id"`
	OwnerName string        `json:"owner_name"`
	Alerts    []*util.Alert `json:"alerts"`
}

func (web *Web) createPolicy(c echo.Context) error {
	policyRaw := c.FormValue("policy")

	policy := &Policy{}
	err := json.Unmarshal([]byte(policyRaw), &policy)
	if err != nil || policy.Name == "" || len(policy.Alerts) == 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 同一个用户下的模版名不能重复
	n := web.Cql.Query(`SELECT id FROM  alerts_policy WHERE name=? and owner=? ALLOW FILTERING`, policy.Name, li.ID).Iter().NumRows()
	if n > 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: PolicyNameExistC,
			Message: PolicyNameExistE,
		})
	}
	// 插入
	q := web.Cql.Query(`INSERT INTO  alerts_policy (id,name,owner,alerts,update_date) VALUES (uuid(),?,?,?,?)`, policy.Name, li.ID, policy.Alerts, time.Now().Unix())
	err = q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (web *Web) editPolicy(c echo.Context) error {
	policyRaw := c.FormValue("policy")

	policy := &Policy{}
	err := json.Unmarshal([]byte(policyRaw), &policy)
	if err != nil || policy.Name == "" || len(policy.Alerts) == 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 查询目标policy是否已经存在
	var owner string
	err = web.Cql.Query(`SELECT owner FROM  alerts_policy WHERE id=?`, policy.ID).Scan(&owner)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	// 必须是owner本人或者管理员才能编辑
	if owner != li.ID && li.Priv != g.PRIV_SUPER_ADMIN && li.Priv != g.PRIV_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// 插入
	q := web.Cql.Query(`UPDATE  alerts_policy SET name=?,alerts=?,update_date=? WHERE id=? and owner=? IF EXISTS`, policy.Name, policy.Alerts, time.Now().Unix(), policy.ID, owner)
	err = q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (web *Web) queryPolicies(c echo.Context) error {
	// 获取当前用户
	li := web.getLoginInfo(c)

	// 若该用户是管理员，可以获取所有组
	var iter *gocql.Iter
	if li.Priv == g.PRIV_NORMAL {
		iter = web.Cql.Query(`SELECT id,name,owner,alerts FROM alerts_policy WHERE owner=?`, li.ID).Iter()
	} else {
		iter = web.Cql.Query(`SELECT id,name,owner,alerts FROM alerts_policy`).Iter()
	}

	var id, name, owner string
	var alerts []*util.Alert

	policies := make([]*Policy, 0)
	for iter.Scan(&id, &name, &owner, &alerts) {
		ownerNameR, _ := web.usersMap.Load(owner)
		policies = append(policies, &Policy{id, name, owner, ownerNameR.(*User).Name, alerts})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   policies,
	})
}

func (web *Web) queryPolicy(c echo.Context) error {
	pid := c.FormValue("id")

	var id, name, owner string
	var alerts []*util.Alert

	q := web.Cql.Query(`SELECT id,name,owner,alerts FROM alerts_policy WHERE id=?`, pid)
	err := q.Scan(&id, &name, &owner, &alerts)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	ownerNameR, _ := web.usersMap.Load(owner)
	policy := &Policy{id, name, owner, ownerNameR.(*User).Name, alerts}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   policy,
	})
}

func (web *Web) deletePolicy(c echo.Context) error {
	id := c.FormValue("id")

	if id == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 删除
	q := web.Cql.Query(`DELETE FROM  alerts_policy WHERE id=? and owner=?`,
		id, li.ID)
	err := q.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (web *Web) createGroup(c echo.Context) error {
	name := c.FormValue("name")
	channel := c.FormValue("channel")
	usersS := c.FormValue("users")

	if name == "" || channel == "" || usersS == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	var users []string
	err := json.Unmarshal([]byte(usersS), &users)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 插入
	q1 := web.Cql.Query(`INSERT INTO  alerts_group (id,name,owner,channel,users,update_date) VALUES (uuid(),?,?,?,?,?)`, name, li.ID, channel, users, time.Now().Unix())
	err = q1.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q1.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (web *Web) editGroup(c echo.Context) error {
	id := c.FormValue("id")
	name := c.FormValue("name")
	channel := c.FormValue("channel")
	usersS := c.FormValue("users")

	if id == "" || name == "" || channel == "" || usersS == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	var users []string
	err := json.Unmarshal([]byte(usersS), &users)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 更新
	q1 := web.Cql.Query(`UPDATE alerts_group SET name=?,channel=?,users=?,update_date=? WHERE id=? and owner=?`,
		name, channel, users, time.Now().Unix(), id, li.ID)
	err = q1.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q1.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

func (web *Web) deleteGroup(c echo.Context) error {
	id := c.FormValue("id")

	if id == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := web.getLoginInfo(c)

	// 删除
	q1 := web.Cql.Query(`DELETE FROM  alerts_group WHERE id=? and owner=?`,
		id, li.ID)
	err := q1.Exec()
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q1.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
	})
}

type Group struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	OwnerID   string   `json:"owner_id"`
	OwnerName string   `json:"owner_name"`
	Channel   string   `json:"channel"`
	Users     []string `json:"users"`
}

func (web *Web) queryGroups(c echo.Context) error {
	// 获取当前用户
	li := web.getLoginInfo(c)

	// 若该用户是管理员，可以获取所有组
	var iter *gocql.Iter
	if li.Priv == g.PRIV_NORMAL {
		iter = web.Cql.Query(`SELECT id,name,owner,channel,users FROM alerts_group WHERE owner=?`, li.ID).Iter()
	} else {
		iter = web.Cql.Query(`SELECT id,name,owner,channel,users FROM alerts_group`).Iter()
	}

	var id, name, owner, channel string
	var users []string

	groups := make([]*Group, 0)
	for iter.Scan(&id, &name, &owner, &channel, &users) {
		ownerNameR, _ := web.usersMap.Load(owner)
		groups = append(groups, &Group{id, name, owner, ownerNameR.(*User).Name, channel, users})
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   groups,
	})
}
