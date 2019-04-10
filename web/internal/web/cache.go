package service

import (
	"time"

	app "github.com/imdevlab/tracing/web/internal/application"
)

type cache struct {
	appList       []*app.Stat
	appListUpdate time.Time
}
