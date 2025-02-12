package stringx

import (
	"strconv"
	"strings"
	"time"
)

func ParseInt(str string, def ...int) int {
	if value, err := strconv.Atoi(str); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseInt64(str string, def ...int64) int64 {
	if value, err := strconv.ParseInt(str, 10, 64); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseFloat(str string, def ...float64) float64 {
	if value, err := strconv.ParseFloat(str, 64); err == nil {
		return value
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func ParseBool(str string, def ...bool) bool {
	switch str {
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

func ParseTime(str string, def ...time.Time) time.Time {
	if len(str) == 10 && str[4:5] == "-" {
		if location, err := time.ParseInLocation("2006-01-02", str, time.Local); err == nil {
			return location
		}
	} else if location, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local); err == nil {
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
func Between(str, start, end string) (from, to int) {
	from, to = -1, -1
	if start == end {
		if indices := Indices(str, start); len(indices) > 1 {
			from, to = indices[0], indices[1]
		} else if len(indices) == 1 {
			from = indices[0]
		}
		return
	}
	var l, sl, el = len(str), len(start), len(end)
	if sl > l || el > l {
		return
	}
	// x:start个数  y:end个数
	var x, y int
	for i := 0; i < l; i++ {
		if str[i] == start[0] {
			if str[i:i+sl] == start {
				x++
				if x == 1 {
					from = i
				}
				i = i + sl - 1
			}
		}
		if str[i] == end[0] {
			if str[i:i+el] == end {
				y++
				if y == x || x == 1 {
					to = i
					break
				}
				i = i + el - 1
			}
		}
	}
	if to == -1 {
		from = -1
	}
	return
}

// Indices 获取所有下标, x：命中数量
func Indices(str, key string) []int {
	var indices []int
	l, m, n := len(str), len(key), 0
	for i := 0; i <= l-m; i++ {
		if str[i] == key[0] && str[i:i+m] == key {
			indices = append(indices, i)
			n++
			i = i + m - 1
		}
	}
	return indices
}

// Index 获取子串的下标
// position：表示获取位置，默认position=1即正序第1处，position=-1即倒序第1处
func Index(str, key string, position ...int) int {
	if l, m := len(str), len(key); l >= m {
		var x, y = 1, 0 // x：目标获取位置，y：sep出现次数计数
		if len(position) > 0 {
			x = position[0]
		}
		for i := 0; i <= l-m; i++ {
			if x > 0 {
				if str[i] == key[0] && str[i:i+m] == key {
					y++
					if x == y {
						return i
					}
				}
			} else {
				j := l - i
				if str[j-1] == key[m-1] && str[j-m:j] == key {
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
func IndexStrict(str, key string) int {
	kl, loop, index := len(key), true, 0
	for loop {
		if newIndex := Index(str, key, 1); newIndex >= 0 {
			if HasAdjacent(str, key, " ", newIndex) {
				index, loop = index+newIndex, false
			} else {
				index = newIndex + kl
				str = str[index:]
			}
		} else {
			index, loop = -1, false // 没找到直接跳出
		}
	}
	return index
}

// HasAdjacent 判断目标kew在文本中当前位置是否有相邻字符
func HasAdjacent(str, key, adjacent string, index int) bool {
	sl, kl, al := len(str), len(key), len(adjacent)
	if index == 0 {
		return str[kl:kl+al] == adjacent
	} else if index == sl-kl {
		return str[index-al:index] == adjacent
	} else {
		return str[index-al:index] == adjacent && str[index+kl:index+kl+al] == adjacent
	}
}

// AddPrefix 添加前缀
func AddPrefix(str, prefix string) string {
	if strings.HasPrefix(str, prefix) {
		return str
	}
	return prefix + str
}

// AddSuffix 添加前缀
func AddSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str
	}
	return str + suffix
}

// Split 字符串分割
func Split(str string, sep string) []string {
	slice := strings.Split(str, sep)
	for i, s := range slice {
		slice[i] = strings.TrimSpace(s)
	}
	return slice
}

// Contains 字符串是否包含
func Contains(str string, keys ...string) (string, int) {
	var hit, index = "", -1
	for _, key := range keys {
		if i := Index(str, key); i >= 0 {
			if i < index {
				hit, index = key, i
			} else if index == -1 {
				hit, index = key, i
			}
		}
	}
	return hit, index
}

// ContainsAny 字符串是否包含
func ContainsAny(str string, keys ...string) bool {
	for _, key := range keys {
		if Index(str, key) >= 0 {
			return true
		}
	}
	return false
}

// ContainsBoth 字符串是否包含
func ContainsBoth(str string, seps ...string) bool {
	for _, sep := range seps {
		if Index(str, sep) == -1 {
			return false
		}
	}
	return true
}

// HasEmpty 是否有空
func HasEmpty(str ...string) bool {
	for _, item := range str {
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
func IfNot(str, value, def string) string {
	if str != value {
		return def
	} else {
		return str
	}
}

// IfZero 为空时取默认值
func IfZero(str, def string) string {
	if str == "" {
		return def
	} else {
		return str
	}
}

// Reverse 反转
func Reverse(str string) string {
	runes := []rune(str)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// SubString 字符串截取
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

// Cut 分割字符串（reverse=true从右往左）
// position：表示分割位置，默认position=1即正序第1处，position=-1即倒序第1处
func Cut(str, sep string, position ...int) (string, string) {
	if i := Index(str, sep, position...); i >= 0 {
		return str[:i], str[i+len(sep):]
	}
	return str, ""
}

// Insert 插入字符串
func Insert(str, insert string, position ...int) string {
	if len(position) > 0 {
		if i := position[0]; i > 0 && i < len(str) {
			return str[:i] + insert + str[i:]
		}
	}
	return str + insert
}

// Fill 字符以固定长度填充（默认填充左边）
func Fill(str, fill string, length int, right ...bool) string {
	strLen, addLen := len(str), len(fill)
	fillLen := length - strLen
	if fillLen <= 0 && addLen == 0 {
		return str
	}
	fillStr := strings.Builder{}
	for i := 0; i < fillLen; i++ {
		fillStr.WriteString(string(fill[i%addLen]))
	}
	if len(right) > 0 && right[0] {
		return str + fillStr.String()
	} else {
		return fillStr.String() + str
	}
}

// ParseUrlParams 解析url参数为map
func ParseUrlParams(str string) map[string]string {
	if str != "" {
		var params = make(map[string]string)
		var kvs = strings.Split(str, "&")
		for _, kv := range kvs {
			k, v := Cut(kv, "=")
			params[k] = v
		}
		return params
	}
	return nil
}

// ToSnake 转下划线
func ToSnake(str string) string {
	data := make([]byte, 0, len(str)*2)
	j := false
	num := len(str)
	for i := 0; i < num; i++ {
		d := str[i]
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
func ToLowerCamel(str string) string {
	ucc := ToUpperCamel(str)
	return strings.ToLower(string(ucc[0])) + ucc[1:]
}

// ToUpperCamel 转大驼峰
func ToUpperCamel(str string) string {
	str = strings.ToLower(str)
	data := make([]byte, 0, len(str))
	j := false
	k := false
	num := len(str) - 1
	for i := 0; i <= num; i++ {
		d := str[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && str[i+1] >= 'a' && str[i+1] <= 'z' {
			j = true
			continue
		}
		if k && d == '_' && num > i && str[i+1] >= '0' && str[i+1] <= '9' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// Similarity 文本相似度计算
func Similarity(source, target string) float64 {
	sl, tl := len(source), len(target)
	if (sl == 0 && tl == 0) || source == target {
		return 1.0
	}
	matrix := make([][]int, sl+1)
	for i := range matrix {
		matrix[i] = make([]int, tl+1)
		matrix[i][0] = i
	}

	for j := 0; j <= tl; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= sl; i++ {
		for j := 1; j <= tl; j++ {
			cost := 0
			if source[i-1] != target[j-1] {
				cost = 1
			}
			matrix[i][j] = minInThree(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	distance := matrix[sl][tl]
	maxLen := float64(sl)
	if tl > sl {
		maxLen = float64(tl)
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
