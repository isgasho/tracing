package web

import (
	"fmt"
	"log"
	"math"
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
	napps := make([]*AppStat, 0)

	now := time.Now()
	// 查询缓存数据是否存在和过期
	// if web.cache.appList == nil || now.Sub(web.cache.appListUpdate).Seconds() > CacheUpdateIntv {
	// 取过去6分钟的数据
	start := now.Unix() - 450
	q := `SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM rpc_stats WHERE input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, start, now.Unix()).Iter()

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

	for _, app := range apps {
		app.ErrorPercent, _ = utils.DecimalPrecision(app.errCount / float64(app.Count))
		app.AverageElapsed, _ = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex, _ = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
		napps = append(napps, app)
	}
	fmt.Println("query from database:", napps)

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err)
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

type DashResult struct {
	Timeline    []string  `json:"timeline"`
	CountList   []int     `json:"count_list"`
	ElapsedList []float64 `json:"elapsed_list"`
	ApdexList   []float64 `json:"apdex_list"`
	ErrorList   []float64 `json:"error_list"`
}

func (web *Web) appDash(c echo.Context) error {
	appName := c.FormValue("app_name")
	start, end, err := startEndDate(c)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusOK,
			ErrCode: g.ParamInvalidC,
			Message: "日期参数不合法",
		})
	}

	// 查询时间间隔要转换为30的倍数，然后按照倍数聚合相应的点，最终形成30个图绘制点
	//计算间隔
	intv := int(end.Sub(start).Minutes())

	if intv%30 != 0 {
		start = end.Add(-(time.Duration(intv/30+1)*30*time.Minute - time.Minute))
	} else {
		start = start.Add(time.Minute)
	}

	// 把start-end分为30个时间点
	timeline := make([]string, 0)
	timeBucks := make(map[string]*AppStat)
	current := start
	step := int(end.Sub(start).Minutes())/30 + 1
	for {
		if current.Unix() > end.Unix() {
			break
		}
		cs := time2String(current)
		timeline = append(timeline, cs)
		timeBucks[cs] = &AppStat{}
		current = current.Add(time.Duration(step) * time.Minute)
	}

	// 读取相应数据，按照时间填到对应的桶中
	q := `SELECT total_elapsed,count,err_count,satisfaction,tolerate,input_date FROM rpc_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()
	defer iter.Close()

	// apps := make(map[string]*AppStat)
	var count int
	var tElapsed, errCount, satisfaction, tolerate int
	var inputDate int64
	for iter.Scan(&tElapsed, &count, &errCount, &satisfaction, &tolerate, &inputDate) {
		t := time.Unix(inputDate, 0)
		// 计算该时间落在哪个时间桶里
		i := int(t.Sub(start).Minutes()) / step
		t1 := start.Add(time.Minute * time.Duration(i*step))

		ts := time2String(t1)
		app := timeBucks[ts]
		app.Count += count
		app.totalElapsed += float64(tElapsed)
		app.errCount += float64(errCount)
		app.satisfaction += float64(satisfaction)
		app.tolerate += float64(tolerate)
	}

	// 对每个桶里的数据进行计算
	for _, app := range timeBucks {
		app.ErrorPercent, _ = utils.DecimalPrecision(app.errCount / float64(app.Count))
		app.AverageElapsed, _ = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex, _ = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
		app.Count = app.Count / step
	}

	// 把结果数据按照时间点顺序存放
	//请求次数列表
	countList := make([]int, 0)
	//耗时列表
	elapsedList := make([]float64, 0)
	//apdex列表
	apdexList := make([]float64, 0)
	//错误率列表
	errorList := make([]float64, 0)

	for _, ts := range timeline {
		app := timeBucks[ts]
		if math.IsNaN(app.AverageElapsed) {
			app.AverageElapsed = 0
		}
		if math.IsNaN(app.Apdex) {
			app.Apdex = 1
		}
		if math.IsNaN(app.ErrorPercent) {
			app.ErrorPercent = 0
		}
		countList = append(countList, app.Count)
		elapsedList = append(elapsedList, app.AverageElapsed)
		apdexList = append(apdexList, app.Apdex)
		errorList = append(errorList, app.ErrorPercent)

	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data: DashResult{
			Timeline:    timeline,
			CountList:   countList,
			ElapsedList: elapsedList,
			ApdexList:   apdexList,
			ErrorList:   errorList,
		},
	})
}

func time2String(t time.Time) string {
	return t.Format("01-02 15:04")
}
