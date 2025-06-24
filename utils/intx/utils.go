package intx

// Default 用于函数中的不定参数取默认值
func Default(def int, x ...int) int {
	if len(x) == 0 {
		return def
	} else {
		return x[0]
	}
}

// IfZero 为空时取默认值
func IfZero(v, def int) int {
	if v == 0 {
		return def
	} else {
		return v
	}
}

// IfNot 不等时取默认值
func IfNot(v, neq, def int) int {
	if v != neq {
		return def
	} else {
		return v
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
