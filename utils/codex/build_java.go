package codex

import (
	"strings"

	"github.com/quanxiaoxuan/quanx/utils/sqlx"
	"github.com/quanxiaoxuan/quanx/utils/stringx"
)

// 构建java实体类
func BuildJavaClass(table string, fieldList FieldList) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	sb := strings.Builder{}
	sb.WriteString(LineSep)
	sb.WriteString("@Data")
	sb.WriteString(LineSep)
	sb.WriteString("public class ")
	sb.WriteString(stringx.UpperCamelCase(table))
	sb.WriteString(" implements Serializable {")
	for _, field := range fieldList {
		low := stringx.LowerCamelCase(field.Name)
		jt := sqlx.Pg2JavaTypeMap[field.Type]
		sb.WriteString(LineSep)
		sb.WriteString(Tab)
		sb.WriteString("@ApiModelProperty(value = \"")
		sb.WriteString(field.Comment)
		sb.WriteString("\", name = \"")
		sb.WriteString(low)
		sb.WriteString("\")")
		sb.WriteString(LineSep)
		sb.WriteString(Tab)
		sb.WriteString("private ")
		sb.WriteString(jt)
		sb.WriteString(" ")
		sb.WriteString(low)
		sb.WriteString(";")
	}
	sb.WriteString(LineSep)
	sb.WriteString("}")
	sb.WriteString(LineSep)
	return sb.String()
}
