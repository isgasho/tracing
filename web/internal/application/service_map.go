package app

import (
	"net/http"

	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/web/internal/tests"
	"github.com/labstack/echo"
)

func QueryServiceMap(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   tests.ServiceMapData,
	})
}
