package mathx

import "math"

// 计算比率
func Ratio(numerator, denominator, prec int) float64 {
	if denominator > 0 {
		ratio := float64(numerator) * 100 / float64(denominator)
		return Ground(ratio, prec)
	}
	return 0
}

// 四舍五入
func Ground(target float64, prec int) float64 {
	if target == 0 || prec == 0 {
		return target
	} else {
		p := math.Pow10(prec)
		return math.Floor(target*p+0.5) / p
	}
}
