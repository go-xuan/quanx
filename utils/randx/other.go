package randx

import (
	"strings"
	"time"
)

// Bool 随机bool
func Bool() bool {
	return IntRange(0, 1) == 0
}

// Split 分割取值
func Split(str, sep string) string {
	list := strings.Split(str, sep)
	return list[IntRange(0, len(list)-1)]
}

// SelectByte 选择字节
func SelectByte(str string) byte {
	return str[IntRange(0, len(str)-1)]
}

// Enum 枚举
func Enum(enums []string) string {
	return enums[IntRange(0, len(enums)-1)]
}

// Time 一天内随机时间
func Time() time.Time {
	now := time.Now()
	return TimeRange(now.Add(time.Hour*-24), now)
}

// Date 随机日期(从1970-01-01到当前日期)
func Date() string {
	diff := Int64Range(1, time.Now().Unix())
	return time.Unix(diff, 0).Format("2006-01-02")
}

// TimeRange 随机时间
func TimeRange(min, max time.Time) time.Time {
	if max.After(min) {
		return time.Unix(0, Int64Range(min.UnixNano(), max.UnixNano()))
	}
	return min
}

// DateRange 随机时间
func DateRange(min, max time.Time) string {
	if max.After(min) {
		return time.Unix(0, Int64Range(min.UnixNano(), max.UnixNano())).Format("2006-01-02")
	}
	return min.Format("2006-01-02")
}
