package randx

// 随机bool
func Bool() bool {
	return IntRange(0, 1) == 0
}
