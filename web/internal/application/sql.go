package app

/* SQL统计 */

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type SqlStat struct {
	ID             int     `json:"id"`
	SQL            string  `json:"sql"`
	MaxElapsed     int     `json:"max_elapsed"`
	MinElapsed     int     `json:"min_elapsed"`
	Count          int     `json:"count"`
	AverageElapsed float64 `json:"average_elapsed"`
	ErrorCount     int     `json:"error_count"`
}

func SqlStats(c echo.Context) error {
	start, end, err := misc.StartEndDate(c)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusOK,
			ErrCode: g.ParamInvalidC,
			Message: "日期参数不合法",
		})
	}

	appName := c.FormValue("app_name")
	if appName == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	q := `SELECT sql,max_elapsed,min_elapsed,elapsed,count,err_count FROM sql_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := misc.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

	var sqlID, maxE, minE, count, errCount, elapsed int
	ad := make(map[int]*SqlStat)
	for iter.Scan(&sqlID, &maxE, &minE, &elapsed, &count, &errCount) {
		am, ok := ad[sqlID]
		if !ok {
			ad[sqlID] = &SqlStat{sqlID, "", maxE, minE, count, utils.DecimalPrecision(float64(elapsed / count)), errCount}
		} else {
			// 取最大值
			if maxE > am.MaxElapsed {
				am.MaxElapsed = maxE
			}
			// 取最小值
			if minE < am.MinElapsed {
				am.MinElapsed = minE
			}

			am.Count += count
			am.ErrorCount += errCount
			// 平均 = 过去的平均 * 过去总次数  + 最新的平均 * 最新的次数/ (过去总次数 + 最新次数)
			am.AverageElapsed = utils.DecimalPrecision((am.AverageElapsed*float64(am.Count) + float64(elapsed)) / float64((am.Count + count)))
		}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	ads := make([]*SqlStat, 0, len(ad))
	for _, am := range ad {
		am.SQL = misc.GetSqlByID(appName, am.ID)

		ads = append(ads, am)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ads,
	})
}

func SqlDashboard(c echo.Context) error {
	appName := c.FormValue("app_name")
	sqlID := c.FormValue("sql_id")
	start, end, err := misc.StartEndDate(c)
	if err != nil {
		g.L.Info("日期参数不合法", zap.String("start", c.FormValue("start")), zap.String("end", c.FormValue("end")), zap.Error(err))
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

	// 把start-end分为30个桶
	timeline := make([]string, 0)
	timeBucks := make(map[string]*Stat)
	current := start
	var step int
	if end.Sub(start).Minutes() <= 60 {
		step = 1
	} else {
		step = int(end.Sub(start).Minutes())/30 + 1
	}

	for {
		if current.Unix() > end.Unix() {
			break
		}
		cs := misc.TimeToChartString(current)
		timeline = append(timeline, cs)
		timeBucks[cs] = &Stat{}
		current = current.Add(time.Duration(step) * time.Minute)
	}

	// 读取相应数据，按照时间填到对应的桶中
	q := `SELECT elapsed,count,err_count,input_date FROM sql_stats WHERE app_name = ?  and sql = ? and input_date > ? and input_date < ? `
	iter := misc.Cql.Query(q, appName, sqlID, start.Unix(), end.Unix()).Iter()

	// apps := make(map[string]*AppStat)
	var count int
	var tElapsed, errCount int
	var inputDate int64
	for iter.Scan(&tElapsed, &count, &errCount, &inputDate) {
		t := time.Unix(inputDate, 0)
		// 计算该时间落在哪个时间桶里
		i := int(t.Sub(start).Minutes()) / step
		t1 := start.Add(time.Minute * time.Duration(i*step))

		ts := misc.TimeToChartString(t1)
		app := timeBucks[ts]
		app.Count += count
		app.totalElapsed += float64(tElapsed)
		app.errCount += float64(errCount)
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, misc.Cql.Closed())
	}

	// 对每个桶里的数据进行计算
	for _, app := range timeBucks {
		app.ErrorPercent = 100 * utils.DecimalPrecision(app.errCount/float64(app.Count))
		app.AverageElapsed = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Count = app.Count / step
	}

	// 把结果数据按照时间点顺序存放
	//请求次数列表
	countList := make([]int, 0)
	//耗时列表
	elapsedList := make([]float64, 0)
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
		errorList = append(errorList, app.ErrorPercent)

	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data: DashResult{
			Timeline:    timeline,
			CountList:   countList,
			ElapsedList: elapsedList,
			ErrorList:   errorList,
		},
	})
}
