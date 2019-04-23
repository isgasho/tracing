package app

/* 应用首页列表 */
import (
	"log"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/imdevlab/tracing/web/internal/session"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

func List(c echo.Context) error {
	napps := make([]*Stat, 0)

	now := time.Now()
	// 查询缓存数据是否存在和过期
	// if web.cache.appList == nil || now.Sub(web.cache.appListUpdate).Seconds() > CacheUpdateIntv {
	// 取过去30分钟的数据
	start := now.Unix() - 30*60
	q := `SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE input_date > ? and input_date < ? `
	iter := misc.Cql.Query(q, start, now.Unix()).Iter()

	apps := make(map[string]*Stat)
	var appName string
	var count int
	var tElapsed, errCount, satisfaction, tolerate int

	for iter.Scan(&appName, &tElapsed, &count, &errCount, &satisfaction, &tolerate) {
		app, ok := apps[appName]
		if !ok {
			apps[appName] = &Stat{
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

	for _, app := range apps {
		app.ErrorPercent = 100 * utils.DecimalPrecision(app.errCount/float64(app.Count))
		app.AverageElapsed = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
		napps = append(napps, app)
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, misc.Cql.Closed())
	}
	// 	// 更新缓存
	// 	if len(napps) > 0 {
	// 		web.cache.appList = napps
	// 		web.cache.appListUpdate = now
	// 	}
	// } else {
	// 	napps = web.cache.appList
	// 	fmt.Println("query from cache:", napps)
	// }

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   napps,
	})
}

func ListWithSetting(c echo.Context) error {
	li := session.GetLoginInfo(c)
	appShow, appNames := UserSetting(li.ID)

	stats := make([]*Stat, 0)

	now := time.Now()
	// 取过去6分钟的数据
	start := now.Unix() - 30*60

	statsMap := make(map[string]*Stat)
	var q *gocql.Query
	var appName string
	var count int
	var tElapsed, errCount, satisfaction, tolerate int

	if appShow == 1 {
		q = misc.Cql.Query(`SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE input_date > ? and input_date < ? `, start, now.Unix())
		iter := q.Iter()

		for iter.Scan(&appName, &tElapsed, &count, &errCount, &satisfaction, &tolerate) {
			app, ok := statsMap[appName]
			if !ok {
				statsMap[appName] = &Stat{
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

		if err := iter.Close(); err != nil {
			log.Println("close iter error:", err)
		}
	} else {
		for _, an := range appNames {
			err := misc.Cql.Query(`SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE app_name =? and input_date > ? and input_date < ? `, an, start, now.Unix()).Scan(&appName, &tElapsed, &count, &errCount, &satisfaction, &tolerate)
			if err != nil {
				log.Println("select app stats error:", err)
			}

			app, ok := statsMap[appName]
			if !ok {
				statsMap[appName] = &Stat{
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

	}

	for _, stat := range statsMap {
		stat.ErrorPercent = utils.DecimalPrecision(stat.errCount/float64(stat.Count)) * 100
		stat.AverageElapsed = utils.DecimalPrecision(stat.totalElapsed / float64(stat.Count))
		stat.Apdex = utils.DecimalPrecision((stat.satisfaction + stat.tolerate/2) / float64(stat.Count))
		stats = append(stats, stat)
	}

	// 获取所有应用，不在之前列表的应用，所有数据置为0
	allApps := appList()
	for _, name := range allApps {
		_, ok := statsMap[name]
		if !ok {
			stat := &Stat{
				Name:  name,
				Apdex: 1,
			}
			stats = append(stats, stat)
		}
	}
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   stats,
	})
}

func appList() []string {
	q := `SELECT app_name FROM apps `
	iter := misc.Cql.Query(q).Iter()

	appNames := make([]string, 0)
	var appName string
	for iter.Scan(&appName) {
		appNames = append(appNames, appName)
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return appNames
}
