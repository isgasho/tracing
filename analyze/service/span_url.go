package service

import (
	"github.com/mafanr/g"
	"go.uber.org/zap"
)

// SpanURLs ...
type SpanURLs struct {
	urls map[string]*SpanURL
}

// NewSpanURLs ...
func NewSpanURLs() *SpanURLs {
	return &SpanURLs{
		urls: make(map[string]*SpanURL),
	}
}

// urlCounter ...
func (spanUrls *SpanURLs) urlCounter(urlStr string, elapsed int, isError int) error {
	url, ok := spanUrls.urls[urlStr]
	if !ok {
		url = NewSpanURL()
		spanUrls.urls[urlStr] = url
	}
	url.elapsed += elapsed
	url.count++
	if isError != 0 {
		url.errCount++
	}

	if elapsed > url.maxElapsed {
		url.maxElapsed = url.elapsed
	}

	if url.minElapsed == 0 || url.minElapsed > elapsed {
		url.minElapsed = elapsed
	}

	url.averageElapsed = url.elapsed / url.count
	return nil
}

var gInserRPCRecord string = `INSERT INTO rpc_record (app_name, start_time, url, elapsed, max_elapsed, min_elapsed, average_elapsed, count, err_count)
 VALUES (?,?,?,?,?,?,?,?,?);`

// urlRecord ...
func (spanUrls *SpanURLs) urlRecord(app *App, recordTime int64) error {
	for urlStr, url := range spanUrls.urls {
		if err := gAnalyze.db.Session.Query(gInserRPCRecord,
			app.AppName,
			recordTime,
			urlStr,
			url.elapsed,
			url.maxElapsed,
			url.minElapsed,
			url.averageElapsed,
			url.count,
			url.errCount,
		).Exec(); err != nil {
			g.L.Warn("urlRecord error", zap.String("error", err.Error()))
		}
	}
	return nil
}

// SpanURL ...
type SpanURL struct {
	averageElapsed int
	elapsed        int
	count          int
	errCount       int
	minElapsed     int
	maxElapsed     int
}

// NewSpanURL ...
func NewSpanURL() *SpanURL {
	return &SpanURL{}
}
