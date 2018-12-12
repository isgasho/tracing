package service

import (
	"sync"
)

// Cache ...
type Cache struct {
	sync.RWMutex
	App map[string]map[int]string
}

// NewCache ...
func NewCache() *Cache {
	return &Cache{
		App: make(map[string]map[int]string),
	}
}
