package codex

import (
	"strings"

	"github.com/go-xuan/quanx/utilx/sqlx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

// 构建go结构体代码
func BuildGoStruct(table string, fieldList []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	table = stringx.UpperCamelCase(table)
	sb := strings.Builder{}
	sb.WriteString("\ntype ")
	sb.WriteString(table)
	sb.WriteString(" struct {")
	for _, field := range fieldList {
		sb.WriteString("\n\t")
		sb.WriteString(stringx.UpperCamelCase(field.Name))
		sb.WriteString(" ")
		sb.WriteString(sqlx.Pg2GoTypeMap[field.Type])
		sb.WriteString(" `json:\"")
		sb.WriteString(stringx.LowerCamelCase(field.Name))
		sb.WriteString("\"` // ")
		sb.WriteString(field.Comment)
	}
	sb.WriteString("\n}\n\n")
	return sb.String()
}

// 构建go结构体代码
func BuildGormStruct(table string, fieldList []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	table = stringx.UpperCamelCase(table)
	sb := strings.Builder{}
	sb.WriteString("\ntype ")
	sb.WriteString(table)
	sb.WriteString(" struct {")
	for _, field := range fieldList {
		sb.WriteString("\n\t")
		sb.WriteString(stringx.UpperCamelCase(field.Name))
		sb.WriteString(" ")
		sb.WriteString(sqlx.Pg2GoTypeMap[field.Type])
		sb.WriteString(" `json:\"")
		sb.WriteString(stringx.LowerCamelCase(field.Name))
		sb.WriteString("\" gorm:\"type:")
		sb.WriteString(sqlx.Pg2GormTypeMap[field.Type])
		sb.WriteString("; comment:")
		sb.WriteString(field.Comment)
		sb.WriteString(";\"`")
	}
	sb.WriteString("\n}\n\n")
	return sb.String()
}
