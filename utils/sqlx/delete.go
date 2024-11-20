package sqlx

func parseDeleteSQL(sql string, indent ...int) *DeleteSqlParser {
	return &DeleteSqlParser{}
}

type DeleteSqlParser struct {
	ParserBase
	Table *Table       // 删除表
	Where []*Condition // 查询条件
}

func (s *DeleteSqlParser) Beautify() string {
	return ""
}
