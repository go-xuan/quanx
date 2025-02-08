package gormx

import (
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/os/errorx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the gorm handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func addDB(source string, config *Config, db *gorm.DB) {
	if _handler == nil {
		_handler = &Handler{
			multi:   false,
			config:  config,
			db:      db,
			configs: make(map[string]*Config),
			dbs:     make(map[string]*gorm.DB),
		}
	} else {
		_handler.multi = true
	}
	_handler.configs[source] = config
	_handler.dbs[source] = db
}

// Handler 数据库连接句柄
type Handler struct {
	multi   bool
	config  *Config
	db      *gorm.DB
	configs map[string]*Config
	dbs     map[string]*gorm.DB
}

func (h *Handler) DB(source ...string) *gorm.DB {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if db, ok := h.dbs[source[0]]; ok {
			return db
		}
	}
	return h.db
}

func (h *Handler) GetConfig(source ...string) *Config {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) Close(source ...string) error {
	if sqlDB, err := h.DB(source...).DB(); err != nil {
		return errorx.Wrap(err, "get sql.Config failed")
	} else if err = sqlDB.Close(); err != nil {
		return errorx.Wrap(err, "close sql.Config failed")
	}
	return nil
}

func (h *Handler) Sources() []string {
	var sources []string
	for source := range h.configs {
		sources = append(sources, source)
	}
	return sources
}

// IsInitialized 是否初始化
func IsInitialized() bool {
	return _handler != nil
}

// GetConfig 获取数据库配置
func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

// DB 获取数据库连接
func DB(source ...string) *gorm.DB {
	return this().DB(source...)
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
