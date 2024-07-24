package intx

import "strconv"

// IfZero 为空时取默认值
func IfZero(s, def int) int {
	if s == 0 {
		return def
	} else {
		return s
	}
}

// IfNot 不等时取默认值
func IfNot(s, v, def int) int {
	if s != v {
		return def
	} else {
		return s
	}
}

// Min 三数取小
func Min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

// Max 三数取大
func Max(a, b, c int) int {
	if a >= b && a >= c {
		return a
	} else if b >= a && b >= c {
		return b
	}
	return c
}

func String(i int) string {
	return strconv.Itoa(i)
}
