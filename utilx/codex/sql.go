package codex

import (
	"fmt"
	"strings"

	"github.com/go-xuan/quanx/utilx/stringx"
)

// 构建新增sql语句
func BuildInsertSql(table string, fieldList []*Field) string {
	sb := strings.Builder{}
	iv := strings.Builder{}
	sb.WriteString("insert into ")
	sb.WriteString(table)
	sb.WriteString("\n(")
	for i, field := range fieldList {
		if field.Default != "" || field.Name == "id" {
			continue
		}
		if i > 0 {
			sb.WriteString("\n")
			iv.WriteString("\n")
		}
		sb.WriteString(field.Name)
		sb.WriteString(",")
		iv.WriteString("#{param.")
		iv.WriteString(stringx.LowerCamelCase(field.Name))
		iv.WriteString("},")
	}
	sb.WriteString(",)\nvalues (")
	sb.WriteString(iv.String())
	sb.WriteString(",)\n\n")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 构建更新sql语句
func BuildUpdateSql(table string, fieldList []*Field) string {
	sb := strings.Builder{}
	sb.WriteString("update ")
	sb.WriteString(table)
	sb.WriteString("\n<set>")
	for _, field := range fieldList {
		if field.Default != "" || field.Name == "id" {
			continue
		}
		paramName := "param." + stringx.LowerCamelCase(field.Name)
		sb.WriteString("\n\t")
		sb.WriteString(fmt.Sprintf(`<if test="%s != null and %s != ''"> %s = #{%s}, </if>`, paramName, paramName, field.Name, paramName))
	}
	sb.WriteString("\n\tupdate_by = #{param.updateBy},")
	sb.WriteString("\n\tupdate_user_name = #{param.updateUserName},")
	sb.WriteString("\n\tupdate_time = #{param.updateTime},")
	sb.WriteString("\n</set>")
	sb.WriteString("\nwhere id = #{item.id} \n\n")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 原始select
func BuildSelectSql(table string, fieldList []*Field) string {
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range fieldList {
		if i > 0 {
			sb.WriteString("\n")
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(",")
	}
	sb.WriteString(",\n  from ")
	sb.WriteString(table)
	sb.WriteString("\n\n")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

// 带字段别名的select
func BuildSelectSqlAlias(table string, columns []*Field) string {
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range columns {
		if i > 0 {
			sb.WriteString("\n")
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(" as ")
		sb.WriteString(stringx.LowerCamelCase(field.Name))
		sb.WriteString(",")
	}
	sb.WriteString(",\n  from ")
	sb.WriteString(table)
	sb.WriteString("\n\n")
	return strings.ReplaceAll(sb.String(), ",,", "")
}
