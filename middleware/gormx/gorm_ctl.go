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
}

// 初始化表结构
func (ctl *Controller) InitTable(dst ...interface{}) (err error) {
	if len(dst) > 0 {
		for _, item := range dst {
			if ctl.DB.Migrator().HasTable(item) {
				err = ctl.DB.Migrator().AutoMigrate(item)
			} else {
				err = ctl.DB.Migrator().CreateTable(item)
			}
			if err != nil {
				return
			}
		}
	}
	return
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
