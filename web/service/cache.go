package service

import "time"

type cache struct {
	appList       []*AppStat
	appListUpdate time.Time
}
