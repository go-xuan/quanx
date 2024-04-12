package randx

import (
	"math"
	
	"github.com/go-xuan/quanx/utils/floatx"
)

// 随机整数
func Int() int {
	return IntRange(1, math.MaxInt)
}

// 随机整数
func Int64() int64 {
	return Int64Range(1, math.MaxInt64)
}

// 随机浮点数
func Float64() float64 {
	return Float64Range(math.SmallestNonzeroFloat64, math.MaxFloat64, 6)
}

// 随机整数
func IntRange(min, max int) int {
	if min == max {
		return min
	}
	return NewRand().Intn(1+max-min) + min
}

// 随机整数
func Int64Range(min, max int64) int64 {
	if min == max {
		return min
	}
	return NewRand().Int63n(1+max-min) + min
}

// 随机浮点数
func Float64Range(min, max float64, prec int) float64 {
	if min == max {
		return min
	}
	float := NewRand().Float64()*(max-min) + min
	return floatx.Ground(float, prec)
}
