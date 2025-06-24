package randx

import (
	"github.com/go-xuan/quanx/utils/floatx"
)

// Int 随机整数
func Int() int {
	return IntRange(1, 1<<16)
}

// Int64 随机整数
func Int64() int64 {
	return Int64Range(1, 1<<16)
}

// Float64 随机浮点数
func Float64() float64 {
	return Float64Range(1, 1<<16, 6)
}

// IntRange 随机整数
func IntRange(min, max int) int {
	if min != max {
		return NewRand().Intn(max-min+1) + min
	}
	return min
}

// Int64Range 随机整数
func Int64Range(min, max int64) int64 {
	if min != max {
		return NewRand().Int63n(max-min+1) + min
	}
	return min
}

// Float64Range 随机浮点数
func Float64Range(min, max float64, prec int) float64 {
	if min != max {
		float := NewRand().Float64()*(max-min+1) + min
		return floatx.Ground(float, prec)
	}
	return min
}
