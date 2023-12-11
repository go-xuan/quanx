package gormx

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

var instance *Handler

// Gorm处理器
type Handler struct {
	Multi     bool // 是否多连接
	DB        *gorm.DB
	Config    *Config
	DBMap     map[string]*gorm.DB
	ConfigMap map[string]*Config
}

func This() *Handler {
	if instance == nil {
		panic("The gorm instance has not been initialized, please check the relevant config")
	}
	return instance
}

// 初始化表结构
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

func (h *Handler) GetDB(source ...string) *gorm.DB {
	if len(source) > 0 {
		if db, ok := h.DBMap[source[0]]; ok {
			return db
		}
	}
	return h.DB
}

func (h *Handler) GetConfig(source ...string) *Config {
	if len(source) > 0 {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}
