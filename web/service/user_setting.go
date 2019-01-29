package service

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"go.uber.org/zap"
)

func (web *Web) setUser(c echo.Context) error {
	appNameS := c.FormValue("app_names")
	appShow, _ := strconv.Atoi(c.FormValue("app_show"))

	li := web.getLoginInfo(c)
	q := web.Cql.Query(`UPDATE  account SET app_show=?,app_names=? WHERE id=?`, appShow, appNameS, li.ID)
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

type AppSetting struct {
	AppShow  int    `json:"app_show"`
	AppNames string `json:"app_names"`
}

func (web *Web) getAppSetting(c echo.Context) error {
	li := web.getLoginInfo(c)

	var appNames string
	var appShow int
	q := web.Cql.Query(`SELECT app_show,app_names from  account  WHERE id=?`, li.ID)
	err := q.Scan(&appShow, &appNames)
	if err != nil {
		g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: g.DatabaseC,
			Message: g.DatabaseE,
		})
	}

	as := AppSetting{appShow, appNames}
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   as,
	})
}
