package sqlx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-xuan/quanx/types/intx"
	"github.com/go-xuan/quanx/types/stringx"
	log "github.com/sirupsen/logrus"
)

// 解析sql字符串
func parseSelectSQL(sql string, indent ...int) *SelectSqlParser {
	parser := newSelectSqlParser(sql, indent...) // 初始化
	parser.prepareSQL()                          // sql准备
	parser.extractLimit()                        // 提取limit
	parser.extractOrderBy()                      // 提取order by
	parser.extractFields()                       // 提取查询字段
	parser.extractTable()                        // 提取主表
	parser.extractJoins()                        // 提取关联子表
	parser.extractWhere()                        // 提取where
	parser.extractGroupBy()                      // 提取group By
	parser.extractHaving()                       // 提取having
	return parser
}

// 初始化查询SQL对象实例
func newSelectSqlParser(sql string, indent ...int) *SelectSqlParser {
	return &SelectSqlParser{
		ParserBase: ParserBase{
			originSql: sql,
			tempSql:   sql,
		},
		indent: intx.Default(0, indent...) + 6}
}

type SelectSqlParser struct {
	ParserBase
	Table    *Table       // 查询主表
	Fields   []*Field     // 查询字段
	Joins    []*Join      // 关联子表
	Where    []*Condition // 查询条件
	GroupBy  []string     // 分组条件
	Having   []*Condition // 分组筛选条件
	OrderBy  []string     // 排序条件
	Limit    string       // 限数条件
	Distinct bool         // 是否distinct
	indent   int          // 缩进量
}

// Beautify SQL美化输出
func (p *SelectSqlParser) Beautify() string {
	var sql = strings.Builder{}
	sql.WriteString(p.buildSelectSql())
	sql.WriteString(p.buildFromSql())
	sql.WriteString(p.buildConditionSql(WHERE))
	sql.WriteString(p.buildGroupOrderSql(GROUP))
	sql.WriteString(p.buildConditionSql(HAVING))
	sql.WriteString(p.buildGroupOrderSql(ORDER))
	sql.WriteString(p.buildLimitSql())
	if p.replacer != nil {
		return p.replacer.Replace(sql.String())
	} else {
		return sql.String()
	}
}

// 提取查询字段
func (p *SelectSqlParser) extractFields() *SelectSqlParser {
	sql := p.tempSql
	if x, y := stringx.Between(sql, SELECT, FROM); x > 0 {
		fieldsSql := sql[x+1 : y]
		if y-x > 9 && fieldsSql[:8] == DISTINCT {
			p.Distinct = true
			fieldsSql = fieldsSql[9:]
		}
		// 判断是否有字段包含括号（子查询或者函数等内部可能会包含","逗号，从而影响字段拆分）
		var sqlList, lastSql = SplitButIgnoreInBracket(fieldsSql, Comma)
		sqlList = append(sqlList, lastSql)
		var fields []*Field
		for _, fieldSql := range sqlList {
			var name, alias string
			fieldSql = strings.TrimSpace(fieldSql)
			if i := KeywordIndex(fieldSql, AS); i >= 0 {
				name, alias = fieldSql[:i], fieldSql[i:]
			} else if fieldSql[len(fieldSql)-1:] == RightBracket {
				name = fieldSql
			} else if i = strings.LastIndex(fieldSql, Blank); i > 0 {
				name, alias = fieldSql[:i], fieldSql[i+1:]
			} else {
				name = fieldSql
			}
			if stringx.Index(name, ReplacePrefix) >= 0 {
				name = p.replacer.Replace(name)
			}
			fields = append(fields, &Field{Name: name, Alias: alias})
		}
		p.Fields = fields
		p.tempSql = sql[y:]
	}
	log.Info("提取查询字段完成")
	return p
}

// 提取查询主表
func (p *SelectSqlParser) extractTable() *SelectSqlParser {
	// 判断第一个from所在位置
	if f := KeywordIndex(p.tempSql, FROM); f >= 0 {
		sql := p.tempSql[f+4:] // 截取掉from，但是保留表名前面的空格
		var table = &Table{}
		var x, y, z int
		if sql[1:2] == LeftBracket { // 如果from后面跟括号，表示是子查询
			x, y = stringx.Between(sql, LeftBracket, RightBracket) // 截取最近的子查询，根据一个对括号内容进行截取
			table.Select = parseSelectSQL(sql[x:y], p.indent+2)
		} else { // from后面直接跟表名
			x = stringx.Index(sql, Blank, 1) // 表名前空格下标
			y = stringx.Index(sql, Blank, 2) // 表名后空格下标
			if y >= 0 {
				table.Name, sql = sql[x+1:y], sql[y+1:] // 如果表名后面还跟有其他sql，则截取掉表名，继续处理
			} else {
				table.Name, sql = sql[x+1:], ""
			}
		}
		if sql != "" {
			var alias string
			if _, i := ContainsKeywords(sql, LEFT, RIGHT, INNER, OUTER, JOIN, WHERE, GroupBy, OrderBy, LIMIT); i >= 0 {
				// 判断是否是复杂查询
				alias, sql = sql[:i], sql[i:]
			} else { // 简单查询
				alias, sql = sql, Empty
			}
			alias = strings.TrimSpace(alias)
			if z = KeywordIndex(alias, AS); z >= 0 {
				z = z + 3
			} else {
				z = stringx.Index(alias, Blank) + 1
			}
			if z > 0 {
				table.Alias = alias[z:]
			} else {
				table.Alias = alias
			}
		}
		p.Table = table
		p.tempSql = sql
	}
	log.Info("提取查询主表完成")
	return p
}

