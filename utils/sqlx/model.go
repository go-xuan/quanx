package sqlx

import "fmt"

// ConditionParser 查询条件解析
type ConditionParser struct {
	Type    string // or/and
	Content string //
}

// JoinParser 关联表解析
type JoinParser struct {
	Table *TableParser // join表对象
	Type  string       // join类型left/right/inner
	On    string       // 关联条件
}

// TableParser 主表解析
type TableParser struct {
	Select *SelectParser // 子查询
	Name   string        // 表名
	Alias  string        // 表别名
	Schema string        // schema名
}

func (p *TableParser) WithAliasSQL() string {
	if p.Select != nil {
		return fmt.Sprintf("( %s ) as %s", p.Select.Beautify(), p.Alias)
	} else {
		return fmt.Sprintf("%s as %s", p.Name, p.Alias)
	}
}

// FieldParser 字段解析
type FieldParser struct {
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
