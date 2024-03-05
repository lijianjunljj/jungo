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

func MondayStart() time.Time {
	now := time.Now().UTC()                                        // 获取当前时间并转换为UTC时区
	weekday := now.Weekday()                                       // 获取今天是星期几
	delta := int(time.Monday-weekday+7) % 7                        // 计算距离本周一还有多少天
	monday := now.AddDate(0, 0, delta*-1).Truncate(24 * time.Hour) // 将当前时间向前调整到本周一零点
	return monday
}
