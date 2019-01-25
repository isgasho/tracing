package service

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
	"github.com/mafanr/vgo/web/misc"
	"go.uber.org/zap"
)

type ApiStat struct {
	API            string  `json:"api"`
	MaxElapsed     int     `json:"max_elapsed"`
	MinElapsed     int     `json:"min_elapsed"`
	Count          int     `json:"count"`
	AverageElapsed float64 `json:"average_elapsed"`
	ErrorCount     int     `json:"error_count"`
}

// type ApiStats []*ApiStat

// func (a ApiStats) Len() int { // 重写 Len() 方法
// 	return len(a)
// }
// func (a ApiStats) Swap(i, j int) { // 重写 Swap() 方法
// 	a[i], a[j] = a[j], a[i]
// }
// func (a ApiStats) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
// 	return a[j].AverageElapsed < a[i].AverageElapsed
// }

// 单个应用下，所有接口的统计信息
func (web *Web) apiStats(c echo.Context) error {
	appName := c.FormValue("app_name")
	start, end, err := misc.StartEndDate(c)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusOK,
			ErrCode: g.ParamInvalidC,
			Message: "日期参数不合法",
		})
	}

	// 读取相应数据，按照时间填到对应的桶中
	q := `SELECT api,max_elapsed,min_elapsed,average_elapsed,count,err_count FROM api_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

	// apps := make(map[string]*AppStat)
	var maxElapsed, minElapsed, count, errCount int
	var aElapsed float64
	var url string
	ass := make(map[string]*ApiStat)
	for iter.Scan(&url, &maxElapsed, &minElapsed, &aElapsed, &count, &errCount) {
		as, ok := ass[url]
		if !ok {
			ass[url] = &ApiStat{url, maxElapsed, minElapsed, count, aElapsed, errCount}
		} else {
			// 取最大值
			if maxElapsed > as.MaxElapsed {
				as.MaxElapsed = maxElapsed
			}
			// 取最小值
			if minElapsed < as.MinElapsed {
				as.MinElapsed = minElapsed
			}

			as.Count += count
			as.ErrorCount += errCount
			// 平均 = 过去的平均 * 过去总次数  + 最新的平均 * 最新的次数/ (过去总次数 + 最新次数)
			as.AverageElapsed, _ = utils.DecimalPrecision((as.AverageElapsed*float64(as.Count) + aElapsed*float64(count)) / float64((as.Count + count)))
		}
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, web.Cql.Closed())
	}

	// 对每个桶里的数据进行计算
	apiStats := make([]*ApiStat, 0, len(ass))
	for _, as := range ass {
		apiStats = append(apiStats, as)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   apiStats,
	})
}

type ApiMethod struct {
	ID             int     `json:"-"`
	API            string  `json:"api"`
	ServiceType    int     `json:"service_type"`
	RatioElapsed   int     `json:"ratio_elapsed"`
	Elapsed        int     `json:"elapsed"`
	MaxElapsed     int     `json:"max_elapsed"`
	MinElapsed     int     `json:"min_elapsed"`
	Count          int     `json:"count"`
	AverageElapsed float64 `json:"average_elapsed"`
	ErrorCount     int     `json:"error_count"`

	Line   int    `json:"line"`
	Method string `json:"method"`
}

func (web *Web) apiDetail(c echo.Context) error {
	start, end, err := misc.StartEndDate(c)
	if err != nil {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusOK,
			ErrCode: g.ParamInvalidC,
			Message: "日期参数不合法",
		})
	}

	appName := c.FormValue("app_name")
	url := c.FormValue("api")

	if appName == "" || url == "" {
		return c.JSON(http.StatusOK, g.Result{
			Status:  http.StatusBadRequest,
			ErrCode: g.ParamInvalidC,
			Message: g.ParamInvalidE,
		})
	}

	q := `SELECT api_id,ser_type,elapsed,max_elapsed,min_elapsed,average_elapsed,count,err_count FROM api_details_stats WHERE app_name = ? and api = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, url, start.Unix(), end.Unix()).Iter()

	var apiID, serType, elapsed, maxE, minE, count, errCount int
	var averageE float64
	var totalElapsed int
	ad := make(map[int]*ApiMethod)
	for iter.Scan(&apiID, &serType, &elapsed, &maxE, &minE, &averageE, &count, &errCount) {
		am, ok := ad[apiID]
		if !ok {
			ad[apiID] = &ApiMethod{apiID, url, serType, 0, elapsed, maxE, minE, count, averageE, errCount, 0, ""}
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
			am.ErrorCount += errCount
			// 平均 = 过去的平均 * 过去总次数  + 最新的平均 * 最新的次数/ (过去总次数 + 最新次数)
			am.AverageElapsed, _ = utils.DecimalPrecision((am.AverageElapsed*float64(am.Count) + averageE*float64(count)) / float64((am.Count + count)))
		}

		totalElapsed += elapsed
	}

	ads := make([]*ApiMethod, 0, len(ad))
	for _, am := range ad {
		// 计算耗时占比
		am.RatioElapsed = am.Elapsed * 100 / totalElapsed
		// 通过apiID 获取api name
		q := web.Cql.Query(`SELECT api_info,line FROM app_apis WHERE app_name = ? and api_id=?`, appName, am.ID)
		var apiInfo string
		var line int
		err := q.Scan(&apiInfo, &line)
		if err != nil {
			g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
			continue
		}

		am.Line = line
		am.Method = apiInfo

		ads = append(ads, am)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ads,
	})
}
