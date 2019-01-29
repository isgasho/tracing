package service

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/g/utils"
	"github.com/mafanr/vgo/web/misc"
	"go.uber.org/zap"
)

func (web *Web) appMethods(c echo.Context) error {
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

	q := `SELECT method_id,api,service_type,elapsed,max_elapsed,min_elapsed,average_elapsed,count,err_count FROM api_details_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

	var apiID, serType, elapsed, maxE, minE, count, errCount int
	var averageE float64
	var totalElapsed int
	var api string
	ad := make(map[int]*ApiMethod)
	for iter.Scan(&apiID, &api, &serType, &elapsed, &maxE, &minE, &averageE, &count, &errCount) {
		am, ok := ad[apiID]
		if !ok {
			ad[apiID] = &ApiMethod{apiID, api, serType, 0, elapsed, maxE, minE, count, averageE, errCount, 0, ""}
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
		q := web.Cql.Query(`SELECT method_info,line FROM app_methods WHERE app_name = ? and method_id=?`, appName, am.ID)
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

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ads,
	})
}
