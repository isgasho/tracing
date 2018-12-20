package service

import (
	"github.com/mafanr/g"
)

// Stats 离线计算
type Stats struct {
}

// NewStats ...
func NewStats() *Stats {
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

func (s *Stats) counter(route chan *App) {

}
