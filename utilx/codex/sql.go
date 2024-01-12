package codex

import (
	"fmt"
	"strings"

	"github.com/go-xuan/quanx/utilx/stringx"
)

var BaseField = []string{
	"create_time",
	"update_time",
	"create_user_id",
	"update_user_id",
	"create_by",
	"update_by",
}

// 构建新增sql语句
func BuildInsertSql(table string, fields []*Field) string {
	sb := strings.Builder{}
	iv := strings.Builder{}
	sb.WriteString("insert into ")
	sb.WriteString(table)
	sb.WriteString("\n(")
	var i int
	for _, field := range fields {
		if (field.Name == "id" || IsBaseField(field.Name)) && field.Default != "" {
			continue
		}
		if i > 0 {
			sb.WriteString("\n")
			iv.WriteString("\n")
		}
		sb.WriteString(field.Name)
		sb.WriteString(",")
		iv.WriteString("#{create.")
		iv.WriteString(stringx.ToLowerCamel(field.Name))
		iv.WriteString("},")
		i++
	}
	sb.WriteString(",)\nvalues \n(")
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
		if field.Name == "id" || IsBaseField(field.Name) {
			continue
		}
		name := "update." + stringx.ToLowerCamel(field.Name)
		sb.WriteString("\n\t")
		sb.WriteString(fmt.Sprintf(`<if test="%s != null and %s != ''"> %s = #{%s}, </if>`, name, name, field.Name, name))
	}
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

// 别名查询select
func BuildSelectAliasSql(table string, columns []*Field) string {
	sb := strings.Builder{}
	sb.WriteString("select ")
	for i, field := range columns {
		if i > 0 {
			sb.WriteString("\n")
			sb.WriteString("       ")
		}
		sb.WriteString(field.Name)
		sb.WriteString(" as ")
		sb.WriteString(stringx.ToLowerCamel(field.Name))
		sb.WriteString(",")
	}
	sb.WriteString(",\n  from ")
	sb.WriteString(table)
	sb.WriteString("\n\n")
	return strings.ReplaceAll(sb.String(), ",,", "")
}

func IsBaseField(name string) bool {
	for _, s := range BaseField {
		if s == name {
			return true
		}
	}
	return false
}
