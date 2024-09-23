package gormx

import (
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/os/errorx"
)

var _handler *Handler

func Initialized() bool {
	return _handler != nil
}

func this() *Handler {
	if _handler == nil {
		panic("the gorm handler has not been initialized, please check the relevant config")
	}
	return _handler
}

// GetConfig 获取数据库配置
func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

// GetDB 获取数据库连接
func GetDB(source ...string) *gorm.DB {
	return this().GetDB(source...)
}

// Sources 所有数据源
func Sources() []string {
	return this().Sources()
}

// InitTable 初始化表结构以及表数据
func InitTable(source string, dst ...Tabler) error {
	return this().InitTableWithTabler(source, dst...)
}

// Close 关闭数据库连接
func Close(source ...string) error {
	return this().Close(source...)
}

// Handler 数据库连接句柄
type Handler struct {
	multi     bool // 是否多数据源连接
	config    *Config
	configMap map[string]*Config
	gormDB    *gorm.DB
	gormMap   map[string]*gorm.DB
}

func (h *Handler) GetDB(source ...string) *gorm.DB {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if db, ok := h.gormMap[source[0]]; ok {
			return db
		}
	}
	return h.gormDB
}

func (h *Handler) GetConfig(source ...string) *Config {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configMap[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) Close(source ...string) error {
	if sqlDB, err := h.GetDB(source...).DB(); err != nil {
		return errorx.Wrap(err, "get sql.Config failed")
	} else if err = sqlDB.Close(); err != nil {
		return errorx.Wrap(err, "close sql.Config failed")
	}
	return nil
}

func (h *Handler) Sources() []string {
	var sources []string
	for source := range h.configMap {
		sources = append(sources, source)
	}
	return sources
}
