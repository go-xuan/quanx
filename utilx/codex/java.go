package codex

import (
	"strings"

	"github.com/go-xuan/quanx/utilx/sqlx"
	"github.com/go-xuan/quanx/utilx/stringx"
)

// 构建java实体类
func BuildJavaClass(table string, fieldList []*Field) string {
	table = strings.TrimPrefix(table, `t_`)
	table = strings.TrimSuffix(table, `_t`)
	sb := strings.Builder{}
	sb.WriteString("\n@Data\npublic class ")
	sb.WriteString(stringx.UpperCamelCase(table))
	sb.WriteString(" implements Serializable {")
	for _, field := range fieldList {
		low := stringx.LowerCamelCase(field.Name)
		sb.WriteString("\n\t")
		sb.WriteString("@ApiModelProperty(value = \"")
		sb.WriteString(field.Comment)
		sb.WriteString("\", name = \"")
		sb.WriteString(low)
		sb.WriteString("\")\n\tprivate ")
		sb.WriteString(sqlx.Pg2JavaTypeMap[field.Type])
		sb.WriteString(" ")
		sb.WriteString(low)
		sb.WriteString(";")
	}
	sb.WriteString("\n}\n\n")
	return sb.String()
}
