package stringx

import (
	"strings"
)

// 字符串是否包含
func ContainsAny(str string, args ...string) bool {
	for _, sub := range args {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}

// 字符串是否包含
func ContainsBoth(str string, args ...string) bool {
	for _, sub := range args {
		if !strings.Contains(str, sub) {
			return false
		}
	}
	return true
}

// 是否有空
func HasEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return true
		}
	}
	return false
}

// 反转
func Reverse(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// 字符串截取
func SubString(str string, start, end int) string {
	var r = []rune(str)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return str
	}
	return string(r[start:end])
}

// 将字符以最后一个符号拆分
func SplitByFirst(str, sep string) (string, string) {
	if str == "" {
		return "", ""
	}
	if strings.Contains(str, sep) && sep != "" {
		i := strings.Index(str, sep)
		return str[:i], str[i+len(sep):]
	} else {
		return str, ""
	}
}

// 将字符以最后一个符号拆分
func SplitByLast(str, sep string) (string, string) {
	if str == "" {
		return "", ""
	}
	if strings.Contains(str, sep) && sep != "" {
		i := strings.LastIndex(str, sep)
		return str[:i], str[i+len(sep):]
	} else {
		return str, ""
	}
}

// 字符填充(将字符以固定长度填充)
func StringFill(str, add string, length int, onLeft bool) string {
	strLen := len(str)
	addLen := len(add)
	fillLen := length - strLen
	if fillLen <= 0 && addLen == 0 {
		return str
	}
	fillStr := strings.Builder{}
	for i := 0; i < fillLen; i++ {
		fillStr.WriteString(string(add[i%addLen]))
	}
	if onLeft {
		return fillStr.String() + str
	} else {
		return str + fillStr.String()
	}
}

// 转下划线
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 转小驼峰
func LowerCamelCase(s string) string {
	u := UpperCamelCase(s)
	return strings.ToLower(string(u[0])) + u[1:]
}

// 转大驼峰
func UpperCamelCase(s string) string {
	s = strings.ToLower(s)
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		if k && d == '_' && num > i && s[i+1] >= '0' && s[i+1] <= '9' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 文本相似度计算
func TextSimilarity(source, target string) float64 {
	sLen, tLen := len(source), len(target)
	if (sLen == 0 && tLen == 0) || source == target {
		return 1.0
	}
	matrix := make([][]int, sLen+1)
	for i := range matrix {
		matrix[i] = make([]int, tLen+1)
		matrix[i][0] = i
	}

	for j := 0; j <= tLen; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= sLen; i++ {
		for j := 1; j <= tLen; j++ {
			cost := 0
			if source[i-1] != target[j-1] {
				cost = 1
			}
			matrix[i][j] = min(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	distance := matrix[sLen][tLen]
	maxLen := float64(sLen)
	if tLen > sLen {
		maxLen = float64(tLen)
	}
	return 1.0 - float64(distance)/maxLen
}

func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}
