package blink

import (
	"github.com/mafanr/g"
)

// Blink 实时计算
type Blink struct {
}

// New ...
func New() *Blink {
	return &Blink{}
}

// Start ...
func (s *Blink) Start() error {
	g.L.Info("Start Blink")
	return nil
}

// Close ...
func (s *Blink) Close() error {
	g.L.Info("Close Blink")
	return nil
}
