package service

import (
	"github.com/mafanr/g"
)

// Blink ...
type Blink interface {
	Start() error
	Close() error
}

// blink 实时计算
type blink struct {
}

// newBlink ...
func newBlink() *blink {
	return &blink{}
}

// Start ...
func (s *blink) Start() error {
	g.L.Info("Start blink")
	return nil
}

// Close ...
func (s *blink) Close() error {
	g.L.Info("Close blink")
	return nil
}
