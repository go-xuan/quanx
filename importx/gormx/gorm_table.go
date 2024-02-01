package gormx

type Tabler[T any] interface {
	TableName() string     // 表名
	TableComment() string  // 表注释
	InitData() interface{} // 表初始数据
}

// 初始化表结构（基于接口实现）
func (h *Handler) InitGormTable(source string, dst ...Tabler[any]) (err error) {
	var db, conf = h.DBMap[source], h.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, table := range dst {
				if db.Migrator().HasTable(table) {
					if err = db.Migrator().AutoMigrate(table); err != nil {
						return
					}
					var count int64
					if err = db.Model(table).Count(&count).Error; err != nil {
						return
					} else if count == 0 {
						// 初始化表数据
						if initData := table.InitData(); initData != nil {
							err = db.Create(initData).Error
							if err != nil {
								return
							}
						}
					}
				} else {
					if err = db.Migrator().CreateTable(table); err != nil {
						return
					}
					// 添加表备注
					if name, comment := table.TableName(), table.TableComment(); name != "" && comment != "" {
						err = db.Exec(conf.CommentTableSql(name, comment)).Error
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
