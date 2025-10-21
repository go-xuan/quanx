package gormx

import (
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitTabler 初始化表数据
type InitTabler interface {
	InitData() any
}

// CommentTabler 添加表备注
type CommentTabler interface {
	TableComment() string
}

// InitTable 初始化表
func InitTable(source string, tablers ...interface{}) error {
	if client := GetClient(source); client != nil {
		if db, conf := client.GetInstance(), client.GetConfig(); db != nil && conf != nil && len(tablers) > 0 {
			migrator := db.Migrator()
			for _, tabler := range tablers {
				if schemaTabler, ok := tabler.(schema.Tabler); ok {
					if migrator.HasTable(tabler) {
						if err := db.Migrator().AutoMigrate(tabler); err != nil {
							if err.Error() != "insufficient arguments" {
								return errorx.Wrap(err, "table auto migrate error")
							}
						}
						if err := initTableData(db, tabler); err != nil {
							return errorx.Wrap(err, "table init data error")
						}
					} else {
						if err := db.Migrator().CreateTable(tabler); err != nil {
							return errorx.Wrap(err, "table create error")
						}
						if err := alterTableComment(db, schemaTabler, conf.Type); err != nil {
							return errorx.Wrap(err, "alter table comment error")
						}
						if err := initTableData(db, tabler); err != nil {
							return errorx.Wrap(err, "table init data error")
						}
					}
				} else {
					if err := migrator.AutoMigrate(tabler); err != nil {
						return errorx.Wrap(err, "table create error")
					}
				}
			}
		}
	}
	return nil
}

// 初始化表数据
func initTableData(db *gorm.DB, tabler interface{}) error {
	if initTabler, ok := tabler.(InitTabler); ok {
		var count int64
		if err := db.Model(tabler).Count(&count).Error; err != nil {
			return errorx.Wrap(err, "table count error")
		} else if initData := initTabler.InitData(); initData != nil && count == 0 {
			if err = db.Create(initData).Error; err != nil {
				return errorx.Wrap(err, "table insert error")
			}
		}
	}
	return nil
}

// 添加表备注
func alterTableComment(db *gorm.DB, tabler schema.Tabler, dbType string) error {
	if commentTabler, ok := tabler.(CommentTabler); ok {
		if name, comment := tabler.TableName(), commentTabler.TableComment(); name != "" && comment != "" {
			if err := db.Exec(commentTableSql(dbType, name, comment)).Error; err != nil {
				return errorx.Wrap(err, "table alter comment error")
			}
		}
	}
	return nil
}

func commentTableSql(dbType, table, comment string) string {
	switch strings.ToLower(dbType) {
	case MYSQL:
		return "alter table " + table + " comment = '" + comment + "'"
	case POSTGRES, PGSQL:
		return "comment on table " + table + " is '" + comment + "'"
	}
	return ""
}