// 提取关联子表
func (p *SelectSqlParser) extractJoins() *SelectSqlParser {
	sql := p.tempSql
	var sqlList []string
	sqlList, sql = SplitButIgnoreInBracket(sql, JOIN)
	var lastJoin string
	if x, y := stringx.Between(sql, LeftBracket, RightBracket); x == 2 {
		lastJoin, sql = sql[:y], sql[y:]
	}
	if _, i := ContainsKeywords(sql, WHERE, GroupBy, OrderBy, LIMIT); i >= 0 {
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
				var space = p.indent - 1
				if joinType == Empty {
					space = space - 5
				}
				join.Type = joinType

				if a := LastIndexOfKeywords(joinSql, LEFT, RIGHT, INNER, OUTER); a >= 0 {
					joinType = strings.TrimSpace(joinSql[a:])
					joinSql = joinSql[:a-1]
				}
				if a := KeywordIndex(joinSql, ON, -1); a >= 0 {
					join.On, joinSql = joinSql[a+3:], joinSql[:a-1]
				}

				join.Table = newSelectSqlParser(FROM+joinSql, space).extractTable().Table
				joins = append(joins, join)
			}
		}
		p.Joins = joins
		p.tempSql = sql
	}
	log.Info("提取关联子表完成")
	return p
}

// 提取查询条件
func (p *SelectSqlParser) extractWhere() *SelectSqlParser {
	if f := KeywordIndex(p.tempSql, WHERE); f >= 0 {
		sql := p.tempSql[f+5:]
		var whereSql string
		if _, i := ContainsKeywords(sql, GroupBy, OrderBy, LIMIT); i >= 0 {
			whereSql, sql = sql[:i], sql[i:]
		} else {
			whereSql, sql = sql, Empty
		}
		var sqlList, lastSql = SplitButIgnoreInBracket(whereSql, AND)
		sqlList = append(sqlList, lastSql)
		if len(sqlList) > 0 {
			var conditions []*Condition
			for _, conditionSql := range sqlList {
				conditions = append(conditions, &Condition{Content: strings.TrimSpace(conditionSql)})
			}
			p.Where = conditions
		}
		p.tempSql = sql
	}
	log.Info("提取查询条件完成")
	return p
}

// 提取group By
func (p *SelectSqlParser) extractGroupBy() *SelectSqlParser {
	sql := p.tempSql
	if i := KeywordIndex(sql, GroupBy); i >= 0 {
		var groupBySql string
		if _, j := ContainsKeywords(sql, HAVING, OrderBy, LIMIT); j >= 0 {
			groupBySql, sql = sql[i+9:j], sql[j:]
		} else {
			groupBySql, sql = sql, Empty
		}
		p.GroupBy = strings.Split(groupBySql, Comma)
	}
	p.tempSql = sql
	log.Info("提取分组条件完成")
	return p
}

// 提取having
func (p *SelectSqlParser) extractHaving() *SelectSqlParser {
	if len(p.GroupBy) > 0 {
		if f := KeywordIndex(p.tempSql, HAVING); f >= 0 {
			sql := p.tempSql[f+6:]
			var havingSql string
			if _, i := ContainsKeywords(sql, OrderBy, LIMIT); i >= 0 {
				havingSql, sql = sql[:i], sql[i:]
			} else {
				havingSql, sql = sql, Empty
			}
			var sqlList, lastSql = SplitButIgnoreInBracket(havingSql, AND)
			sqlList = append(sqlList, lastSql)
			if len(sqlList) > 0 {
				var conditions []*Condition
				for _, conditionSql := range sqlList {
					conditions = append(conditions, &Condition{Content: strings.TrimSpace(conditionSql)})
				}
				p.Having = conditions
			}
			p.tempSql = sql
		}
	}
	log.Info("提取分组过滤条件完成")
	return p
}

// 提取order by
func (p *SelectSqlParser) extractOrderBy() *SelectSqlParser {
	sql := p.tempSql
	if i := KeywordIndex(sql, OrderBy, -1); i > 0 {
		var orderBySql string
		if j := KeywordIndex(sql, RightBracket, -1); i > j {
			orderBySql, sql = sql[i+9:], sql[:i-1]
		}
		if orderBySql != Empty {
			p.OrderBy = strings.Split(orderBySql, Comma)
		}
	}
	p.tempSql = sql
	log.Info("提取排序条件完成")
	return p
}

