package web

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/vgo/misc"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"
)

// 后台服务
// Stats 离线计算
type Web struct {
	Cql *gocql.Session
}

// New ...
func New() *Web {
	return &Web{}
}

// Start ...
func (s *Web) Start() error {
	// 初始化Cql连接
	// connect to the cluster
	cluster := gocql.NewCluster(misc.Conf.Storage.Cluster...)
	cluster.Keyspace = misc.Conf.Storage.Keyspace
	cluster.Consistency = gocql.Quorum
	//设置连接池的数量,默认是2个（针对每一个host,都建立起NumConns个连接）
	cluster.NumConns = misc.Conf.Storage.NumConns

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

		// 回调相关
		//同步回调接口
		e.POST("/login", func(c echo.Context) error {
			return c.JSON(http.StatusOK, g.Result{
				Status: http.StatusOK,
				Data:   "hello login",
			})
		})

		// 应用查询接口
		//查询应用列表
		e.GET("/apm/query/appList", s.appList)
		//根据搜索条件，查询应用

		e.GET("/apm/query/serviceMap", queryServiceMap)
		e.GET("/apm/query/traces", queryTraces)
		e.GET("/apm/query/trace", queryTrace)

		e.Logger.Fatal(e.Start(misc.Conf.Web.Addr))
	}()

	go func() {
		config := newrelic.NewConfig("vgo-web", "fb07ffe86ca0ab409e91faae48e13b253e87be23")
		app, err := newrelic.NewApplication(config)
		if err != nil {
			g.L.Fatal("start relic apm error", zap.Error(err))
		}
		http.HandleFunc(newrelic.WrapHandleFunc(app, "/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello world"))
		}))
		http.ListenAndServe(":8001", nil)
	}()

	return nil
}

// Close ...
func (s *Web) Close() error {
	return nil
}
