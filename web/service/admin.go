package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"go.uber.org/zap"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
	Priv   string `json:"priv"`
}

func (web *Web) userList(c echo.Context) error {
	// 查询所有用户
	q := `SELECT id,name,mobile,email FROM account`
	iter := web.Cql.Query(q).Iter()

	users := make(map[string]*User)
	var id, name, mobile, email, priv string

	for iter.Scan(&id, &name, &mobile, &email) {
		users[id] = &User{
			ID:     id,
			Name:   name,
			Mobile: mobile,
			Email:  email,
			Priv:   g.PRIV_NORMAL,
		}
	}

	// 查询相应权限
	q = `SELECT id,priv FROM admin`
	iter = web.Cql.Query(q).Iter()

	for iter.Scan(&id, &priv) {
		users[id].Priv = priv
	}

	nusers := make([]*User, 0)
	for _, u := range users {
		nusers = append(nusers, u)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   nusers,
	})
}

func (web *Web) setAdmin(c echo.Context) error {
	userID := c.FormValue("user_id")
	li := web.getLoginInfo(c)
	// 判断当前用户是否超级管理员
	if li.Priv != g.PRIV_SUPER_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// 将目标用户设置为管理员
	q := web.Cql.Query(`INSERT INTO admin (id,priv) VALUES (?,?)`, userID, g.PRIV_ADMIN)
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

func (web *Web) cancelAdmin(c echo.Context) error {
	userID := c.FormValue("user_id")
	li := web.getLoginInfo(c)
	// 判断当前用户是否超级管理员
	if li.Priv != g.PRIV_SUPER_ADMIN {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusForbidden,
			ErrCode: g.ForbiddenC,
			Message: g.ForbiddenE,
		})
	}

	// 将目标用户设置为管理员
	q := web.Cql.Query(`DELETE FROM admin WHERE id=?`, userID)
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
