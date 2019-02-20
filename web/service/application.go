package service

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
	"github.com/mafanr/vgo/web/misc"
	"go.uber.org/zap"
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
	q := `SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE input_date > ? and input_date < ? `
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

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, web.Cql.Closed())
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

func (web *Web) appListWithSetting(c echo.Context) error {
	li := web.getLoginInfo(c)
	appShow, appNames := web.userAppSetting(li.ID)

	napps := make([]*AppStat, 0)

	now := time.Now()
	// 查询缓存数据是否存在和过期
	// if web.cache.appList == nil || now.Sub(web.cache.appListUpdate).Seconds() > CacheUpdateIntv {
	// 取过去6分钟的数据
	start := now.Unix() - 450

	apps := make(map[string]*AppStat)
	var q *gocql.Query
	var appName string
	var count int
	var tElapsed, errCount, satisfaction, tolerate int

	if appShow == 1 {
		q = web.Cql.Query(`SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE input_date > ? and input_date < ? `, start, now.Unix())
		iter := q.Iter()

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

		if err := iter.Close(); err != nil {
			log.Println("close iter error:", err)
		}
	} else {
		for _, an := range appNames {
			err := web.Cql.Query(`SELECT app_name,total_elapsed,count,err_count,satisfaction,tolerate FROM api_stats WHERE app_name =? and input_date > ? and input_date < ? `, an, start, now.Unix()).Scan(&appName, &tElapsed, &count, &errCount, &satisfaction, &tolerate)
			if err != nil {
				log.Println("select app stats error:", err)
			}

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

	}

	for _, app := range apps {
		app.ErrorPercent, _ = utils.DecimalPrecision(app.errCount / float64(app.Count))
		app.AverageElapsed, _ = utils.DecimalPrecision(app.totalElapsed / float64(app.Count))
		app.Apdex, _ = utils.DecimalPrecision((app.satisfaction + app.tolerate/2) / float64(app.Count))
		napps = append(napps, app)
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
	start, end, err := misc.StartEndDate(c)
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

	// 把start-end分为30个桶
	timeline := make([]string, 0)
	timeBucks := make(map[string]*AppStat)
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
		timeBucks[cs] = &AppStat{}
		current = current.Add(time.Duration(step) * time.Minute)
	}

	// 读取相应数据，按照时间填到对应的桶中
	q := `SELECT total_elapsed,count,err_count,satisfaction,tolerate,input_date FROM api_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

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

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, web.Cql.Closed())
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

func (web *Web) appNames(c echo.Context) error {
	q := `SELECT app_name FROM apps `
	iter := web.Cql.Query(q).Iter()

	appNames := make([]string, 0)
	var appName string
	for iter.Scan(&appName) {
		appNames = append(appNames, appName)
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   appNames,
	})
}

func (web *Web) appNamesWithSetting(c echo.Context) error {
	li := web.getLoginInfo(c)
	appShow, appNames := web.userAppSetting(li.ID)

	ans := make([]string, 0)
	if appShow == 1 { // 显示全部应用
		q := `SELECT app_name FROM apps `
		iter := web.Cql.Query(q).Iter()

		var appName string
		for iter.Scan(&appName) {
			ans = append(ans, appName)
		}
		if err := iter.Close(); err != nil {
			g.L.Warn("close iter error:", zap.Error(err))
		}

	} else {
		ans = appNames
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ans,
	})
}

type AgentStat struct {
	AgentID     string `json:"agent_id"`
	HostName    string `json:"host_name"`
	IsLive      bool   `json:"is_live"`
	IsContainer bool   `json:"is_container"`
}

func (web *Web) agentList(c echo.Context) error {
	appName := c.FormValue("app_name")
	q := `SELECT agent_id,host_name,is_live,is_container FROM agents WHERE app_name=?`
	iter := web.Cql.Query(q, appName).Iter()

	var agentID, hostName string
	var isLive, isContainer bool

	agents := make([]*AgentStat, 0)
	for iter.Scan(&agentID, &hostName, &isLive, &isContainer) {
		agents = append(agents, &AgentStat{
			AgentID:     agentID,
			HostName:    hostName,
			IsLive:      isLive,
			IsContainer: isContainer,
		})
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   agents,
	})
}

func time2String(t time.Time) string {
	return t.Format("01-02 15:04")
}

// 获取用户的应用设定
func (web *Web) userAppSetting(user string) (int, []string) {
	q := web.Cql.Query(`SELECT app_show,app_names FROM account WHERE id=?`, user)
	var appShow int
	var appNameS string
	err := q.Scan(&appShow, &appNameS)
	if err != nil {
		return 1, nil
	}

	appNames := make([]string, 0)
	err = json.Unmarshal([]byte(appNameS), &appNames)
	if err != nil {
		return 1, nil
	}

	return appShow, appNames
}

func (web *Web) appApis(c echo.Context) error {
	appName := c.FormValue("app_name")
	if appName == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	q := `SELECT api FROM app_apis WHERE app_name=?`
	iter := web.Cql.Query(q, appName).Iter()

	var api string
	apis := make([]string, 0)
	for iter.Scan(&api) {
		apis = append(apis, api)
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   apis,
	})
}
