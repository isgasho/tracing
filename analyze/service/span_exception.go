package service

// SpanExceptions ...
type SpanExceptions struct {
	exceptions map[int32]*SpanException
}

// NewSpanExceptions ...
func NewSpanExceptions() *SpanExceptions {
	return &SpanExceptions{
		exceptions: make(map[int32]*SpanException),
	}
}

// exceptionCounter ...
func (spanExceptions *SpanExceptions) exceptionCounter(urlStr string, elapsed int, isError int) error {
	// url, ok := spanExceptions.exceptions[urlStr]
	// if !ok {
	// 	url = NewSpanURL()
	// 	spanUrls.urls[urlStr] = url
	// }
	// url.elapsed += elapsed
	// url.count++
	// if isError != 0 {
	// 	url.errCount++
	// }

	// if elapsed > url.maxElapsed {
	// 	url.maxElapsed = url.elapsed
	// }

	// if url.minElapsed == 0 || url.minElapsed > elapsed {
	// 	url.minElapsed = elapsed
	// }

	// url.averageElapsed = url.elapsed / url.count
	return nil
}

// SpanException ...
type SpanException struct {
	serType        int
	elapsed        int
	maxElapsed     int
	minElapsed     int
	averageElapsed int
	count          int
	errCount       int
}

// NewSpanException ...
func NewSpanException() *SpanException {
	return &SpanException{}
}
