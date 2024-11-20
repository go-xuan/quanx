package sqlx

import (
	"github.com/go-xuan/quanx/types/stringx"
	log "github.com/sirupsen/logrus"
	"strings"
)

func parseUpdateSQL(sql string) *UpdateSqlParser {
	parser := &UpdateSqlParser{
		ParserBase: ParserBase{
			originSql: sql,
			tempSql:   sql,
		},
	}
	parser.prepareSQL() // sql准备
	return parser
}

type UpdateSqlParser struct {
	ParserBase
	Table  *Table       // 更新表
	Fields []*Field     // 更新字段
	Where  []*Condition // 查询条件
}

func (p *UpdateSqlParser) Beautify() string {
	return ""
}

// 提取更新字段
func (p *UpdateSqlParser) extractFields() *UpdateSqlParser {
	sql := p.tempSql

	// 根据where关键字进行拆分
	if index := stringx.Index(sql, WHERE); index > 0 {
		p.tempSql = sql[index:]
		sql = sql[:index]
	}

	// 截取where关键字前面的sql片段
	if index := stringx.Index(sql, SET); index > 0 {
		sql = sql[index:]
		var fieldList, lastField = SplitButIgnoreInBracket(sql, Comma)
		fieldList = append(fieldList, lastField)
		var fields []*Field
		for _, field := range fieldList {
			var name, value string
			field = strings.TrimSpace(field)
			if i := stringx.Index(field, Equals); i >= 0 {
				name, value = field[:i], field[i+1:]
			}
			if stringx.Index(name, ReplacePrefix) >= 0 {
				name = p.replacer.Replace(name)
			}
			fields = append(fields, &Field{Name: name, Alias: value})
		}
		p.Fields = fields
	}
	log.Info("提取更新字段完成")
	return p
}
