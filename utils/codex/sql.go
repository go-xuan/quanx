package codex

import (
	"github.com/go-xuan/quanx/public/constx"
	"strings"

	"github.com/go-xuan/quanx/utils/sqlx"
	"github.com/go-xuan/quanx/utils/stringx"
)

// 构建CK远程表建表语句
func BuildCkCreateSql(table, engine string, fieldList FieldList) string {
	sb := strings.Builder{}
	sb.WriteString("drop")
	sb.WriteString(" table ")
	sb.WriteString(table)
	sb.WriteString(";")
	sb.WriteString(constx.NextLine)
	sb.WriteString("create")
	sb.WriteString(" table ")
	sb.WriteString(table)
	sb.WriteString(constx.NextLine)
	sb.WriteString("(")
	sb.WriteString(constx.NextLine)
	for i, field := range fieldList {
		if i > 0 {
			sb.WriteString(",")
			sb.WriteString(constx.NextLine)
		}
		sb.WriteString("   `")
		sb.WriteString(field.Name)
		sb.WriteString("` ")
		sb.WriteString(sqlx.Pg2CkTypeMap[field.Type])
		if field.Default != "" {
			sb.WriteString(" default ")
			sb.WriteString(field.Default)
		}
		sb.WriteString(" comment '")
		sb.WriteString(field.Comment)
		sb.WriteString("'")
	}
	sb.WriteString(constx.NextLine)
	sb.WriteString(")")
	sb.WriteString(constx.NextLine)
	sb.WriteString(engine)
	return sb.String()
}

// 构建新增sql语句
func BuildInsertSql(table string, fieldList FieldList) string {
	sb := strings.Builder{}
	iv := strings.Builder{}
	sb.WriteString("insert")
	sb.WriteString(" into ")
	sb.WriteString(table)
	sb.WriteString(constx.NextLine)
	sb.WriteString("  (")
	for i, field := range fieldList {
		if i > 0 {
			sb.WriteString(constx.NextLine)
			sb.WriteString("   ")
			iv.WriteString(constx.NextLine)
			iv.WriteString("   ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(",")
		iv.WriteString("#{item.")
		iv.WriteString(stringx.LowerCamelCase(field.Name))
		iv.WriteString("},")
	}
	sb.WriteString(",)")
	sb.WriteString(constx.NextLine)
	sb.WriteString("values")
	sb.WriteString(constx.NextLine)
	sb.WriteString("  (")
	sb.WriteString(iv.String())
	sb.WriteString(",)")
	sb.WriteString(constx.NextLine)
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 构建更新sql语句
func BuildUpdateSql(table string, fieldList FieldList) string {
	sb := strings.Builder{}
	sb.WriteString("update ")
	sb.WriteString(table)
	sb.WriteString(constx.NextLine)
	sb.WriteString("   set ")
	for i, field := range fieldList {
		if i > 0 {
			sb.WriteString(constx.NextLine)
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(" = #{item.")
		sb.WriteString(stringx.LowerCamelCase(field.Name))
		sb.WriteString("},")
	}
	sb.WriteString(",")
	sb.WriteString(constx.NextLine)
	sb.WriteString(" where id = #{item.id} ")
	sb.WriteString(constx.NextLine)
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 原始select
func BuildSelectSql(table string, fieldList FieldList) string {
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range fieldList {
		if i > 0 {
			sb.WriteString(constx.NextLine)
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(",")
	}
	sb.WriteString(",")
	sb.WriteString(constx.NextLine)
	sb.WriteString("  from ")
	sb.WriteString(table)
	sb.WriteString(constx.NextLine)
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 带字段别名的select
func BuildSelectSqlAlias(table string, columns FieldList) string {
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range columns {
		if i > 0 {
			sb.WriteString(constx.NextLine)
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(" as ")
		sb.WriteString(stringx.LowerCamelCase(field.Name))
		sb.WriteString(",")
	}
	sb.WriteString(",")
	sb.WriteString(constx.NextLine)
	sb.WriteString("  from ")
	sb.WriteString(table)
	sb.WriteString(constx.NextLine)
	return strings.ReplaceAll(sb.String(), ",,", "")
}
