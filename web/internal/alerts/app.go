package alerts

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/imdevlab/g"
	app "github.com/imdevlab/tracing/web/internal/application"
	ecode "github.com/imdevlab/tracing/web/internal/error_code"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/imdevlab/tracing/web/internal/session"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

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

func CreateApp(c echo.Context) error {
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
	li := session.GetLoginInfo(c)

	// 查询目标appAlert是否已经存在
	n := misc.Cql.Query(`SELECT name FROM  alerts_app WHERE name=?`, appName).Iter().NumRows()
	if n > 0 {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusInternalServerError,
			ErrCode: ecode.AppAlertNameExistC,
			Message: ecode.AppAlertNameExistE,
		})
	}

	// 插入
	q := misc.Cql.Query(`INSERT INTO  alerts_app (name,owner,policy_id,channel,users,update_date) VALUES (?,?,?,?,?,?)`, appName, li.ID, policy, channel, users, time.Now().Unix())
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

func EditApp(c echo.Context) error {
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
	li := session.GetLoginInfo(c)

	// 查询目标appAlert是否已经存在
	var owner string
	err = misc.Cql.Query(`SELECT owner FROM  alerts_app WHERE name=?`, appName).Scan(&owner)
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
	q := misc.Cql.Query(`UPDATE alerts_app SET policy_id=?,channel=?,users=?,update_date=? WHERE name=? and owner=? IF EXISTS`, policy, channel, users, time.Now().Unix(), appName, owner)
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

func DeleteApp(c echo.Context) error {
	name := c.FormValue("name")

	if name == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	// 获取当前用户
	li := session.GetLoginInfo(c)

	// 查询目标appAlert是否已经存在
	var owner string
	err := misc.Cql.Query(`SELECT owner FROM  alerts_app WHERE name=?`, name).Scan(&owner)
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
	q := misc.Cql.Query(`DELETE FROM  alerts_app WHERE name=? and owner=?`,
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

func AppList(c echo.Context) error {
	tp := c.FormValue("type")
	// 获取当前用户
	li := session.GetLoginInfo(c)

	apps := make([]*AppAlert, 0)
	polies := make(map[string]string)
	var userNames []string
	var name, owner, channel, policy string
	var users []string

	var iter *gocql.Iter
	switch tp {
	case "1": // 查看全部应用告警
		iter = misc.Cql.Query(`SELECT name,owner,policy_id,channel,users FROM alerts_app`).Iter()
	case "2": // 用户自己创建的
		iter = misc.Cql.Query(`SELECT name,owner,policy_id,channel,users FROM alerts_app WHERE owner=?`, li.ID).Iter()
	case "3": // 用户设定的应用列表
		_, appNames := app.UserSetting(li.ID)
		for _, an := range appNames {
			q := misc.Cql.Query(`SELECT name,owner,policy_id,channel,users FROM alerts_app WHERE name=?`, an)
			err := q.Scan(&name, &owner, &policy, &channel, &users)
			if err != nil {
				g.L.Warn("query database error:", zap.Error(err))
				continue
			}
			var on string
			ownerNameR, ok := session.UsersMap.Load(owner)
			if ok {
				on = ownerNameR.(*session.User).Name
			}
			for _, u := range users {
				var un string
				unr, ok := session.UsersMap.Load(u)
				if ok {
					un = unr.(*session.User).Name
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
			ownerNameR, ok := session.UsersMap.Load(owner)
			if ok {
				on = ownerNameR.(*session.User).Name
			}
			for _, u := range users {
				var un string
				unr, ok := session.UsersMap.Load(u)
				if ok {
					un = unr.(*session.User).Name
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
		misc.Cql.Query(`SELECT name FROM alerts_policy WHERE id = ?`, pid).Scan(&pname)
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
