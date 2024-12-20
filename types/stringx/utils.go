package stringx

import (
	"strconv"
	"strings"
	"time"
)

func ParseInt(s string, def ...int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseInt64(s string, def ...int64) int64 {
	if value, err := strconv.ParseInt(s, 10, 64); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseFloat(s string, def ...float64) float64 {
	if value, err := strconv.ParseFloat(s, 64); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseBool(value string, def ...bool) bool {
	switch value {
	case "1", "t", "T", "true", "TRUE", "True", "是", "yes", "YES", "Yes":
		return true
	default:
		if len(def) > 0 {
			return def[0]
		} else {
			return false
		}
	}
}

func ParseTime(value string, def ...time.Time) time.Time {
	if len(value) == 10 && value[4:5] == "-" {
		if location, err := time.ParseInLocation("2006-01-02", value, time.Local); err == nil {
			return location
		}
	} else if location, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local); err == nil {
		return location
	} else if len(def) > 0 {
		return def[0]
	}
	return time.Time{}
}

func FormatInt(i int) string {
	return strconv.Itoa(i)
}

func FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// Between 获取起始字符首次出现和结尾字符末次出现的下标
func Between(s, start, end string) (from, to int) {
	from, to = -1, -1
	if start == end {
		if indices := Indices(s, start); len(indices) > 1 {
			from, to = indices[0], indices[1]
		} else if len(indices) == 1 {
			from = indices[0]
		}
		return
	}
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
		}
		if s[i] == end[0] {
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

// Indices 获取所有下标, x：命中数量
func Indices(s, sep string) []int {
	var indices []int
	l, m, n := len(s), len(sep), 0
	for i := 0; i <= l-m; i++ {
		if s[i] == sep[0] && s[i:i+m] == sep {
			indices = append(indices, i)
			n++
			i = i + m - 1
		}
	}
	return indices
}

// Index 获取子串的下标
// position：表示获取位置，默认position=1即正序第1处，position=-1即倒序第1处
func Index(s, sep string, position ...int) int {
	if l, m := len(s), len(sep); l >= m {
		var x, y = 1, 0 // x：目标获取位置，y：sep出现次数计数
		if len(position) > 0 {
			x = position[0]
		}
		for i := 0; i <= l-m; i++ {
			if x > 0 {
				if s[i] == sep[0] && s[i:i+m] == sep {
					y++
					if x == y {
						return i
					}
				}
			} else {
				j := l - i
				if s[j-1] == sep[m-1] && s[j-m:j] == sep {
					y--
					if x == y {
						return j - m
					}
				}
			}
		}
	}
	return -1
}

// IndexStrict 获取子串下标（严格模式：忽略单词中的子串）
func IndexStrict(s, sep string) int {
	kl, loop, index := len(sep), true, 0
	for loop {
		if newIndex := Index(s, sep, 1); newIndex >= 0 {
			sl := len(s)
			if newIndex == 0 && s[kl:kl+1] == " " {
				index, loop = index+newIndex, false
			} else if newIndex == sl-kl && s[newIndex-1:newIndex] == " " {
				index, loop = index+newIndex, false
			} else if s[newIndex-1:newIndex] == " " && s[newIndex+kl:newIndex+kl+1] == " " {
				index, loop = index+newIndex, false
			} else {
				// 当前index无效则缩减原sql继续loop
				index = newIndex + kl
				s = s[index:]
			}
		} else {
			index, loop = -1, false // 没找到直接跳出
		}
	}
	return index
}

// AddPrefix 添加前缀
func AddPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

// AddSuffix 添加前缀
func AddSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s
	}
	return s + suffix
}

// Contains 字符串是否包含
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

// ContainsAny 字符串是否包含
func ContainsAny(s string, seps ...string) bool {
	for _, sep := range seps {
		if Index(s, sep) >= 0 {
			return true
		}
	}
	return false
}

// ContainsBoth 字符串是否包含
func ContainsBoth(s string, seps ...string) bool {
	for _, sep := range seps {
		if Index(s, sep) == -1 {
			return false
		}
	}
	return true
}

// HasEmpty 是否有空
func HasEmpty(s ...string) bool {
	for _, item := range s {
		if item == "" {
			return true
		}
	}
	return false
}

// Default 用于函数中的不定参数取默认值
func Default(def string, x ...string) string {
	if len(x) == 0 {
		return def
	} else {
		return x[0]
	}
}

// IfNot 不等时取默认值
func IfNot(s, v, def string) string {
	if s != v {
		return def
	} else {
		return s
	}
}

// IfZero 为空时取默认值
func IfZero(s, def string) string {
	if s == "" {
		return def
	} else {
		return s
	}
}

// Reverse 反转
func Reverse(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// SubString 字符串截取
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

// Cut 分割字符串（reverse=true从右往左）
// position：表示分割位置，默认position=1即正序第1处，position=-1即倒序第1处
func Cut(s, sep string, position ...int) (string, string) {
	if i := Index(s, sep, position...); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return s, ""
}

// Insert 插入字符串
func Insert(s, insert string, position ...int) string {
	if len(position) > 0 {
		if i := position[0]; i > 0 && i < len(s) {
			return s[:i] + insert + s[i:]
		}
	}
	return s + insert
}

// Fill 字符以固定长度填充（默认填充左边）
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
	if len(right) > 0 && right[0] {
		return s + fillStr.String()
	} else {
		return fillStr.String() + s
	}
}

// ParseUrlParams 解析url参数为map
func ParseUrlParams(args string) map[string]string {
	if strings.Contains(args, "=") {
		var params = make(map[string]string)
		kvs := strings.Split(args, "&")
		for _, kv := range kvs {
			k, v := Cut(kv, "=")
			params[k] = v
		}
		return params
	}
	return nil
}

// ToSnake 转下划线
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

// ToLowerCamel 转小驼峰
func ToLowerCamel(s string) string {
	ucc := ToUpperCamel(s)
	return strings.ToLower(string(ucc[0])) + ucc[1:]
}

// ToUpperCamel 转大驼峰
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

// TextSimilarity 文本相似度计算
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
			matrix[i][j] = minInThree(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	distance := matrix[sLen][tLen]
	maxLen := float64(sLen)
	if tLen > sLen {
		maxLen = float64(tLen)
	}
	return 1.0 - float64(distance)/maxLen
}

func minInThree(a, b, c int) int {
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
