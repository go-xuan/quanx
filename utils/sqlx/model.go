package sqlx

import (
	"fmt"
	"github.com/go-xuan/quanx/types/stringx"
	"regexp"
	"strconv"
	"strings"
)

type SQL struct {
	Origin   string
	Sql      string
	Replacer *strings.Replacer
	Distinct bool
	Space    int
	Fields   []*Field
	From     *Table
	Joins    []*Join
	Where    []*Condition
	GroupBy  []string
	Having   []*Condition
	OrderBy  []string
	Limit    string
}

// sql输出
func (s *SQL) ToString() string {
	var sql = strings.Builder{}
	sql.WriteString(s.SelectSql())
	sql.WriteString(s.FromSql())
	sql.WriteString(s.ConditionSql(WHERE))
	sql.WriteString(s.GroupOrderSql(GROUP))
	sql.WriteString(s.ConditionSql(HAVING))
	sql.WriteString(s.GroupOrderSql(ORDER))
	sql.WriteString(s.LimitSql())
	if s.Replacer != nil {
		return s.Replacer.Replace(sql.String())
	} else {
		return sql.String()
	}
}

// sql预备：将所有参数值替换为占位符, 并转换所有关键字为小写
func (s *SQL) prepareSQL() *SQL {
	sql := s.Sql
	values := regexp.MustCompile(`'[^']*'`).FindAllString(sql, -1)
	if len(values) > 0 {
		var all []string
		for i, value := range values {
			var key = ReplacePrefix + strconv.Itoa(i+1)
			sql = strings.Replace(sql, value, key, 1)
			all = append(all, key, value)
		}
		s.Replacer = strings.NewReplacer(all...)
	}
	var toLower []string
	for _, keyword := range KEYWORDS {
		toLower = append(toLower, strings.ToUpper(keyword), keyword)
	}
	s.Sql = strings.NewReplacer(toLower...).Replace(sql)
	return s
}

func (s *SQL) extractFields() *SQL {
	sql := s.Sql
	if x, y := stringx.Between(sql, SELECT, FROM); x > 0 {
		fieldsSql := sql[x+1 : y]
		if y-x > 9 && fieldsSql[:8] == DISTINCT {
			s.Distinct = true
			fieldsSql = fieldsSql[9:]
		}
		// 判断是否有字段包含括号（子查询/函数等内部可能会包含","逗号，从而影响字段拆分）
		var sqlList, lastSql = SplitAndIgnoreBracket(fieldsSql, Comma)
		sqlList = append(sqlList, lastSql)
		var fields []*Field
		for _, fieldSql := range sqlList {
			var name, alias string
			fieldSql = strings.TrimSpace(fieldSql)
			if i := stringx.Index(fieldSql, " as "); i >= 0 {
				name, alias = fieldSql[:i], fieldSql[i+1:]
			} else if fieldSql[len(fieldSql)-1:] == RightBracket {
				name = fieldSql
			} else if i = strings.LastIndex(fieldSql, Blank); i > 0 {
				name, alias = fieldSql[:i], fieldSql[i+1:]
			} else {
				name = fieldSql
			}
			if stringx.Index(name, ReplacePrefix) >= 0 {
				name = s.Replacer.Replace(name)
			}
			fields = append(fields, &Field{Name: name, Alias: alias})
		}
		s.Fields = fields
		s.Sql = sql[y:]
	}
	return s
}

func (s *SQL) extractFrom() *SQL {
	if f := stringx.Index(s.Sql, FROM); f >= 0 {
		sql := s.Sql[f+4:]
		var from = &Table{}
		var x, y, z int
		if sql[1:2] == LeftBracket {
			x, y = stringx.Between(sql, LeftBracket, RightBracket)
			from.Select = newSQL(sql[x:y], s.Space+2)
		} else {
			x = stringx.Index(sql, Blank, 1)
			y = stringx.Index(sql, Blank, 2)
			from.Name = sql[x+1 : y]
		}
		sql = sql[y+1:]
		var aliasSql string
		if _, i := stringx.Contains(sql, LEFT, RIGHT, INNER, JOIN, WHERE, GROUP, ORDER, LIMIT); i >= 0 {
			aliasSql, sql = sql[:i], sql[i:]
		} else {
			aliasSql, sql = sql, Empty
		}
		aliasSql = strings.TrimSpace(aliasSql)
		if z = stringx.Index(aliasSql, AS); z >= 0 {
			z = z + 3
		} else {
			z = stringx.Index(aliasSql, Blank) + 1
		}
		if z > 0 {
			from.Alias = aliasSql[z:]
		} else {
			from.Alias = aliasSql
		}
		s.From = from
		s.Sql = sql
	}
	return s
}

