package service

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/labstack/echo"
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/web/misc"
	"go.uber.org/zap"
)

type Trace struct {
	ID        string `json:"id"`
	API       string `json:"api"`
	Elapsed   int    `json:"y"`
	AgentID   string `json:"agent_id"`
	InputDate int64  `json:"x"`
}

type Traces []*Trace

func (a Traces) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a Traces) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a Traces) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Elapsed < a[i].Elapsed
}

type ChartTraces struct {
	Suc   bool    `json:"is_suc"`
	Xaxis []int64 `json:"timeXticks"`
	Title string  `json:"subTitle"`

	Series []*TraceSeries `json:"series"`
}

type TraceSeries struct {
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Data  []*Trace `json:"data"`
}

// traceSeries":[{"name":"success","color":"rgb(18, 147, 154,.5)","data":[{"x":1545200556716,"y":7,"traceId":"yunbaoParkApp3^1545036617750^4217","agentId":"agencyBookKeep3","startTime":"1545200556716","url":"/agencyBookKeep/financialstatementscs/getOneAccountingDataByParams","traceIp":"127.0.0.1"},
func (web *Web) queryTraces(c echo.Context) error {
	appName := c.FormValue("app_name")
	api := c.FormValue("api")
	min, _ := strconv.Atoi(c.FormValue("min_elapsed"))
	max, _ := strconv.Atoi(c.FormValue("max_elapsed"))
	limit, err := strconv.Atoi(c.FormValue("limit"))
	if err != nil {
		limit = 50
	}

	start, end, _ := misc.StartEndDate(c)

	var q *gocql.Query
	if api == "" {
		if max == 0 {
			q = web.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date FROM app_operation_index WHERE app_name=? and input_date > ? and input_date < ? and elapsed >= ? ALLOW FILTERING`, appName, start.Unix()*1000, end.Unix()*1000, min)
		} else {
			q = web.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date FROM app_operation_index WHERE app_name=? and input_date > ? and input_date < ? and elapsed >= ? and elapsed <= ? ALLOW FILTERING`, appName, start.Unix()*1000, end.Unix()*1000, min, max)
		}
	} else {
		if max == 0 {
			q = web.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date FROM app_operation_index WHERE app_name=? and api=?  and input_date > ? and input_date < ? and elapsed >= ? ALLOW FILTERING`, appName, api, start.Unix()*1000, end.Unix()*1000, min)
		} else {
			q = web.Cql.Query(`SELECT trace_id,api,elapsed,agent_id,input_date FROM app_operation_index WHERE app_name=? and api=?  and input_date > ? and input_date < ? and elapsed >= ? and elapsed <= ? ALLOW FILTERING`, appName, api, start.Unix()*1000, end.Unix()*1000, min, max)
		}
	}

	iter := q.Iter()

	var elapsed int
	var inputDate int64
	var tid, agentID string

	traceMap := make(map[string]*Trace)
	for iter.Scan(&tid, &api, &elapsed, &agentID, &inputDate) {
		traceMap[tid] = &Trace{tid, api, elapsed, agentID, inputDate}
	}

	if err := iter.Close(); err != nil {
		g.L.Warn("close iter error:", zap.Error(err))
	}

	traces := make(Traces, 0, len(traceMap))
	for _, t := range traceMap {
		traces = append(traces, t)
	}

	sort.Sort(traces)

	// 取出耗时最高的limit数量的trace
	if limit < len(traces) {
		traces = traces[:limit]
	}

	ct := &ChartTraces{}
	if len(traces) == 0 {
		ct.Suc = false
	} else {
		ct.Suc = true
		ct.Xaxis = []int64{traces[0].InputDate / 1000, traces[len(traces)-1].InputDate / 1000}
		ct.Title = fmt.Sprintf("success: %d, error: %d", len(traces), 0)

		sucData := &TraceSeries{
			Name:  "success",
			Color: "rgb(18, 147, 154,.5)",
			Data:  traces,
		}

		errData := &TraceSeries{
			Name:  "error",
			Color: "rgba(223, 83, 83, .5)",
			Data:  make([]*Trace, 0),
		}

		ct.Series = []*TraceSeries{sucData, errData}
	}

	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   ct,
	})
}

func queryTrace(c echo.Context) error {
	return c.JSON(http.StatusOK, g.Result{
		Status: http.StatusOK,
		Data:   traceTest,
	})
}
