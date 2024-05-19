package timex

import (
	"strings"
	"time"

	"github.com/go-xuan/quanx/types/anyx"
)

type Unit uint

const (
	TimestampFmt = "20060102150405"
	TimeFmt      = "2006-01-02 15:04:05"
	DateFmt      = "2006-01-02"
	MonthFmt     = "2006-01"
	AllShengXiao = "鼠,牛,虎,兔,龙,蛇,马,羊,猴,鸡,狗,猪"
)

const (
	Nano Unit = iota
	Micro
	Milli
	Second
	Minute
	Hour
	Day
	Week
	Month
	Year
)

// 时间转字符
func Format(time time.Time, format ...string) string {
	return time.Format(anyx.Default(TimeFmt, format...))
}

// 字符转时间
func ToTime(timeStr string) time.Time {
	var format = anyx.If(strings.Contains(timeStr, "-") && len(timeStr) == 10, DateFmt, TimeFmt)
	return TimeParse(timeStr, format)
}

// 时间格式化
func TimeParse(timeStr string, format ...string) time.Time {
	var layout = anyx.Default(TimeFmt, format...)
	if parseTime, err := time.ParseInLocation(layout, timeStr, time.Local); err != nil {
		return time.Unix(0, 0)
	} else {
		return parseTime
	}
}

// 时间戳(秒级)转字符
func SecondFormat(second int64, format string) string {
	return time.Unix(second, 0).Format(format)
}

// 当前时间字符串
func NowString() string {
	return time.Now().Format(TimeFmt)
}

// 今天
func TodayStr() string {
	return time.Now().Format(DateFmt)
}

// 今天开始时间
func TodayStart() time.Time {
	return DayStart(time.Now())
}

// 当天开始时间（yyyy-mm-dd 00:00:00）
func DayStart(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// 当天结束时间（yyyy-mm-dd 23:59:59）
func DayEnd(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, 0, time.Local)
}

// 昨天
func YesterdayStr() string {
	return time.Now().AddDate(0, 0, -1).Format(DateFmt)
}

// 是否闰年
func IsLeapYear(year int) bool {
	if year <= 0 {
		return false
	}
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		return true
	}
	return false
}

// 获取生肖
func ShengXiao(year int) string {
	for {
		if year < 4 {
			year = year + 12
		} else {
			break
		}
	}
	return strings.Split(AllShengXiao, ",")[(year-4)%12]
}

// 间隔时间
func TimeDiff(start, end time.Time, unit Unit) int64 {
	switch unit {
	case Year:
		return int64(YearInterval(start, end))
	case Month:
		return int64(MonthInterval(start, end))
	case Day:
		return (end.Unix() - start.Unix()) / 86400
	case Hour:
		return (end.Unix() - start.Unix()) / 3600
	case Minute:
		return (end.Unix() - start.Unix()) / 60
	case Second:
		return end.Unix() - start.Unix()
	case Milli:
		return end.UnixMilli() - start.UnixMilli()
	case Micro:
		return end.UnixMicro() - start.UnixMicro()
	case Nano:
		return end.UnixNano() - start.UnixNano()
	default:
		return -1
	}
}

// 间隔天数
func DayInterval(start, end time.Time) int {
	return int((end.Unix() - start.Unix()) / 86400)
}

// 间隔月份数
func MonthInterval(start, end time.Time) int {
	y1, m1, d1 := start.Date()
	y2, m2, d2 := end.Date()
	diff := (y2-y1)*12 + int(m2-m1)
	if d1 <= d2 {
		diff++
	}
	return diff
}

// 间隔年数
func YearInterval(start, end time.Time) int {
	return end.Year() - start.Year()
}

// 时间切片
func TimeSlice(start, end time.Time, unit Unit) []string {
	var slice []string
	if unit == Day {
		for start.Unix() <= end.Unix() {
			slice = append(slice, start.Format(DateFmt))
			start = start.AddDate(0, 0, 1)
		}
	} else if unit == Month {
		for start.Unix() <= end.Unix() {
			slice = append(slice, start.Format(MonthFmt))
			start = start.AddDate(0, 1, 0)
		}
	}
	return slice
}

// 获取特定范围内起止时间(当天/本周/本月/本年)
func TimeRange(t time.Time, unit Unit) (start, end time.Time) {
	switch unit {
	case Year:
		start = DayStart(t.AddDate(0, 0, -t.YearDay()+1))
		end = start.AddDate(1, 0, 0).Add(-time.Second)
	case Month:
		start = DayStart(t.AddDate(0, 0, -t.Day()+1))
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	case Week:
		start = DayStart(t.AddDate(0, 0, int(time.Monday-t.Weekday())))
		end = start.AddDate(0, 0, 7).Add(-time.Second)
	default:
		start = DayStart(t)
		end = start.AddDate(0, 0, 1).Add(-time.Second)
	}
	return
}
