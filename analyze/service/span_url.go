package service

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
