package gormx

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var CTL *Controller

// Gorm控制器
type Controller struct {
	Config *Config
	DB     *gorm.DB
	Init   func()
}

// 初始化表结构
func (ctl *Controller) InitTable(dst interface{}) error {
	if ctl.DB.Migrator().HasTable(dst) {
		return ctl.DB.Migrator().AutoMigrate(dst)
	} else {
		return ctl.DB.Migrator().CreateTable(dst)
	}
}

// 初始化方法
func (ctl *Controller) SetInit(f func()) {
	ctl.Init = f
}

// 执行方法
func (ctl *Controller) ExecInit() {
	if ctl.Init != nil {
		ctl.Init()
	}
}

// 初始化Gorm控制器
func Init(conf *Config) {
	if conf.Type == "" {
		log.Error("数据库类型配置不可为空 : ", conf.Format())
		return
	}
	newDB, err := conf.NewGormDB()
	if err != nil {
		log.Error("gorm连接失败! ", conf.Format())
		log.Error("error : ", err)
	} else {
		CTL = &Controller{Config: conf, DB: newDB}
		log.Info("gorm连接成功! ", conf.Format())
	}
}
