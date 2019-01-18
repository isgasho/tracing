package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
	"github.com/mafanr/vgo/util"
	"go.uber.org/zap"
)

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

	// 插入
	q := web.Cql.Query(`UPDATE  alerts_policy SET name=?,alerts=?,update_date=? WHERE id=? and owner=?`, policy.Name, policy.Alerts, time.Now().Unix(), policy.ID, li.ID)
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

	// 判断group是否已经存在
	q := `SELECT name from alerts_group WHERE name=?`
	count := web.Cql.Query(q, name).Iter().NumRows()
	if count > 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.AlreadyExistC,
			Message: g.AlreadyExistE,
		})
	}

	// 插入
	q1 := web.Cql.Query(`INSERT INTO  alerts_group (name,owner,channel,users,update_date) VALUES (?,?,?,?,?)`, name, li.ID, channel, users, utils.Time2StringSecond(time.Now()))
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

	// 更新
	q1 := web.Cql.Query(`UPDATE alerts_group SET channel=?,users=?,update_date=? WHERE name=? and owner=?`,
		channel, users, utils.Time2StringSecond(time.Now()), name, li.ID)
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

	// 删除
	q1 := web.Cql.Query(`DELETE FROM  alerts_group WHERE name=? and owner=?`,
		name, li.ID)
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
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Channel string   `json:"channel"`
	Users   []string `json:"users"`
}

func (web *Web) queryGroups(c echo.Context) error {
	// 获取当前用户
	li := web.getLoginInfo(c)

	// 若该用户是管理员，可以获取所有组
	var iter *gocql.Iter
	if li.Priv == g.PRIV_NORMAL {
		iter = web.Cql.Query(`SELECT name,owner,channel,users FROM alerts_group WHERE owner=?`, li.ID).Iter()
	} else {
		iter = web.Cql.Query(`SELECT name,owner,channel,users FROM alerts_group`).Iter()
	}

	var name, owner, channel string
	var users []string

	groups := make([]*Group, 0)
	for iter.Scan(&name, &owner, &channel, &users) {
		groups = append(groups, &Group{name, owner, channel, users})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   groups,
	})
}
