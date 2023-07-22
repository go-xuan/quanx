package gormx

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var CTL *Control

// Gorm控制器
type Control struct {
	Config *Config
	DB     *gorm.DB
	Init   func()
}

// 初始化表结构
func (ctl *Control) InitTable(dst interface{}) error {
	if ctl.DB.Migrator().HasTable(dst) {
		return ctl.DB.Migrator().AutoMigrate(dst)
	} else {
		return ctl.DB.Migrator().CreateTable(dst)
	}
}

// 初始化方法
func (ctl *Control) SetInit(f func()) {
	ctl.Init = f
}

// 执行方法
func (ctl *Control) ExecInit() {
	if ctl.Init != nil {
		ctl.Init()
	}
}

// 初始化Gorm控制器
func InitGormCTL(conf *Config) {
	if conf.Type == "" {
		return
	}
	var err error
	msg := conf.Format()
	if CTL == nil {
		CTL, err = conf.NewGormCTL()
		if err != nil {
			log.Error("初始化gorm连接-失败! ", msg)
			log.Error("error : ", err)
		} else {
			log.Info("初始化gorm连接-成功! ", msg)
		}
	} else {
		var newDB *gorm.DB
		newDB, err = conf.NewGormDB()
		if err != nil {
			log.Error("更新gorm连接-失败! ", msg)
			log.Error("error : ", err)
		} else {
			CTL.DB = newDB
			CTL.Config = conf
			log.Error("更新gorm连接-成功! ", msg)
		}
	}
}
