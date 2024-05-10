package randx

import "strings"

// 随机bool
func Bool() bool {
	return IntRange(0, 1) == 0
}

// 枚举
func Enum(enums, sep string) string {
	list := strings.Split(enums, sep)
	return list[IntRange(0, len(list)-1)]
}

// 单选
func Radio(options []string) string {
	return options[IntRange(0, len(options)-1)]
}
