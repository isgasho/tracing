package service

import "github.com/imdevlab/tracing/collector/stats"

// Points stats point
type Points struct {
	points map[int64]*stats.Stats
}

func newPoints() *Points {
	return &Points{}
}
