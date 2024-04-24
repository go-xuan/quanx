package sqlx

import (
	"regexp"
	"strings"

	"github.com/go-xuan/quanx/types/anyx"
)

func SqlFormat(sql string) string {
	return Format(sql).ToString()
}

func Format(sql string, space ...int) *SQL {
	sql = strings.ReplaceAll(sql, NewLine, Blank)
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, Blank)
	sql = strings.TrimSpace(sql)
	return newSQL(sql, space...)
}

func newSQL(sql string, space ...int) *SQL {
	return initSQL(sql, space...).
		prepareSQL(). // sql预备
		extractLimit(). // 提取limit
		extractOrderBy(). // 提取order by
		extractFields(). // 提取查询字段
		extractFrom(). // 提取from主表
		extractJoins(). // 提取表关联
		extractWhere(). // 提取where
		extractGroupBy(). // 提取group By
		extractHaving() // 提取having
}

func initSQL(sql string, space ...int) *SQL {
	return &SQL{Sql: sql, Origin: sql, Space: anyx.Default(space, 0) + 6}
}

func SplitAndIgnoreBracket(sql, sep string) ([]string, string) {
	var slice []string
	// l：总长度  k:sep长度  m:上个拆分点  n:括号个数
	var l, k, m, n = len(sql), len(sep), 0, 0
	for i := 0; i < l-k; i++ {
		if sql[i] == sep[0] && sql[i:i+k] == sep {
			if n == 0 {
				slice = append(slice, sql[m:i])
				m = i + k
			}
		} else if sql[i:i+1] == LeftBracket {
			n++
		} else if sql[i:i+1] == RightBracket && n > 0 {
			n--
		}
	}
	return slice, sql[m:]
}
