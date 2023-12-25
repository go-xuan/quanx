package gormx

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

var handler *Handler

// Gorm处理器
type Handler struct {
	Multi     bool // 是否多数据源连接
	DB        *gorm.DB
	Config    *Database
	DBMap     map[string]*gorm.DB
	ConfigMap map[string]*Database
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
	if len(source) > 0 {
		if db, ok := h.DBMap[source[0]]; ok {
			return db
		}
	}
	return h.DB
}

func (h *Handler) GetConfig(source ...string) *Database {
	if len(source) > 0 {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}

// 初始化表结构（基于反射）
func (h *Handler) InitTable(source string, dst ...interface{}) (err error) {
	var db, conf = h.DBMap[source], h.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Debug {
			for _, model := range dst {
				if db.Migrator().HasTable(model) {
					err = db.Migrator().AutoMigrate(model)
					if err != nil {
						return
					}
				} else {
					err = db.Migrator().CreateTable(model)
					if err != nil {
						return
					}
					// 添加表备注
					var refValue = reflect.ValueOf(model)
					if method := refValue.MethodByName("Comment"); method.IsValid() {
						tableName := refValue.MethodByName("TableName").Call([]reflect.Value{})[0].String()
						comment := method.Call([]reflect.Value{})[0].String()
						err = db.Exec(fmt.Sprintf("COMMENT ON TABLE %s IS '%s'", tableName, comment)).Error
						if err != nil {
							return
						}
					}
					// 初始化表数据
					if method := refValue.MethodByName("InitData"); method.IsValid() {
						initData := method.Call([]reflect.Value{})[0].Interface()
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
