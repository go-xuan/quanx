package gormx

type Tabler interface {
	TableName() string    // 表名
	TableComment() string // 表注释
	InitData() any        // 表初始数据
}

// InitGormTable 初始化表结构以及表数据（基于接口实现）
func (h *Handler) InitGormTable(source string, dst ...Tabler) (err error) {
	var db, conf = h.DBMap[source], h.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, table := range dst {
				migrator := db.Migrator()
				if migrator.HasTable(table) {
					if err = migrator.AutoMigrate(table); err != nil {
						return
					}
					var count int64
					if err = db.Model(table).Count(&count).Error; err != nil {
						return
					} else if count == 0 {
						// 初始化表数据
						if initData := table.InitData(); initData != nil {
							if err = db.Create(initData).Error; err != nil {
								return
							}
						}
					}
				} else {
					if err = migrator.CreateTable(table); err != nil {
						return
					}
					// 添加表备注
					if name, comment := table.TableName(), table.TableComment(); name != "" && comment != "" {
						if err = db.Exec(conf.CommentTableSql(name, comment)).Error; err != nil {
							return
						}
					}
					// 初始化表数据
					if initData := table.InitData(); initData != nil {
						if err = db.Create(initData).Error; err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}
