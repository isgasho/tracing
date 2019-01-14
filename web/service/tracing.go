package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
)

func queryTraces(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   tracesTest,
	})
}

func queryTrace(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   traceTest,
	})
}
