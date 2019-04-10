package service

// Points stats point
type Points struct {
	points map[int64]*Stats
}

func newPoints() *Points {
	return &Points{}
}
