package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
	"go.uber.org/zap"
)

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
	if err != nil || len(users) == 0 {
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
	if err != nil || len(users) == 0 {
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
