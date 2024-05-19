package gormx

import (
	"gorm.io/gorm"
	"reflect"

	"github.com/go-xuan/quanx/common/constx"
)

var handler *Handler

// Gorm处理器
type Handler struct {
	Multi     bool // 是否多数据源连接
	Config    *Database
	DB        *gorm.DB
	ConfigMap map[string]*Database
	DBMap     map[string]*gorm.DB
}

func DB(name ...string) *gorm.DB {
	return This().GetDB(name...)
}

func This() *Handler {
	if !Initialized() {
		panic("The gorm handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Initialized() bool {
	return handler != nil
}

func (h *Handler) GetDB(source ...string) *gorm.DB {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if db, ok := h.DBMap[source[0]]; ok {
			return db
		}
	}
	return h.DB
}

func (h *Handler) GetConfig(source ...string) *Database {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}

// 初始化表结构（基于反射）
func (h *Handler) InitTable(source string, dst ...any) (err error) {
	var db, conf = h.DBMap[source], h.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, model := range dst {
				if db.Migrator().HasTable(model) {
					if err = db.Migrator().AutoMigrate(model); err != nil {
						return
					}
				} else {
					if err = db.Migrator().CreateTable(model); err != nil {
						return
					}
					// 添加表备注
					var refValue = reflect.ValueOf(model)
					if method := refValue.MethodByName("TableComment"); method.IsValid() {
						tableName := refValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
						comment := method.Call([]reflect.Value{})[0].String()
						if err = db.Exec(conf.CommentTableSql(tableName, comment)).Error; err != nil {
							return
						}
					}
					// 初始化表数据
					if method := refValue.MethodByName("InitData"); method.IsValid() {
						initData := method.Call([]reflect.Value{})[0].Interface()
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
