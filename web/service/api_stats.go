package service

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/web/misc"
)

type ApiStat struct {
	URL            string `json:"url"`
	MaxElapsed     int    `json:"max_elapsed"`
	MinElapsed     int    `json:"min_elapsed"`
	Count          int    `json:"count"`
	AverageElapsed int    `json:"average_elapsed"`
	ErrorCount     int    `json:"error_count"`
}

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
	q := `SELECT url,max_elapsed,min_elapsed,average_elapsed,count,err_count FROM rpc_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

	// apps := make(map[string]*AppStat)
	var maxElapsed, minElapsed, count, aElapsed, errCount int
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
			as.AverageElapsed = (as.AverageElapsed*as.Count + aElapsed*count) / (as.Count + count)
		}
	}

	if err := iter.Close(); err != nil {
		log.Println("close iter error:", err, web.Cql.Closed())
	}

	// 对每个桶里的数据进行计算
	apiStats := make([]*ApiStat, len(ass))
	for _, as := range ass {
		apiStats = append(apiStats, as)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   apiStats,
	})
}
