package gormx

import (
	"reflect"

	"github.com/go-xuan/quanx/os/errorx"
)

type Tabler interface {
	TableName() string    // 表名
	TableComment() string // 表注释
	InitData() any        // 表初始数据
}

// InitTableWithTabler 初始化表结构以及表数据（基于接口实现）
func (h *Handler) InitTableWithTabler(source string, dst ...Tabler) error {
	var db, conf = h.gormMap[source], h.configMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, table := range dst {
				migrator := db.Migrator()
				if migrator.HasTable(table) {
					if err := migrator.AutoMigrate(table); err != nil {
						if err.Error() != "insufficient arguments" {
							return errorx.Wrap(err, "table auto migrate error")
						}
					}
					var count int64
					if err := db.Model(table).Count(&count).Error; err != nil {
						return errorx.Wrap(err, "table count error")
					} else if count == 0 {
						// 初始化表数据
						if initData := table.InitData(); initData != nil {
							if err = db.Create(initData).Error; err != nil {
								return errorx.Wrap(err, "table init data error")
							}
						}
					}
				} else {
					if err := migrator.CreateTable(table); err != nil {
						return errorx.Wrap(err, "table create error")
					}
					if name, comment := table.TableName(), table.TableComment(); name != "" && comment != "" {
						if err := db.Exec(conf.CommentTableSql(name, comment)).Error; err != nil {
							return errorx.Wrap(err, "table alter comment error")
						}
					}
					if initData := table.InitData(); initData != nil {
						if err := db.Create(initData).Error; err != nil {
							return errorx.Wrap(err, "table init data error")
						}
					}
				}
			}
		}
	}
	return nil
}

// InitTableWithAny 初始化表结构（基于反射）
func (h *Handler) InitTableWithAny(source string, dst ...any) error {
	var db, conf = h.gormMap[source], h.configMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, model := range dst {
				if db.Migrator().HasTable(model) {
					if err := db.Migrator().AutoMigrate(model); err != nil {
						if err.Error() != "insufficient arguments" {
							return errorx.Wrap(err, "table auto migrate error")
						}
					}
				} else {
					if err := db.Migrator().CreateTable(model); err != nil {
						return errorx.Wrap(err, "table create error")
					}
					// 添加表备注
					var refValue = reflect.ValueOf(model)
					if method := refValue.MethodByName("TableComment"); method.IsValid() {
						tableName := refValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
						comment := method.Call([]reflect.Value{})[0].String()
						if err := db.Exec(conf.CommentTableSql(tableName, comment)).Error; err != nil {
							return errorx.Wrap(err, "table alter comment error")
						}
					}
					// 初始化表数据
					if method := refValue.MethodByName("InitData"); method.IsValid() {
						initData := method.Call([]reflect.Value{})[0].Interface()
						if err := db.Create(initData).Error; err != nil {
							return errorx.Wrap(err, "table init data error")
						}
					}
				}
			}
		}
	}
	return nil
}
