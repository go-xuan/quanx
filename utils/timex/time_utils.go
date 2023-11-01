package timex

import (
	"strconv"
	"strings"
	"time"

	"github.com/quanxiaoxuan/quanx/common/constx"
)

type Unit uint

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
	return time.Format(constx.TimeFmt)
}

// 字符转时间
func ToTime(timeStr string) time.Time {
	format := constx.TimeFmt
	if strings.Contains(timeStr, "-") && len(timeStr) == 10 {
		format = constx.DateFmt
	}
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
	return time.Now().Format(constx.TimeFmt)
}

// 今天
func TodayStr() string {
	return time.Now().Format(constx.DateFmt)
}

// 今天开始时间
func TodayStart() time.Time {
	return time.Now().Truncate(24 * time.Hour).Add(-8 * time.Hour)
}

// 昨天
func YesterdayStr() string {
	return time.Now().AddDate(0, 0, -1).Format(constx.DateFmt)
}

// 获取生肖
func AnimalIs(year int) string {
	for {
		if year < 4 {
			year = year + 12
		} else {
			break
		}
	}
	diff := (year - 4) % 12
	animals := strings.Split(constx.Animal, ",")
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
	case Nano:
		diff = endTime.UnixNano() - startTime.UnixNano()
	case Micro:
		diff = endTime.UnixMicro() - startTime.UnixMicro()
	case Milli:
		diff = endTime.UnixMilli() - startTime.UnixMilli()
	case Second:
		diff = endTime.Unix() - startTime.Unix()
	case Minute:
		diff = (endTime.Unix() - startTime.Unix()) / 60
	case Hour:
		diff = (endTime.Unix() - startTime.Unix()) / 3600
	case Day:
		diff = (endTime.Unix() - startTime.Unix()) / 86400
	case Month:
		diff = MonthInterval(startTime, endTime)
	case Year:
		diff = int64(YearInterval(startTime, endTime))
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
			captureMonths = append(captureMonths, tempTime.Format(constx.MonthFmt))
			tempTime = tempTime.AddDate(0, 1, 0)
		} else if tempTime.Equal(endTime) {
			captureMonths = append(captureMonths, tempTime.Format(constx.MonthFmt))
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
		dateList = append(dateList, tempTime.Format(constx.DateFmt))
		tempTime = tempTime.AddDate(0, 0, 1)
	}
	return dateList
}

// 获取特定范围内起止时间(当天/本周/本月/本年)
func TimeRange(unit Unit) (startTime time.Time, endTime time.Time) {
	now := time.Now()
	switch unit {
	case Day:
		startTime = now.Truncate(24 * time.Hour).Add(-8 * time.Hour)
		endTime = startTime.AddDate(0, 0, 1).Add(-time.Second)
	case Week:
		startTime = now.AddDate(0, 0, int(time.Monday-now.Weekday())).Truncate(24 * time.Hour).Add(-8 * time.Hour)
		endTime = startTime.AddDate(0, 0, 7).Add(-time.Second)
	case Month:
		startTime = now.AddDate(0, 0, -now.Day()+1).Truncate(24 * time.Hour).Add(-8 * time.Hour)
		endTime = startTime.AddDate(0, 1, 0).Add(-time.Second)
	case Year:
		startTime = now.AddDate(0, 0, -now.YearDay()+1).Truncate(24 * time.Hour).Add(-8 * time.Hour)
		endTime = startTime.AddDate(1, 0, 0).Add(-time.Second)
	}
	return
}
