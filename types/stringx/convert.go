package stringx

import (
	"strconv"
)

func ToInt(s string) int {
	if value, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return value
	}
}

func ToInt64(s string) int64 {
	if value, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return value
	}
}

func ToFloat64(s string) float64 {
	if value, err := strconv.ParseFloat(s, 64); err != nil {
		return 0
	} else {
		return value
	}
}

func ParseInt(i int) string {
	return strconv.Itoa(i)
}

func ParseInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ParseFloat64(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
