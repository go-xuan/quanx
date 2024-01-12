package codex

import (
	"strings"

	"github.com/go-xuan/quanx/utilx/sqlx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

// 构建go结构体代码
func BuildGoStruct(table string, fields []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	table = stringx.ToUpperCamel(table)
	sb := strings.Builder{}
	sb.WriteString("\ntype ")
	sb.WriteString(table)
	sb.WriteString(" struct {")
	for _, field := range fields {
		sb.WriteString("\n\t")
		sb.WriteString(stringx.ToUpperCamel(field.Name))
		sb.WriteString(" ")
		sb.WriteString(sqlx.DB2Go()[field.Type])
		sb.WriteString(" `json:\"")
		sb.WriteString(stringx.ToLowerCamel(field.Name))
		sb.WriteString("\"` // ")
		sb.WriteString(field.Comment)
	}
	sb.WriteString("\n}\n\n")
	return sb.String()
}

// 构建go结构体代码
func BuildGormStruct(table string, fields []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	table = stringx.ToUpperCamel(table)
	sb := strings.Builder{}
	sb.WriteString("\ntype ")
	sb.WriteString(table)
	sb.WriteString(" struct {")
	for _, field := range fields {
		sb.WriteString("\n\t")
		sb.WriteString(stringx.ToUpperCamel(field.Name))
		sb.WriteString(" ")
		sb.WriteString(sqlx.DB2Go()[field.Type])
		sb.WriteString(" `json:\"")
		sb.WriteString(stringx.ToLowerCamel(field.Name))
		sb.WriteString("\" gorm:\"type:")
		sb.WriteString(sqlx.DB2Gorm()[field.Type])
		sb.WriteString("; comment:")
		sb.WriteString(field.Comment)
		sb.WriteString(";\"`")
	}
	sb.WriteString("\n}\n\n")
	return sb.String()
}
