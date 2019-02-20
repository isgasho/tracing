package misc

import (
	"time"

	"github.com/labstack/echo"
	"github.com/mafanr/g/utils"
)

// 获取开始和截止日期
func StartEndDate(c echo.Context) (start time.Time, end time.Time, err error) {
	startRaw := c.FormValue("start")
	endRaw := c.FormValue("end")

	// start和end的时间字符串转成秒级时间戳:2019-01-10 00:00:00
	start, err = utils.StringToTime2(startRaw)
	if err != nil {
		return
	}

	end, err = utils.StringToTime2(endRaw)
	if err != nil {
		return
	}
	// utils.OnlyAlphaAndNum
	return
}
