package web

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
)

type AppStat struct {
	Name           string  `json:"name"`
	Count          int     `json:"count"`
	Apdex          float64 `json:"apdex"`
	AverageElapsed float64 `json:"average_elapsed"`
	ErrorPercent   float64 `json:"error_percent"`

	totalElapsed float64

	errCount     float64
	satisfaction float64
	tolerate     float64
}

func (web *Web) appList(c echo.Context) error {
	now := time.Now().Unix()
	// 取过去6分钟的数据
	start := now - 450
	q := `SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM rpc_stats WHERE input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, start, now).Iter()
	defer iter.Close()

	apps := make(map[string]*AppStat)
	var appName string
	var count int
	var tElapsed, errCount, satisfaction, tolerate int

	for iter.Scan(&appName, &tElapsed, &count, &errCount, &satisfaction, &tolerate) {
		app, ok := apps[appName]
		if !ok {
			apps[appName] = &AppStat{
				Name:         appName,
				Count:        count,
				totalElapsed: float64(tElapsed),
				errCount:     float64(errCount),
				satisfaction: float64(satisfaction),
				tolerate:     float64(tolerate),
			}
		} else {
			app.Count += count
			app.totalElapsed += float64(tElapsed)
			app.errCount += float64(errCount)
			app.satisfaction += float64(satisfaction)
			app.tolerate += float64(tolerate)
		}
	}

	napps := make([]*AppStat, 0)
	for _, app := range apps {
		app.ErrorPercent, _ = utils.DecimalPrecision(app.errCount / float64(app.Count))
		app.AverageElapsed, _ = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex, _ = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
		napps = append(napps, app)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   napps,
	})
}
