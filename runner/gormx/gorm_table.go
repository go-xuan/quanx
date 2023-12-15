package gormx

type Table[T any] interface {
	TableName() string
	Comment() string
	InitData() interface{}
}

// 初始化表结构（基于接口实现）
func (h *Handler) InitGormTable(source string, dst ...Table[any]) (err error) {
	var db, conf = h.DBMap[source], h.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, table := range dst {
				if db.Migrator().HasTable(table) {
					err = db.Migrator().AutoMigrate(table)
					if err != nil {
						return
					}
				} else {
					err = db.Migrator().CreateTable(table)
					if err != nil {
						return
					}
					// 添加表备注
					if tableName, comment := table.TableName(), table.Comment(); tableName != "" && comment != "" {
						err = db.Exec(conf.CommentTableSql(tableName, comment)).Error
						if err != nil {
							return
						}
					}
					// 初始化表数据
					if initData := table.InitData(); initData != nil {
						err = db.Create(initData).Error
						if err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}
