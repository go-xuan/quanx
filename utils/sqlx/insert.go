package sqlx

func parseInsertSQL(sql string, indent ...int) *InsertSqlParser {
	return &InsertSqlParser{}
}

type InsertSqlParser struct {
	ParserBase
	Table  *Table           // 插入表
	Fields []*Field         // 插入字段
	Select *SelectSqlParser // 查询条件
}

func (s *InsertSqlParser) Beautify() string {
	return ""
}