func (s *SQL) extractJoins() *SQL {
	sql := s.Sql
	var sqlList []string
	sqlList, sql = SplitAndIgnoreBracket(sql, JOIN)
	var lastJoin string
	if x, y := stringx.Between(sql, LeftBracket, RightBracket); x == 2 {
		lastJoin, sql = sql[:y], sql[y:]
	}
	if _, i := stringx.Contains(sql, WHERE, GROUP, ORDER, LIMIT); i >= 0 {
		lastJoin, sql = lastJoin+sql[:i], sql[i:]
	} else {
		lastJoin, sql = lastJoin+sql, Empty
	}
	sqlList = append(sqlList, lastJoin)
	if len(sqlList) > 0 {
		var joinType string
		var joins []*Join
		for i, joinSql := range sqlList {
			if i == 0 {
				joinType = strings.TrimSpace(joinSql)
			} else {
				var join = &Join{}
				var space = s.Space - 1
				if joinType == Empty {
					space = space - 5
				}
				join.Type = joinType
				if a := stringx.Index(joinSql, LEFT, -1); a >= 0 {
					joinType = LEFT
					joinSql = joinSql[:a-1]
				}
				if a := stringx.Index(joinSql, RIGHT, -1); a >= 0 {
					joinType = RIGHT
					joinSql = joinSql[:a-1]
				}
				if a := stringx.Index(joinSql, INNER, -1); a >= 0 {
					joinType = INNER
					joinSql = joinSql[:a-1]
				}
				if a := stringx.Index(joinSql, ON, -1); a >= 0 {
					join.On, joinSql = joinSql[a+3:], joinSql[:a-1]
				}

				joinSql = FROM + joinSql
				join.Table = initSQL(joinSql, space).extractFrom().From
				joins = append(joins, join)
			}
		}
		s.Joins = joins
		s.Sql = sql
	}
	return s
}

func (s *SQL) extractWhere() *SQL {
	if f := stringx.Index(s.Sql, WHERE); f >= 0 {
		sql := s.Sql[f+5:]
		var whereSql string
		if _, i := stringx.Contains(sql, GROUP, ORDER, LIMIT); i >= 0 {
			whereSql, sql = sql[:i], sql[i:]
		} else {
			whereSql, sql = sql, Empty
		}
		var sqlList, lastSql = SplitAndIgnoreBracket(whereSql, AND)
		sqlList = append(sqlList, lastSql)
		if len(sqlList) > 0 {
			var conditions []*Condition
			for _, conditionSql := range sqlList {
				conditions = append(conditions, &Condition{Content: strings.TrimSpace(conditionSql)})
			}
			s.Where = conditions
		}
		s.Sql = sql
	}
	return s
}

func (s *SQL) extractGroupBy() *SQL {
	sql := s.Sql
	if i := stringx.Index(sql, GROUP); i >= 0 {
		var groupBySql string
		if _, j := stringx.Contains(sql, HAVING, ORDER, LIMIT); j >= 0 {
			groupBySql, sql = sql[i+9:j], sql[j:]
		} else {
			groupBySql, sql = sql, Empty
		}
		s.GroupBy = strings.Split(groupBySql, Comma)
	}
	s.Sql = sql
	return s
}

func (s *SQL) extractHaving() *SQL {
	if len(s.GroupBy) > 0 {
		if f := stringx.Index(s.Sql, HAVING); f >= 0 {
			sql := s.Sql[f+6:]
			var havingSql string
			if _, i := stringx.Contains(sql, ORDER, LIMIT); i >= 0 {
				havingSql, sql = sql[:i], sql[i:]
			} else {
				havingSql, sql = sql, Empty
			}
			var sqlList, lastSql = SplitAndIgnoreBracket(havingSql, AND)
			sqlList = append(sqlList, lastSql)
			if len(sqlList) > 0 {
				var conditions []*Condition
				for _, conditionSql := range sqlList {
					conditions = append(conditions, &Condition{Content: strings.TrimSpace(conditionSql)})
				}
				s.Having = conditions
			}
			s.Sql = sql
		}
	}
	return s
}

func (s *SQL) extractOrderBy() *SQL {
	sql := s.Sql
	if i := stringx.Index(sql, ORDER, -1); i > 0 {
		var orderBySql string
		if j := stringx.Index(sql, RightBracket, -1); i > j {
			orderBySql, sql = sql[i+9:], sql[:i-1]
		}
		if orderBySql != Empty {
			s.OrderBy = strings.Split(orderBySql, Comma)
		}
	}
	s.Sql = sql
	return s
}

func (s *SQL) extractLimit() *SQL {
	sql := s.Sql
	i := stringx.Index(sql, LIMIT, -1)
	if j := stringx.Index(sql, RightBracket, -1); i > 0 && i > j {
		s.Limit, sql = sql[i+6:], sql[:i]
	}
	s.Sql = sql
	return s
}

