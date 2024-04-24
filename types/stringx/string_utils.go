package stringx

import (
	"strings"

	"github.com/go-xuan/quanx/types/anyx"
)

// 获取起始字符首次出现和结尾字符末次出现的下标
func Between(s, start, end string) (from, to int) {
	if start == end {
		indices := AllIndex(s, start, 2)
		from, to = indices[0], indices[1]
		return
	}
	from, to = -1, -1
	var l, m, n = len(s), len(start), len(end)
	if m > l || n > l {
		return
	}
	// x:start个数  y:end个数
	var x, y int
	for i := 0; i < l; i++ {
		if s[i] == start[0] {
			if s[i:i+m] == start {
				x++
				if x == 1 {
					from = i + m
				}
				i = i + m - 1
			}
		} else if s[i] == end[0] {
			if s[i:i+n] == end {
				y++
				if y == x || x == 1 {
					to = i
					break
				}
				i = i + n - 1
			}
		}
	}
	if to == -1 {
		from = -1
	}
	return
}

// 获取子串下标
// position：表示获取位置，默认position=1即正序第1处，position=-1即倒序第1处
func Index(s, sep string, position ...int) int {
	l, m, n := len(s), len(sep), 0
	if m > l {
		return -1
	}
	x := anyx.Default(position, 1)
	for i := 0; i <= l-m; i++ {
		if x > 0 {
			if s[i] == sep[0] && s[i:i+m] == sep {
				n++
				if x == n {
					return i
				}
			}
		} else {
			j := l - i
			if s[j-1] == sep[m-1] && s[j-m:j] == sep {
				n--
				if x == n {
					return j - m
				}
			}
		}
	}
	return -1
}

// 获取所有下标
func AllIndex(s, sep string, x int) []int {
	var indices = make([]int, x)
	for i := 0; i < x; i++ {
		indices[i] = -1
	}
	l, m, n := len(s), len(sep), 0
	for i := 0; i <= l-m; i++ {
		if s[i] == sep[0] && s[i:i+m] == sep {
			n++
			if n <= x {
				indices[n-1] = i
			}
			i = i + m - 1
		}
	}
	return indices
}

// 添加前缀
func AddPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// 添加前缀
func AddSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// 字符串是否包含
func Contains(s string, seps ...string) (string, int) {
	var hit, index = "", -1
	for _, sep := range seps {
		if i := Index(s, sep); i >= 0 {
			if i < index {
				hit, index = sep, i
			} else if index == -1 {
				hit, index = sep, i
			}
		}
	}
	return hit, index
}

// 字符串是否包含
func ContainsAny(s string, seps ...string) bool {
	for _, sep := range seps {
		if Index(s, sep) >= 0 {
			return true
		}
	}
	return false
}

// 字符串是否包含
func ContainsBoth(s string, seps ...string) bool {
	for _, sep := range seps {
		if Index(s, sep) == -1 {
			return false
		}
	}
	return true
}

// 是否有空
func HasEmpty(s ...string) bool {
	for _, item := range s {
		if item == "" {
			return true
		}
	}
	return false
}

// 不等时取默认值
func IfNot(s, v, def string) string {
	if s != v {
		return def
	} else {
		return s
	}
}

// 为空时取默认值
func IfZero(s, def string) string {
	if s == "" {
		return def
	} else {
		return s
	}
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
func SubString(s string, start, end int) string {
	var r = []rune(s)
	length := len(r)
	if start < 0 || end > length || start > end {
		return ""
	}
	if start == 0 && end == length {
		return s
	}
	return string(r[start:end])
}

// 分割字符串（reverse=true从右往左）
// position：表示分割位置，默认position=1即正序第1处，position=-1即倒序第1处
func Cut(s, sep string, position ...int) (string, string) {
	if i := Index(s, sep, position...); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return s, ""
}

// 插入字符串
func Insert(s, insert string, position ...int) string {
	if len(position) > 0 {
		if i := position[0]; position[0] > 0 && position[0] < len(s) {
			return s[:i] + insert + s[i:]
		}
	}
	return s + insert
}

// 字符以固定长度填充（默认填充左边）
func Fill(s, fill string, length int, right ...bool) string {
	strLen, addLen := len(s), len(fill)
	fillLen := length - strLen
	if fillLen <= 0 && addLen == 0 {
		return s
	}
	fillStr := strings.Builder{}
	for i := 0; i < fillLen; i++ {
		fillStr.WriteString(string(fill[i%addLen]))
	}
	return anyx.If(len(right) > 0 && right[0], s+fillStr.String(), fillStr.String()+s)
}

// 转下划线
func ToSnake(s string) string {
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
func ToLowerCamel(s string) string {
	ucc := ToUpperCamel(s)
	return strings.ToLower(string(ucc[0])) + ucc[1:]
}

// 转大驼峰
func ToUpperCamel(s string) string {
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
			matrix[i][j] = MinOfThree(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	distance := matrix[sLen][tLen]
	maxLen := float64(sLen)
	if tLen > sLen {
		maxLen = float64(tLen)
	}
	return 1.0 - float64(distance)/maxLen
}

func MinOfThree(a, b, c int) int {
	if a <= b && a <= c {
		return a
	} else if b <= a && b <= c {
		return b
	}
	return c
}

const (
	Upper      = "upper"      // 大写
	Lower      = "lower"      // 小写
	UpperCamel = "upperCamel" // 大驼峰
	LowerCamel = "lowerCamel" // 小驼峰
	Snake      = "snake"      // 蛇形
)

func Transforms(str string, mode ...string) map[string]string {
	var result = make(map[string]string)
	for _, m := range mode {
		result[m] = Transform(str, m)
	}
	return result
}

func Transform(str string, mode string) string {
	switch mode {
	case Upper:
		return strings.ToUpper(str)
	case Lower:
		return strings.ToLower(str)
	case UpperCamel:
		return ToUpperCamel(str)
	case LowerCamel:
		return ToLowerCamel(str)
	case Snake:
		return ToSnake(str)
	default:
		return str
	}
}
