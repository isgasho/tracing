package app

/* 应用Dashboard */
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

type DashResult struct {
	Timeline    []string  `json:"timeline"`
	CountList   []int     `json:"count_list"`
	ElapsedList []float64 `json:"elapsed_list"`
	ApdexList   []float64 `json:"apdex_list"`
	ErrorList   []float64 `json:"error_list"`
}

func Dashboard(c echo.Context) error {
	appName := c.FormValue("app_name")
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
	q := misc.Cql.Query(`SELECT total_elapsed,count,err_count,satisfaction,tolerate,input_date FROM api_stats WHERE app_name = ? and input_date > ? and input_date < ? `, appName, start.Unix(), end.Unix())
	iter := q.Iter()

	// apps := make(map[string]*AppStat)
	var count int
	var tElapsed, errCount, satisfaction, tolerate int
	var inputDate int64
	for iter.Scan(&tElapsed, &count, &errCount, &satisfaction, &tolerate, &inputDate) {
		t := time.Unix(inputDate, 0)
		// 计算该时间落在哪个时间桶里
		i := int(t.Sub(start).Minutes()) / step
		t1 := start.Add(time.Minute * time.Duration(i*step))

		ts := misc.TimeToChartString(t1)
		app := timeBucks[ts]
		app.Count += count
		app.totalElapsed += float64(tElapsed)
		app.errCount += float64(errCount)
		app.satisfaction += float64(satisfaction)
		app.tolerate += float64(tolerate)
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, misc.Cql.Closed())
	}

	// 对每个桶里的数据进行计算
	for _, app := range timeBucks {
		app.ErrorPercent = 100 * utils.DecimalPrecision(app.errCount/float64(app.Count))
		app.AverageElapsed = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
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
