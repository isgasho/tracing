package service

import (
	"github.com/mafanr/g"
	"github.com/mafanr/vgo/analyze/service/blink"
	"github.com/mafanr/vgo/analyze/service/stats"
	"go.uber.org/zap"
)

// Analyze ...
type Analyze struct {
	stats *stats.Stats
	blink *blink.Blink
}

// New ...
func New() *Analyze {
	return &Analyze{
		stats: stats.New(),
		blink: blink.New(),
	}
}

// Start ...
func (anlyze *Analyze) Start() error {

	if err := anlyze.blink.Start(); err != nil {
		g.L.Fatal("Start", zap.String("error", err.Error()))
	}

	if err := anlyze.stats.Start(); err != nil {
		g.L.Fatal("Start", zap.String("error", err.Error()))
	}

	g.L.Info("Start ok!")
	return nil
}

// Close ...
func (anlyze *Analyze) Close() error {

	if anlyze.blink != nil {
		anlyze.blink.Close()
	}

	if anlyze.stats != nil {
		anlyze.stats.Close()
	}

	g.L.Info("Close ok!")
	return nil
}
