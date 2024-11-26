package sqlx

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

type Parser interface {
	Beautify() string
}

type ParserBase struct {
	originSql string            // 原始sql
	tempSql   string            // 临时sql
	replacer  *strings.Replacer // 替换器
}

// sql准备：将所有参数值替换为占位符, 并转换所有关键字为小写
func (p *ParserBase) prepareSQL() {
	sql := p.tempSql
	// 提取sql中所有的参数值，避免参数值值影响后续sql解析
	var replacer *strings.Replacer
	if sql, replacer = parseValuesInSql(sql); replacer != nil {
		p.replacer = replacer
	}
	// 将sql中所有关键字转为小写
	p.tempSql = allKeywordsToLower(sql)
	log.Info("准备SQL完成")
}

func Parse(sql string) Parser {
	sql = strings.ReplaceAll(sql, NewLine, Blank)                // 移除换行
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, Blank) // 去除多余空格
	sql = strings.TrimSpace(sql)                                 // 去除空格
	sqlType := strings.ToLower(sql[:5])                          // 根据sql查询语句开头关键字判断sql类型
	switch sqlType {
	case SELECT:
		return parseSelectSQL(sql)
	case UPDATE:
		return parseUpdateSQL(sql)
	case DELETE:
		return parseDeleteSQL(sql)
	case INSERT:
		return parseInsertSQL(sql)
	case CREATE:
		return parseCreateSQL(sql)
	default:
		panic("")
	}
}
