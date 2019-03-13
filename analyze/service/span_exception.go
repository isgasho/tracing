package service

import (
	"github.com/imdevlab/g"
	"github.com/imdevlab/tracing/analyze/misc"
	"github.com/imdevlab/tracing/proto/pinpoint/thrift/trace"
	"go.uber.org/zap"
)

// SpanExceptions ...
type SpanExceptions struct {
	apiExceptions map[int32]*Exceptions
}

// NewSpanExceptions ...
func NewSpanExceptions() *SpanExceptions {
	return &SpanExceptions{
		apiExceptions: make(map[int32]*Exceptions),
	}
}

// ExceptionRecord ...
func (spanExceptions *SpanExceptions) exceptionRecord(app *App, inputDate int64) error {

	for apiID, apiEx := range spanExceptions.apiExceptions {
		for exStr, exinfo := range apiEx.exceptions {
			query := gAnalyze.cql.Session.Query(misc.InserExceptionRecord,
				app.AppName,
				apiID,
				exStr,
				inputDate,
				exinfo.elapsed,
				exinfo.maxElapsed,
				exinfo.minElapsed,
				exinfo.count,
			)
			if err := query.Exec(); err != nil {
				g.L.Warn("exceptionRecord error", zap.String("error", err.Error()), zap.String("SQL", query.String()))
			}

		}
	}
	return nil
}

// exceptionCounter ...
func (spanExceptions *SpanExceptions) exceptionCounter(events []*trace.TSpanEvent, chunkEvents []*trace.TSpanEvent) error {

	for _, event := range events {
		exinfo := event.GetExceptionInfo()
		if exinfo == nil {
			continue
		}

		apiEx, ok := spanExceptions.apiExceptions[event.GetApiId()]
		if !ok {
			apiEx = NewExceptions()
			spanExceptions.apiExceptions[event.GetApiId()] = apiEx
		}
		ex, ok := apiEx.exceptions[exinfo.GetStringValue()]
		if !ok {
			ex = NewException()
			apiEx.exceptions[exinfo.GetStringValue()] = ex
		}

		ex.elapsed += int(event.GetEndElapsed())
		ex.serviceType = event.GetServiceType()
		ex.count++

		if int(event.GetEndElapsed()) > ex.maxElapsed {
			ex.maxElapsed = int(event.GetEndElapsed())
		}

		if ex.minElapsed == 0 || ex.minElapsed > int(event.GetEndElapsed()) {
			ex.minElapsed = int(event.GetEndElapsed())
		}
	}

	for _, event := range chunkEvents {
		exinfo := event.GetExceptionInfo()
		if exinfo == nil {
			continue
		}

		apiEx, ok := spanExceptions.apiExceptions[event.GetApiId()]
		if !ok {
			apiEx = NewExceptions()
			spanExceptions.apiExceptions[event.GetApiId()] = apiEx
		}
		ex, ok := apiEx.exceptions[exinfo.GetStringValue()]
		if !ok {
			ex = NewException()
			apiEx.exceptions[exinfo.GetStringValue()] = ex
		}

		ex.elapsed += int(event.GetEndElapsed())
		ex.serviceType = event.GetServiceType()
		ex.count++

		if int(event.GetEndElapsed()) > ex.maxElapsed {
			ex.maxElapsed = int(event.GetEndElapsed())
		}

		if ex.minElapsed == 0 || ex.minElapsed > int(event.GetEndElapsed()) {
			ex.minElapsed = int(event.GetEndElapsed())
		}
	}
	return nil
}

// Exceptions ...
type Exceptions struct {
	exceptions map[string]*Exception
}

// NewExceptions ...
func NewExceptions() *Exceptions {
	return &Exceptions{
		exceptions: make(map[string]*Exception),
	}
}

// Exception ...
type Exception struct {
	serviceType int16
	count       int
	elapsed     int
	maxElapsed  int
	minElapsed  int
}

// NewException ...
func NewException() *Exception {
	return &Exception{}
}
