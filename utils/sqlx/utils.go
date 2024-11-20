package sqlx

import (
	"github.com/go-xuan/quanx/types/stringx"
	"regexp"
	"strconv"
	"strings"
)

func parseValuesInSql(sql string) (string, *strings.Replacer) {
	values := regexp.MustCompile(`'[^']*'`).FindAllString(sql, -1)
	if len(values) > 0 {
		var all []string
		for i, value := range values {
			var key = ReplacePrefix + strconv.Itoa(i+1)
			sql = strings.Replace(sql, value, key, 1)
			all = append(all, key, value)
		}
		replacer := strings.NewReplacer(all...)
		return sql, replacer
	}
	return sql, nil
}

// 将所有关键字转为小写
func allKeywordsToLower(sql string) string {
	var oldnew []string
	var KEYWORDS = []string{
		SELECT, FROM, WHERE, SET, JOIN, GROUP, ORDER, HAVING, LIMIT, OFFSET,
		ASC, DESC, CASE, WHEN, THEN, END, INNER, OUTER, LEFT, RIGHT,
		DISTINCT, PARTITION, OVER, AS, AND, ON, OR, IN, NOT, LIKE, By,
	}
	for _, keyword := range KEYWORDS {
		switch keyword {
		case ASC, DESC:
			keyword = " " + keyword
		default:
			keyword = keyword + " "
		}
		oldnew = append(oldnew, strings.ToUpper(keyword), keyword)
	}
	sql = strings.NewReplacer(oldnew...).Replace(sql)
	return sql
}

// SplitButIgnoreInBracket 根据分隔符进行拆分但是忽略括号内的分隔符
func SplitButIgnoreInBracket(sql, sep string) ([]string, string) {
	var slice []string
	// l：总长度  k:sep长度  m:上个拆分点  n:括号个数
	var l, k, m, n = len(sql), len(sep), 0, 0
	for i := 0; i < l-k; i++ {
		if sql[i] == sep[0] && sql[i:i+k] == sep {
			if n == 0 {
				// 当前位置已将前面的括号对全部消完，才是有效的分割位
				slice = append(slice, sql[m:i])
				m = i + k // 将当前拆分点后移一个sep长度
			}
		} else if sql[i:i+1] == LeftBracket {
			n++ // 遍历到左括号则加一
		} else if sql[i:i+1] == RightBracket && n > 0 {
			n-- // 遍历到右括号则消掉
		}
	}
	return slice, sql[m:]
}

// ContainsKeywords 字符串是否包含sql关键字
func ContainsKeywords(sql string, keys ...string) (string, int) {
	var hit, index = "", -1
	for _, key := range keys {
		if i := KeywordIndex(sql, key); i >= 0 {
			if i < index {
				hit, index = key, i
			} else if index == -1 {
				hit, index = key, i
			}
		}
	}
	return hit, index
}

// LastIndexOfKeywords 字符串是否包含sql关键字
func LastIndexOfKeywords(sql string, keys ...string) int {
	for _, key := range keys {
		if i := KeywordIndex(sql, key, -1); i >= 0 {
			return i
		}
	}
	return -1
}

func KeywordIndex(sql, key string, positions ...int) int {
	// 获取所有命中字符下标
	if allIndices := stringx.AllIndex(sql, key); len(allIndices) >= 0 {
		sl, kl := len(sql), len(key)
		// 筛选出有效下标
		var validIndices []int
		for _, index := range allIndices {
			if index == 0 {
				if sql[index+kl:index+kl+1] == Blank {
					validIndices = append(validIndices, index)
				}
			} else if index == sl-kl {
				if sql[index-1:index] == Blank {
					validIndices = append(validIndices, index)
				}
			} else if sql[index-1:index] == Blank && sql[index+kl:index+kl+1] == Blank {
				validIndices = append(validIndices, index)
			}
		}
		if vl := len(validIndices); vl > 0 {
			if len(positions) > 0 {
				if position := positions[0]; position > 0 && position <= vl { // 正序
					return validIndices[position-1]
				} else if position < 0 && vl+position >= 0 { // 倒序
					return validIndices[vl+position]
				} else {
					return -1
				}
			} else {
				return validIndices[0]
			}
		}
	}
	return -1
}

// DB2GoType DB-Go类型映射
func DB2GoType(t string) string {
	switch t {
	case Char, Varchar, Varchar100, Text, Uuid:
		return String
	case Int, Int2, Int4, Tinyint, Smallint, Mediumint:
		return Int
	case Int8, Bigint:
		return Int64
	case Float, Float4, Numeric:
		return Float64
	case Timestamp, Timestampz, Datetime, Time, Date:
		return TimeTime
	case Bool:
		return Bool
	default:
		return String
	}
}
