package dbx

import (
	"fmt"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"
)

// Tabler 表名接口
type Tabler interface {
	TableName() string
}

// InitTabler 初始化表数据接口
type InitTabler interface {
	InitData() any
}

// CommentTabler 添加表备注接口
type CommentTabler interface {
	TableComment() string
}

// InitTable 初始化表
func InitTable(tables map[string][]interface{}) error {
	if !Initialized() || len(tables) == 0 {
		return nil
	}
	var err error
	this().Range(func(source string, client Client) bool {
		if list, ok := tables[source]; ok && len(list) > 0 {
			if err = InitGormTable(client, list...); err != nil {
				return true
			}
		}
		return false
	})
	if err != nil {
		return errorx.Wrap(err, "init tables failed")
	}
	return nil
}

// InitGormTable 初始化gorm表
func InitGormTable(client Client, tables ...interface{}) error {
	if client == nil || len(tables) == 0 {
		return nil
	}
	db, ok := client.GetInstance().(*gorm.DB)
	if !ok {
		return nil
	}
	migrator, typ := db.Migrator(), client.GetConfig().Type
	for _, table := range tables {
		if tabler, ok := table.(Tabler); ok {
			if migrator.HasTable(table) {
				if err := migrator.AutoMigrate(table); err != nil {
					if err.Error() != "insufficient arguments" {
						return errorx.Wrap(err, "auto migrate table error")
					}
				}
				if err := initGormTableData(db, tabler); err != nil {
					return errorx.Wrap(err, "init table data error")
				}
			} else {
				if err := migrator.CreateTable(table); err != nil {
					return errorx.Wrap(err, "create table error")
				}
				if err := alterGormTableComment(db, tabler, typ); err != nil {
					return errorx.Wrap(err, "alter table comment error")
				}
				if err := initGormTableData(db, tabler); err != nil {
					return errorx.Wrap(err, "init table data error")
				}
			}
		} else {
			if err := migrator.AutoMigrate(table); err != nil {
				return errorx.Wrap(err, "create table error")
			}
		}
	}
	return nil
}

// 初始化表数据
func initGormTableData(db *gorm.DB, tabler Tabler) error {
	if initTabler, ok := tabler.(InitTabler); ok {
		var count int64
		if err := db.Model(tabler).Count(&count).Error; err != nil {
			return errorx.Wrap(err, "count table data error")
		} else if data := initTabler.InitData(); data != nil && count == 0 {
			if err = db.Create(data).Error; err != nil {
				return errorx.Wrap(err, "insert table data error")
			}
		}
	}
	return nil
}

// 添加表备注
func alterGormTableComment(db *gorm.DB, tabler Tabler, typ string) error {
	if commentTabler, ok := tabler.(CommentTabler); ok {
		if name, comment := tabler.TableName(), commentTabler.TableComment(); name != "" && comment != "" {
			var sql string
			switch strings.ToLower(typ) {
			case MYSQL:
				sql = fmt.Sprintf("alter table %s comment = '%s'", name, comment)
			case POSTGRES, PGSQL:
				sql = fmt.Sprintf("comment on table %s is '%s'", name, comment)
			default:
				return errorx.Sprintf("unsupported database type: %s", typ)
			}
			if err := db.Exec(sql).Error; err != nil {
				return errorx.Wrap(err, "alter table comment error")
			}
		}
	}
	return nil
}
