package app

import (
	"log"
	"math"
	"net/http"
	"time"

	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"github.com/imdevlab/tracing/pkg/constant"
	"github.com/imdevlab/tracing/web/internal/misc"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type Exception struct {
	ID             int     `json:"id"`
	Exception      string  `json:"exception"`
	ServiceType    string  `json:"service_type"`
	Elapsed        int     `json:"elapsed"`
	MaxElapsed     int     `json:"max_elapsed"`
	MinElapsed     int     `json:"min_elapsed"`
	Count          int     `json:"count"`
	AverageElapsed float64 `json:"average_elapsed"`

	Method string `json:"method"`
	Class  string `json:"class"`
}

func ExceptionStats(c echo.Context) error {
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

	q := misc.Cql.Query(`SELECT method_id,class_id,service_type,total_elapsed,max_elapsed,min_elapsed,count  FROM exception_stats WHERE app_name = ? and input_date > ? and input_date < ? `, appName, start.Unix(), end.Unix())
	iter := q.Iter()

	var methodID, exceptionID, serType, elapsed, maxE, minE, count int
	ad := make(map[int]*Exception)
	for iter.Scan(&methodID, &exceptionID, &serType, &elapsed, &maxE, &minE, &count) {
		am, ok := ad[methodID]
		if !ok {
			ave, _ := utils.DecimalPrecision(float64(elapsed / count))
			ad[methodID] = &Exception{exceptionID, "", constant.ServiceType[serType], elapsed, maxE, minE, count, ave, "", ""}
		} else {
			am.Elapsed += elapsed
			// 取最大值
			if maxE > am.MaxElapsed {
				am.MaxElapsed = maxE
			}
			// 取最小值
			if minE < am.MinElapsed {
				am.MinElapsed = minE
			}

			am.Count += count
			// 平均 = 过去的平均 * 过去总次数  + 最新的平均 * 最新的次数/ (过去总次数 + 最新次数)
			am.AverageElapsed, _ = utils.DecimalPrecision((am.AverageElapsed*float64(am.Count) + float64(elapsed)) / float64((am.Count + count)))
		}
	}

	ads := make([]*Exception, 0, len(ad))
	for _, am := range ad {
		// 获取method信息
		methodInfo := misc.GetMethodByID(appName, am.ID)
		class, method := misc.SplitMethod(methodInfo)

		am.Method = method
		am.Class = class

		// 获取exception信息
		am.Exception = misc.GetClassByID(appName, am.ID)
		ads = append(ads, am)
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
		return err
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ads,
	})
}

func ExceptionDashboard(c echo.Context) error {
	appName := c.FormValue("app_name")
	eid := c.FormValue("exception_id")
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
		cs := time2String(current)
		timeline = append(timeline, cs)
		timeBucks[cs] = &Stat{}
		current = current.Add(time.Duration(step) * time.Minute)
	}

	// 读取相应数据，按照时间填到对应的桶中
	q := misc.Cql.Query(`SELECT total_elapsed,count,input_date FROM exception_stats WHERE app_name = ? and class_id = ?  and input_date > ? and input_date < ?  ALLOW FILTERING `, appName, eid, start.Unix(), end.Unix())
	iter := q.Iter()

	// apps := make(map[string]*AppStat)
	var count int
	var tElapsed int
	var inputDate int64
	for iter.Scan(&tElapsed, &count, &inputDate) {
		t := time.Unix(inputDate, 0)
		// 计算该时间落在哪个时间桶里
		i := int(t.Sub(start).Minutes()) / step
		t1 := start.Add(time.Minute * time.Duration(i*step))

		ts := time2String(t1)
		app := timeBucks[ts]
		app.Count += count
		app.totalElapsed += float64(tElapsed)
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, misc.Cql.Closed())
	}

	// 对每个桶里的数据进行计算
	for _, app := range timeBucks {
		ep, _ := utils.DecimalPrecision(app.errCount / float64(app.Count))
		app.ErrorPercent = 100 * ep
		app.AverageElapsed, _ = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Count = app.Count / step
	}

	// 把结果数据按照时间点顺序存放
	//请求次数列表
	countList := make([]int, 0)
	//耗时列表
	elapsedList := make([]float64, 0)

	for _, ts := range timeline {
		app := timeBucks[ts]
		if math.IsNaN(app.AverageElapsed) {
			app.AverageElapsed = 0
		}

		countList = append(countList, app.Count)
		elapsedList = append(elapsedList, app.AverageElapsed)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data: DashResult{
			Timeline:    timeline,
			CountList:   countList,
			ElapsedList: elapsedList,
		},
	})
}
