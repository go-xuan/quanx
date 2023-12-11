package codex

import (
	"strings"

	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/utilx/sqlx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

// 构建java实体类
func BuildJavaClass(table string, fieldList []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	sb := strings.Builder{}
	sb.WriteString(constx.NextLine)
	sb.WriteString("@Data")
	sb.WriteString(constx.NextLine)
	sb.WriteString("public class ")
	sb.WriteString(stringx.UpperCamelCase(table))
	sb.WriteString(" implements Serializable {")
	for _, field := range fieldList {
		low := stringx.LowerCamelCase(field.Name)
		sb.WriteString(constx.NextLine)
		sb.WriteString(constx.Tab)
		sb.WriteString("@ApiModelProperty(value = \"")
		sb.WriteString(field.Comment)
		sb.WriteString("\", name = \"")
		sb.WriteString(low)
		sb.WriteString("\")")
		sb.WriteString(constx.NextLine)
		sb.WriteString(constx.Tab)
		sb.WriteString("private ")
		sb.WriteString(sqlx.Pg2JavaTypeMap[field.Type])
		sb.WriteString(" ")
		sb.WriteString(low)
		sb.WriteString(";")
	}
	sb.WriteString(constx.NextLine)
	sb.WriteString("}")
	sb.WriteString(constx.NextLine)
	return sb.String()
}
