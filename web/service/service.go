package service

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/web/misc"
	"go.uber.org/zap"
)

// 后台服务
// Stats 离线计算
type Web struct {
	Cql      *gocql.Session
	cache    *cache
	sessions *sync.Map
}

// New ...
func New() *Web {
	return &Web{}
}

type RetryPolicy struct {
	NumRetries int //Number of times to retry a query,-1 means always retries
	Interval   int
}

// Attempt tells gocql to attempt the query again based on query.Attempts being less
// than the NumRetries defined in the policy.
func (s *RetryPolicy) Attempt(q gocql.RetryableQuery) bool {
	fmt.Println("start retry")
	time.Sleep(time.Duration(s.Interval) * time.Second)
	if s.NumRetries == -1 {
		return true
	}

	return q.Attempts() <= s.NumRetries
}

func (s *RetryPolicy) GetRetryType(err error) gocql.RetryType {
	return gocql.Retry
}

// Start ...
func (s *Web) Start() error {
	// 用户登录session
	s.sessions = &sync.Map{}

	// 初始化内部缓存
	s.cache = &cache{}
	// 初始化Cql连接
	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
	cluster.Keyspace = misc.Conf.Storage.Keyspace
	cluster.Timeout = 5 * time.Second

	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = 20

	// cluster.RetryPolicy = &RetryPolicy{NumRetries: -1, Interval: 2}

	session, err := cluster.CreateSession()
	if err != nil {
		g.L.Fatal("Init web cql connections error", zap.String("error", err.Error()))
		return err
	}
	s.Cql = session

	go func() {
		e := echo.New()
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowHeaders:     append([]string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}, "X-Token"),
			AllowCredentials: true,
		}))

		e.Pre(middleware.RemoveTrailingSlash())
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))

		// 回调相关
		//同步回调接口
		e.POST("/apm/web/login", s.login)
		e.POST("/apm/web/logout", s.logout)
		// 应用查询接口
		//查询应用列表
		e.GET("/apm/web/appList", s.appList, s.checkLogin)
		//获取指定应用的一段时间内的状态
		e.GET("/apm/web/appDash", s.appDash, s.checkLogin)
		//查询所有的app名
		e.GET("/apm/web/appNames", s.appNames, s.checkLogin)
		//查询所有服务器名
		e.GET("/apm/web/agentList", s.agentList, s.checkLogin)

		e.GET("/apm/web/serviceMap", queryServiceMap, s.checkLogin)
		e.GET("/apm/web/traces", queryTraces, s.checkLogin)
		e.GET("/apm/web/trace", queryTrace, s.checkLogin)

		// 管理员面板
		e.GET("/apm/web/userList", s.userList, s.checkLogin)
		e.POST("/apm/web/setAdmin", s.setAdmin, s.checkLogin)
		e.POST("/apm/web/cancelAdmin", s.cancelAdmin, s.checkLogin)

		e.Logger.Fatal(e.Start(misc.Conf.Web.Addr))
	}()

	return nil
}

// Close ...
func (s *Web) Close() error {
	return nil
}

func (web *Web) checkLogin(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		li := web.getLoginInfo(c)
		if li == nil {
			return c.JSON(http.StatusOK, g.Result{
				Status:  http.StatusUnauthorized,
				ErrCode: g.NeedLoginC,
				Message: g.NeedLoginE,
			})
		}

		return f(c)
	}
}