// 对齐
func (s *SQL) Align(in ...string) string {
	if len(in) == 0 {
		return strings.Repeat(Blank, s.Space)
	} else if key := in[0]; len(key) <= s.Space {
		return fmt.Sprintf("%"+strconv.Itoa(s.Space)+"s", key)
	} else {
		return key
	}
}

// 新行空格
func (s *SQL) NewLineSpace(n int) string {
	return strings.Repeat(Blank, s.Space+n)
}

func (s *SQL) SelectSql() string {
	var sql = strings.Builder{}
	var space = 1
	sql.WriteString(SELECT)
	sql.WriteString(Blank)
	if s.Distinct {
		sql.WriteString(DISTINCT)
		sql.WriteString(Blank)
		space += 9
	}
	var fieldLength, aliasNum int
	for _, field := range s.Fields {
		y := len(field.Name)
		if fieldLength < y {
			fieldLength = y
		}
		if field.Alias != Empty {
			aliasNum++
		}
	}
	for i, field := range s.Fields {
		if i > 0 {
			sql.WriteString(Comma)
			if aliasNum > 3 {
				sql.WriteString(NewLine)
				sql.WriteString(s.NewLineSpace(space))
			} else {
				sql.WriteString(Blank)
			}
		}
		if field.Alias != Empty {
			sql.WriteString(fmt.Sprintf("%"+strconv.Itoa(-fieldLength-1)+"s", field.Name))
			sql.WriteString(field.Alias)
		} else {
			sql.WriteString(field.Name)
		}
	}
	return sql.String()
}

func (s *SQL) FromSql() string {
	sql := strings.Builder{}
	sql.WriteString(NewLine)
	sql.WriteString(s.Align(FROM))
	sql.WriteString(Blank)
	sql.WriteString(s.From.toString())
	for _, join := range s.Joins {
		sql.WriteString(NewLine)
		if join.Type != Empty {
			sql.WriteString(s.Align(join.Type))
			sql.WriteString(Blank)
			sql.WriteString(JOIN)
		} else {
			sql.WriteString(s.Align(JOIN))
		}
		sql.WriteString(Blank)
		sql.WriteString(join.Table.toString())
		sql.WriteString(NewLine)
		sql.WriteString(s.Align(ON))
		sql.WriteString(Blank)
		sql.WriteString(join.On)
	}
	return sql.String()
}

func (s *SQL) ConditionSql(in string) string {
	var conditions []*Condition
	switch in {
	case WHERE:
		conditions = s.Where
	case HAVING:
		conditions = s.Having
	default:
		in = WHERE
		conditions = s.Where
	}
	sql := strings.Builder{}
	if len(conditions) > 0 {
		sql.WriteString(NewLine)
		sql.WriteString(s.Align(in))
		sql.WriteString(Blank)
		for i, cond := range conditions {
			if i > 0 {
				sql.WriteString(NewLine)
				if cond.Type == Empty {
					sql.WriteString(s.Align(AND))
					sql.WriteString(Blank)
				} else {
					sql.WriteString(s.Align(cond.Type))
					sql.WriteString(Blank)
				}
			}
			sql.WriteString(cond.Content)
		}
	}
	return sql.String()
}

func (s *SQL) GroupOrderSql(in string) string {
	var values []string
	switch in {
	case GROUP:
		values = s.GroupBy
	case ORDER:
		values = s.OrderBy
	}
	sql := strings.Builder{}
	if len(values) > 0 {
		sql.WriteString(NewLine)
		sql.WriteString(s.Align(in))
		sql.WriteString(Blank)
		sql.WriteString(By)
		sql.WriteString(Blank)
		for i, value := range values {
			value = strings.TrimSpace(value)
			if i > 0 {
				sql.WriteString(Comma)
				sql.WriteString(NewLine)
				sql.WriteString(s.NewLineSpace(4))
			}
			sql.WriteString(value)
		}
	}
	return sql.String()
}

func (s *SQL) LimitSql() string {
	sql := strings.Builder{}
	if s.Limit != Empty {
		sql.WriteString(NewLine)
		sql.WriteString(s.Align(LIMIT))
		sql.WriteString(Blank)
		sql.WriteString(s.Limit)
	}
	return sql.String()
}

type Field struct {
	Name  string
	Alias string
}

type Table struct {
	Select *SQL
	Name   string
	Alias  string
}

func (f *Table) toString() string {
	if f.Select != nil {
		return LeftBracket + f.Select.ToString() + RightBracket + Blank + f.Alias
	} else {
		return f.Name + Blank + f.Alias
	}
}

type Join struct {
	*Table
	Type string // left/right/inner
	On   string
}

type Condition struct {
	Type    string // or/and
	Content string
}
