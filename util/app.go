package util

// App ...
type App struct {
	Code int32  `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}

// NewApp ...
func NewApp() *App {
	return &App{}
}
