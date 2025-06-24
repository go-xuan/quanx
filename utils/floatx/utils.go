package floatx

import (
	"math"
)

// Ratio 计算百分率
func Ratio(numerator, denominator, prec int) float64 {
	if denominator > 0 {
		ratio := float64(numerator) * 100 / float64(denominator)
		return Ground(ratio, prec)
	}
	return 0
}

// Ground 四舍五入
func Ground(target float64, prec int) float64 {
	if target == 0 {
		return target
	} else {
		pow := math.Pow10(prec)
		return math.Floor(target*pow+0.5) / pow
	}
}

// Min 三数取小
func Min(a, b, c float64) float64 {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

// Max 三数取大
func Max(a, b, c float64) float64 {
	if a >= b && a >= c {
		return a
	} else if b >= a && b >= c {
		return b
	}
	return c
}
