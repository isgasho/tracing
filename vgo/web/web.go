package web

import (
	"net/http"

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
}

// New ...
func New() *Web {
	return &Web{}
}

// Start ...
func (s *Web) Start() error {
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
