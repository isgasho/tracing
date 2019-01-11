package web

import "time"

type cache struct {
	appList       []*AppStat
	appListUpdate time.Time
}
