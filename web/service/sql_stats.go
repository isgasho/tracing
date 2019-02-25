package service

import (
	"net/http"

	"github.com/imdevlab/g"
	"github.com/imdevlab/g/utils"
	"github.com/imdevlab/vgo/web/misc"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type SqlStat struct {
	ID             int     `json:"-"`
	SQL            string  `json:"sql"`
	MaxElapsed     int     `json:"max_elapsed"`
	MinElapsed     int     `json:"min_elapsed"`
	Count          int     `json:"count"`
	AverageElapsed float64 `json:"average_elapsed"`
	ErrorCount     int     `json:"error_count"`
}

func (web *Web) sqlStats(c echo.Context) error {
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

	q := `SELECT sql,max_elapsed,min_elapsed,average_elapsed,count,err_count FROM sql_stats WHERE app_name = ? and input_date > ? and input_date < ? `
	iter := web.Cql.Query(q, appName, start.Unix(), end.Unix()).Iter()

	var sqlID, maxE, minE, count, errCount int
	var averageE float64
	ad := make(map[int]*SqlStat)
	for iter.Scan(&sqlID, &maxE, &minE, &averageE, &count, &errCount) {
		am, ok := ad[sqlID]
		if !ok {
			ad[sqlID] = &SqlStat{sqlID, "", maxE, minE, count, averageE, errCount}
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
			am.AverageElapsed, _ = utils.DecimalPrecision((am.AverageElapsed*float64(am.Count) + averageE*float64(count)) / float64((am.Count + count)))
		}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	ads := make([]*SqlStat, 0, len(ad))
	for _, am := range ad {
		// 通过apiID 获取api name
		q := web.Cql.Query(`SELECT sql_info FROM app_sqls WHERE app_name = ? and sql_id=?`, appName, am.ID)
		var sql string
		err := q.Scan(&sql)
		if err != nil {
			g.L.Info("access database error", zap.Error(err), zap.String("query", q.String()))
			continue
		}

		s, _ := g.B64.DecodeString(sql)
		am.SQL = utils.Bytes2String(s)

		ads = append(ads, am)
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ads,
	})
}
