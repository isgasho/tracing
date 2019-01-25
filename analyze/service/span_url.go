package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/misc"
	"go.uber.org/zap"
)

// SpanAPIs ...
type SpanAPIs struct {
	apis map[string]*SpanAPI
}

// NewSpanAPIs ...
func NewSpanAPIs() *SpanAPIs {
	return &SpanAPIs{
		apis: make(map[string]*SpanAPI),
	}
}

// apiCounter ...
func (spanUrls *SpanAPIs) apiCounter(apiStr string, elapsed int, isError int) error {
	api, ok := spanUrls.apis[apiStr]
	if !ok {
		api = NewSpanAPI()
		spanUrls.apis[apiStr] = api
	}
	api.elapsed += elapsed
	api.count++
	if isError != 0 {
		api.errCount++
	}

	if elapsed > api.maxElapsed {
		api.maxElapsed = api.elapsed
	}

	if api.minElapsed == 0 || api.minElapsed > elapsed {
		api.minElapsed = elapsed
	}

	api.averageElapsed = float64(api.elapsed) / float64(api.count)

	if elapsed < misc.Conf.Stats.SatisfactionTime {
		api.satisfactionCount++
	} else if elapsed > misc.Conf.Stats.TolerateTime {
		api.tolerateCount++
	}

	return nil
}

var gInserRPCRecord string = `INSERT INTO api_stats (app_name, input_date, api, total_elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count, satisfaction, tolerate)
 VALUES (?,?,?,?,?,?,?,?,?,?,?);`

// apiRecord ...
func (spanUrls *SpanAPIs) apiRecord(app *App, recordTime int64) error {
	for apiStr, api := range spanUrls.apis {
		if err := gAnalyze.cql.Session.Query(gInserRPCRecord,
			app.AppName,
			recordTime,
			apiStr,
			api.elapsed,
			api.maxElapsed,
			api.minElapsed,
			api.averageElapsed,
			api.count,
			api.errCount,
			api.satisfactionCount,
			api.tolerateCount,
		).Exec(); err != nil {
			g.L.Warn("apiRecord error", zap.String("error", err.Error()))
		}
	}
	return nil
}

// SpanAPI ...
type SpanAPI struct {
	averageElapsed    float64
	elapsed           int
	count             int
	errCount          int
	minElapsed        int
	maxElapsed        int
	satisfactionCount int
	tolerateCount     int
}

// NewSpanAPI ...
func NewSpanAPI() *SpanAPI {
	return &SpanAPI{}
}
