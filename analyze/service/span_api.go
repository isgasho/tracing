package service

import (
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/analyze/misc"
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

// apiRecord ...
func (spanUrls *SpanAPIs) apiRecord(app *App, recordTime int64) error {
	for apiStr, api := range spanUrls.apis {
		query := gAnalyze.cql.Session.Query(misc.InserAPIRecord,
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
		)
		if err := query.Exec(); err != nil {
			g.L.Warn("apiRecord error", zap.String("error", err.Error()), zap.String("sql", query.String()))
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
