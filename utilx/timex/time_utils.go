package timex

import (
	"github.com/go-xuan/quanx/utilx/anyx"
	"strconv"
	"strings"
	"time"
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
func ToString(time time.Time) string {
	return time.Format(TimeFmt)
}

// 字符转时间
func ToTime(timeStr string) time.Time {
	format := anyx.IfElseValue(strings.Contains(timeStr, "-") && len(timeStr) == 10, DateFmt, TimeFmt)
	return TimeFormat(timeStr, format)
}

// 时间格式化
func TimeFormat(timeStr string, format string) time.Time {
	parseTime, err := time.ParseInLocation(format, timeStr, time.Local)
	if err != nil {
		return time.Unix(0, 0)
	}
	return parseTime
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
	return DateStart(time.Now())
}

// 今天开始时间
func DateStart(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
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
	diff := (year - 4) % 12
	animals := strings.Split(AllShengXiao, ",")
	return animals[diff]
}

// 间隔时间
func TimeDiff(startTime, endTime time.Time, unit Unit) (diff int64) {
	if startTime.After(endTime) {
		temp := startTime
		startTime = endTime
		endTime = temp
	}
	switch unit {
	case Year:
		diff = int64(YearInterval(startTime, endTime))
	case Month:
		diff = MonthInterval(startTime, endTime)
	case Day:
		diff = (endTime.Unix() - startTime.Unix()) / 86400
	case Hour:
		diff = (endTime.Unix() - startTime.Unix()) / 3600
	case Minute:
		diff = (endTime.Unix() - startTime.Unix()) / 60
	case Second:
		diff = endTime.Unix() - startTime.Unix()
	case Milli:
		diff = endTime.UnixMilli() - startTime.UnixMilli()
	case Micro:
		diff = endTime.UnixMicro() - startTime.UnixMicro()
	case Nano:
		diff = endTime.UnixNano() - startTime.UnixNano()
	default:
	}
	return
}

// 间隔年数
func YearInterval(startTime, endTime time.Time) (diff int) {
	if startTime.After(endTime) {
		return 0
	}
	return endTime.Year() - startTime.Year()
}

// 间隔月份数
func MonthInterval(startTime, endTime time.Time) (diff int64) {
	startTime = startTime.AddDate(0, 1, 0)
	for {
		if startTime.Before(endTime) {
			startTime = startTime.AddDate(0, 1, 0)
			diff++
		} else {
			break
		}
	}
	return
}

// 间隔天数
func DayDiff(startTime, endTime time.Time) (diff int64) {
	startTime = startTime.AddDate(0, 0, 1)
	for {
		if startTime.Before(endTime) {
			startTime = startTime.AddDate(0, 0, 1)
			diff++
		} else {
			break
		}
	}
	return
}

// 根据时间范围生成月份切片
func MonthSlice(startTimeUnix, endTimeUnix int64) []string {
	var captureMonths []string
	tempTime := time.Unix(startTimeUnix/1000, 0)
	endTime := time.Unix(endTimeUnix/1000, 0)
	year, month, _ := tempTime.Date()
	tempTime = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	for {
		if tempTime.Before(endTime) {
			captureMonths = append(captureMonths, tempTime.Format(MonthFmt))
			tempTime = tempTime.AddDate(0, 1, 0)
		} else if tempTime.Equal(endTime) {
			captureMonths = append(captureMonths, tempTime.Format(MonthFmt))
			break
		} else {
			break
		}
	}
	return captureMonths
}

// 生成日期切片
func DateSlice(startTimeUnix, endTimeUnix int64) []string {
	var dateList []string
	// 获取开始时间的当天0点0分0秒
	year, month, day := time.Unix(startTimeUnix/1000, 0).Date()
	tempTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	// 获取相差天数
	timeDiff, _ := strconv.Atoi(strconv.FormatInt((endTimeUnix-tempTime.Unix())/86400, 10))
	for i := 0; i <= timeDiff; i++ {
		dateList = append(dateList, tempTime.Format(DateFmt))
		tempTime = tempTime.AddDate(0, 0, 1)
	}
	return dateList
}

// 获取特定范围内起止时间(当天/本周/本月/本年)
func TimeRange(unit Unit) (startTime time.Time, endTime time.Time) {
	now := time.Now()
	switch unit {
	case Year:
		startTime = DateStart(now.AddDate(0, 0, -now.YearDay()+1))
		endTime = startTime.AddDate(1, 0, 0).Add(-time.Second)
	case Month:
		startTime = DateStart(now.AddDate(0, 0, -now.Day()+1))
		endTime = startTime.AddDate(0, 1, 0).Add(-time.Second)
	case Week:
		startTime = DateStart(now.AddDate(0, 0, int(time.Monday-now.Weekday())))
		endTime = startTime.AddDate(0, 0, 7).Add(-time.Second)
	case Day:
		startTime = DateStart(now)
		endTime = startTime.AddDate(0, 0, 1).Add(-time.Second)
	default:
	}
	return
}