// 提取limit
func (p *SelectSqlParser) extractLimit() *SelectSqlParser {
	sql := p.tempSql
	i := KeywordIndex(sql, LIMIT, -1)
	if j := KeywordIndex(sql, RightBracket, -1); i > 0 && i > j {
		p.Limit, sql = sql[i+6:], sql[:i]
	}
	p.tempSql = sql
	log.Info("提取限制数量完成")
	return p
}

// 以当前缩进量对齐
func (p *SelectSqlParser) align(sql ...string) string {
	if len(sql) == 0 {
		return strings.Repeat(Blank, p.indent)
	} else if str := sql[0]; len(str) <= p.indent {
		return strings.Repeat(Blank, p.indent-len(str)) + str
	} else if str1, _ := stringx.Cut(str, Blank); len(str1) <= p.indent {
		return strings.Repeat(Blank, p.indent-len(str1)) + str
	} else {
		return str
	}
}

// 新行以当前缩进量添加空格
func (p *SelectSqlParser) newLineSpace(n int) string {
	return strings.Repeat(Blank, p.indent+n)
}

// 构建查询字段sql
func (p *SelectSqlParser) buildSelectSql() string {
	var sql = strings.Builder{}
	var space = 1
	sql.WriteString(SELECT)
	sql.WriteString(Blank)
	if p.Distinct {
		sql.WriteString(DISTINCT)
		sql.WriteString(Blank)
		space += 9
	}
	var fieldLength, aliasNum int
	for _, field := range p.Fields {
		y := len(field.Name)
		if fieldLength < y {
			fieldLength = y
		}
		if field.Alias != Empty {
			aliasNum++
		}
	}
	for i, field := range p.Fields {
		if i > 0 {
			sql.WriteString(Comma)
			if aliasNum > 0 {
				sql.WriteString(NewLine)
				sql.WriteString(p.newLineSpace(space))
			} else {
				sql.WriteString(Blank)
			}
		}
		if field.Alias != Empty {
			if aliasNum > 0 {
				sql.WriteString(fmt.Sprintf("%"+strconv.Itoa(-fieldLength-1)+"s", field.Name))
			} else {
				sql.WriteString(field.Name)
				sql.WriteString(Blank)
			}
			sql.WriteString(field.Alias)
		} else {
			sql.WriteString(field.Name)
		}
	}
	return sql.String()
}

func (p *SelectSqlParser) buildFromSql() string {
	sql := strings.Builder{}
	sql.WriteString(NewLine)
	sql.WriteString(p.align(FROM))
	sql.WriteString(Blank)
	sql.WriteString(p.Table.WithAliasSQL())
	for _, join := range p.Joins {
		sql.WriteString(NewLine)
		if join.Type != Empty {
			sql.WriteString(p.align(join.Type))
			sql.WriteString(Blank)
			sql.WriteString(JOIN)
		} else {
			sql.WriteString(p.align(JOIN))
		}
		sql.WriteString(Blank)
		sql.WriteString(join.Table.WithAliasSQL())
		sql.WriteString(NewLine)
		sql.WriteString(p.align(ON))
		sql.WriteString(Blank)
		sql.WriteString(join.On)
	}
	return sql.String()
}

func (p *SelectSqlParser) buildConditionSql(in string) string {
	var conditions []*Condition
	switch in {
	case WHERE:
		conditions = p.Where
	case HAVING:
		conditions = p.Having
	default:
		in = WHERE
		conditions = p.Where
	}
	sql := strings.Builder{}
	if len(conditions) > 0 {
		sql.WriteString(NewLine)
		sql.WriteString(p.align(in))
		sql.WriteString(Blank)
		for i, cond := range conditions {
			if i > 0 {
				sql.WriteString(NewLine)
				if cond.Type == Empty {
					sql.WriteString(p.align(AND))
					sql.WriteString(Blank)
				} else {
					sql.WriteString(p.align(cond.Type))
					sql.WriteString(Blank)
				}
			}
			sql.WriteString(cond.Content)
		}
	}
	return sql.String()
}

func (p *SelectSqlParser) buildGroupOrderSql(in string) string {
	var values []string
	switch in {
	case GROUP:
		values = p.GroupBy
	case ORDER:
		values = p.OrderBy
	}
	sql := strings.Builder{}
	if len(values) > 0 {
		sql.WriteString(NewLine)
		sql.WriteString(p.align(in))
		sql.WriteString(Blank)
		sql.WriteString(By)
		sql.WriteString(Blank)
		for i, value := range values {
			value = strings.TrimSpace(value)
			if i > 0 {
				sql.WriteString(Comma)
				sql.WriteString(NewLine)
				sql.WriteString(p.newLineSpace(4))
			}
			sql.WriteString(value)
		}
	}
	return sql.String()
}

func (p *SelectSqlParser) buildLimitSql() string {
	sql := strings.Builder{}
	if p.Limit != Empty {
		sql.WriteString(NewLine)
		sql.WriteString(p.align(LIMIT))
		sql.WriteString(Blank)
		sql.WriteString(p.Limit)
	}
	return sql.String()
}
