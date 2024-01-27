package jun_util

import (
	"time"
)

func TodayStartTimeUnix() int64 {
	now := time.Now()
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return today.Unix()
}
