package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
)

func queryServiceMap(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   serviceMapTest,
	})
}
