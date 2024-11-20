package sqlx

import "fmt"

type Condition struct {
	Type    string // or/and
	Content string //
}

type Join struct {
	Table *Table // join表对象
	Type  string // join类型left/right/inner
	On    string // 关联条件
}

type Table struct {
	Select *SelectSqlParser // 子查询
	Name   string           // 表名
	Alias  string           // 表别名
	schema string           // schema名
}

type Field struct {
	Name      string // 字段名
	Alias     string // 字段别名，仅查询使用
	Table     string // 表名
	Value     string // 字段值
	Type      string // 字段类型
	Precision int    // 长度
	Scale     int    // 小数点
	Nullable  bool   // 允许为空
	Default   string // 默认值
	Comment   string // 注释
}

func (t *Table) WithAliasSQL() string {
	if t.Select != nil {
		return fmt.Sprintf("( %s ) as %s", t.Select.Beautify(), t.Alias)
	} else {
		return fmt.Sprintf("%s as %s", t.Name, t.Alias)
	}
}
