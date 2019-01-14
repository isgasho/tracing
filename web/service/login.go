package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"time"

	"github.com/mafanr/vgo/web/misc"

	"go.uber.org/zap"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/valyala/fasthttp"
)

var defaultRole = "normal"

type Session struct {
	User       *UserInfo
	CreateTime time.Time
}

type UserInfo struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Avatar   string   `json:"avatar"`
	Role     []string `json:"role"`
	SsoToken string   `json:"ssoToken"`
}

func (web *Web) login(c echo.Context) error {
	subToken := c.FormValue("subToken")
	// 通过subtoken获取用户和ssotoken
	body := "{'subToken':'" + subToken + "'}"
	url := misc.Conf.Login.SsoLogin

	b := requestToSso(body, url)
	fmt.Println(string(b))
	tokenInfo := &TokenInfo{}
	err := json.Unmarshal(b, tokenInfo)
	if err != nil {
		g.L.Info("解析sso用户信息失败", zap.Error(err), zap.String("body", string(b)))
		return nil
	}

	uid := tokenInfo.Data.SubTokenObj.UserSession.UserId

	user := &UserInfo{
		ID:       uid,
		Name:     tokenInfo.Data.SubTokenObj.UserSession.UserName,
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		SsoToken: tokenInfo.Data.SubTokenObj.SsoToken,
	}

	if user.ID == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			Message: "sub token无效，请重新登陆",
		})
	}

	//sub token验证成功，保存session
	web.sessions.Store(user.SsoToken, &Session{
		User:       user,
		CreateTime: time.Now(),
	})

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (a *Web) loginMock(c echo.Context) error {
	user := &UserInfo{
		ID:       "13269",
		Name:     "孙飞",
		Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		SsoToken: "0af8c18eed353af38b2e2524f4850f76",
	}

	//sub token验证成功，保存session
	a.sessions.Store(user.SsoToken, &Session{
		User:       user,
		CreateTime: time.Now(),
	})

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   user,
	})
}

func (a *Web) logout(c echo.Context) error {
	token := getToken(c)
	// 删除用户的session
	a.sessions.Delete(token)

	// 请求sso 注销token
	body := "{'ssoToken':'" + token + "', 'clientNo':'1' }"
	url := misc.Conf.Login.SsoLogout
	b := requestToSso(body, url)
	fmt.Println(string(b))
	return c.JSON(http.StatusOK, g.Result{
		Status:  http.StatusOK,
		Message: "退出登陆成功",
	})
}

func (a *Web) userInfo(c echo.Context) error {
	token := getToken(c)
	sess, ok := a.sessions.Load(token)
	if !ok {
		// 用户未登陆或者session失效
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusUnauthorized,
			ErrCode: g.NeedLoginC,
			Message: g.NeedLoginE,
		})
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   sess.(*Session).User,
	})
}

//subToken认证，返回的用户信息
type TokenInfo struct {
	Code int `json:"code"`
	Data struct {
		SubTokenObj struct {
			Message     string `json:"message"`
			SsoToken    string `json:"ssoToken"`
			Status      int    `json:"status"`
			UserSession struct {
				Facility        string `json:"facility"`
				HeadImgUrl      string `json:"headImgUrl"`
				SysDepartmentId string `json:"sysDepartmentId"`
				UserId          string `json:"userId"`
				UserName        string `json:"userName"`
				UserType        string `json:"userType"`
			} `json:"userSession"`
		} `json:"subTokenObj"`
	} `json:"data"`
	Message string `json:"message"`
}

func requestToSso(body string, url string) []byte {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBodyString(body)

	req.SetRequestURI(url)
	var cli = &fasthttp.Client{}
	err := cli.DoTimeout(req, resp, 15*time.Second)
	if err != nil {
		g.L.Info("获取sso用户信息失败", zap.Error(err))
		return nil
	}

	return resp.Body()
}

func getToken(c echo.Context) string {
	return c.Request().Header.Get("X-Token")
}
