package util

import "sync"

// App ...
type App struct {
	AppID  int32  `db:"id" json:"id" msg:"id"`
	Name   string `db:"name" json:"name" msg:"name"`
	Agents sync.Map
	Apis   sync.Map
}

// NewApp ...
func NewApp() *App {
	return &App{}
}
