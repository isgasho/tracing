package service

import (
	"net/http"

	"github.com/imdevlab/g"
	"github.com/labstack/echo"
)

func queryServiceMap(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   serviceMapTest,
	})
}
