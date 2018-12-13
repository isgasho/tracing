package stats

import "github.com/mafanr/g"

// Stats 离线计算
type Stats struct {
}

// New ...
func New() *Stats {
	return &Stats{}
}

// Start ...
func (s *Stats) Start() error {
	g.L.Info("Start Stats")
	return nil
}

// Close ...
func (s *Stats) Close() error {
	g.L.Info("Close Stats")

	return nil
}
