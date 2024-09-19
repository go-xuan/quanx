package gormx

import (
	"github.com/go-xuan/quanx/os/errorx"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/app/constx"
)

var handler *Handler

// Handler Gorm处理器
type Handler struct {
	multi     bool // 是否多数据源连接
	config    *DB
	configMap map[string]*DB
	gormDB    *gorm.DB
	gormMap   map[string]*gorm.DB
}

func GetConfig(source ...string) *DB {
	return This().GetConfig(source...)
}

func GetDB(source ...string) *gorm.DB {
	return This().GetDB(source...)
}

func CloseDB(source ...string) error {
	if sqlDB, err := GetDB(source...).DB(); err != nil {
		return errorx.Wrap(err, "获取*sql.DB连接失败")
	} else if err = sqlDB.Close(); err != nil {
		return errorx.Wrap(err, "关闭*sql.DB连接失败")
	}
	return nil
}

func Close(db *gorm.DB) error {
	if sqlDB, err := db.DB(); err != nil {
		return errorx.Wrap(err, "获取*sql.DB连接失败")
	} else if err = sqlDB.Close(); err != nil {
		return errorx.Wrap(err, "关闭*sql.DB连接失败")
	}
	return nil
}

func This() *Handler {
	if handler == nil {
		panic("the gorm handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Initialized() bool {
	return handler != nil
}

func (h *Handler) GetDB(source ...string) *gorm.DB {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if db, ok := h.gormMap[source[0]]; ok {
			return db
		}
	}
	return h.gormDB
}

func (h *Handler) GetConfig(source ...string) *DB {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if conf, ok := h.configMap[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) Sources() []string {
	var sources []string
	for source := range h.configMap {
		sources = append(sources, source)
	}
	return sources
}
