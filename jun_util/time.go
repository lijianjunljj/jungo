package jun_util

import (
	"github.com/lijianjunljj/gocommon/config"
	"time"
)

const DateTimeFormat = config.DateTimeFormat

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
func TimeUnix2DateTime(second int64) string {
	return time.Unix(second, 0).Format(DateTimeFormat)
}

// 获取本周开始和结束时间
func GetMondayAndSundayUnixTime() (int64, int64) {
	// 获取本周第一天（周一）
	t := time.Now()
	weekDay := int(t.Weekday())

	if weekDay == 0 {
		weekDay = 7
	}

	monday := t.AddDate(0, 0, -weekDay+1)
	// 获取周一的零点时间
	zeroTime := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
	//monday := t.AddDate(0, 0, -int(t.Weekday()))
	//mondayStartOfDay := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.Local)
	unixMonday := zeroTime.Unix()

	// 获取本周最后一天（周日）
	sunday := t.AddDate(0, 0, -weekDay+8)
	sundayStartOfDay := time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 0, 0, 0, 0, time.Local)
	unixSunday := sundayStartOfDay.Unix()

	return unixMonday, unixSunday
}
