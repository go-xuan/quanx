package sqlx

func parseCreateSQL(sql string, indent ...int) *CreateSqlParser {
	return &CreateSqlParser{}
}

type CreateSqlParser struct {
	originSql string       // 原始sql
	table     *Table       // 创建表
	fields    []*Field     // 表字段
	where     []*Condition // 查询条件
}

func (s *CreateSqlParser) Beautify() string {
	return ""
}
