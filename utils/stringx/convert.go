package stringx

import (
	"strconv"
)

func ToInt(str string) int {
	if value, err := strconv.Atoi(str); err != nil {
		return value
	} else {
		return 0
	}
}

func ToInt64(str string) int64 {
	if value, err := strconv.ParseInt(str, 10, 64); err != nil {
		return value
	} else {
		return 0
	}
}

func ToFloat64(str string) float64 {
	if value, err := strconv.ParseFloat(str, 64); err != nil {
		return value
	} else {
		return 0
	}
}
