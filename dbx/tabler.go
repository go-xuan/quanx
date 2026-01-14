package dbx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitTabler 初始化表数据接口
type InitTabler interface {
	InitData() interface{}
}

// CommentTabler 添加表备注接口
type CommentTabler interface {
	TableComment() string
}

// InitGormTable 初始化gorm表
func InitGormTable(db *gorm.DB, tables ...interface{}) error {
	if db == nil || len(tables) == 0 {
		return nil
	}
	migrator := db.Migrator()
	for _, table := range tables {
		if migrator.HasTable(table) {
			if err := migrator.AutoMigrate(table); err != nil {
				if err.Error() != "insufficient arguments" {
					return errorx.Wrap(err, "auto migrate table failed")
				}
			}
		} else {
			if err := migrator.CreateTable(table); err != nil {
				return errorx.Wrap(err, "create table failed")
			}
			if err := alterGormTableComment(db, table); err != nil {
				return errorx.Wrap(err, "alter table comment failed")
			}
		}
		if err := initGormTableData(db, table); err != nil {
			return errorx.Wrap(err, "init table data failed")
		}
	}
	return nil
}

// 初始化表数据
func initGormTableData(db *gorm.DB, table interface{}) error {
	var data interface{}
	if tabler, ok := table.(InitTabler); !ok {
		return nil
	} else if data = tabler.InitData(); data == nil {
		return nil
	}
	var count int64
	if err := db.Model(table).Count(&count).Error; err != nil {
		return errorx.Wrap(err, "count table data failed")
	} else if count == 0 {
		if err = db.Create(data).Error; err != nil {
			return errorx.Wrap(err, "insert table data failed")
		}
	}
	return nil
}

// 添加表备注
func alterGormTableComment(db *gorm.DB, table interface{}) error {
	var name, comment string
	if tabler, ok := table.(schema.Tabler); ok {
		if name = tabler.TableName(); name == "" {
			return nil
		}
	}
	if tabler, ok := table.(CommentTabler); ok {
		if comment = tabler.TableComment(); comment == "" {
			return nil
		}
	}
	var sql string
	switch typ := db.Name(); typ {
	case MYSQL:
		sql = fmt.Sprintf("alter table %s comment = '%s'", name, comment)
	case POSTGRES, PGSQL:
		sql = fmt.Sprintf("comment on table %s is '%s'", name, comment)
	default:
		return errorx.Sprintf("unsupported database type: %s", typ)
	}
	if err := db.Exec(sql).Error; err != nil {
		return errorx.Wrap(err, "alter table comment failed")
	}
	return nil
}
