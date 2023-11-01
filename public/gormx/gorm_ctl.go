package gormx

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var CTL *Controller

// Gorm控制器
type Controller struct {
	Multi     bool // 是否多连接
	DB        *gorm.DB
	Config    *Config
	DBMap     map[string]*gorm.DB
	ConfigMap map[string]*Config
}

// 初始化Gorm控制器
func Init(configs Configs) {
	if len(configs) == 0 {
		log.Error("数据库配置为空!")
		return
	}
	CTL = &Controller{
		DBMap:     make(map[string]*gorm.DB),
		ConfigMap: make(map[string]*Config),
	}
	for i, conf := range configs {
		if conf.Enable {
			if newDB, err := conf.NewGormDB(); err == nil {
				CTL.DBMap[conf.Source] = newDB
				CTL.ConfigMap[conf.Source] = conf
				log.Info("gorm连接成功! ", conf.Format())
				if i == 0 || conf.Source == "default" {
					CTL.DB = newDB
					CTL.Config = conf
				}
			} else {
				log.Error("gorm连接失败! ", conf.Format())
				log.Error("error : ", err)
			}
		}
	}
}

// 初始化表结构
func (ctl *Controller) InitTable(source string, dst ...interface{}) (err error) {
	var db, conf = ctl.DBMap[source], ctl.ConfigMap[source]
	if db != nil && conf != nil && len(dst) > 0 {
		if conf.Init {
			for _, item := range dst {
				if db.Migrator().HasTable(item) {
					err = db.Migrator().AutoMigrate(item)
				} else {
					err = db.Migrator().CreateTable(item)
				}
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (ctl *Controller) GetDB(source ...string) *gorm.DB {
	if len(source) > 0 {
		if db, ok := ctl.DBMap[source[0]]; ok {
			return db
		}
	}
	return ctl.DB
}

func (ctl *Controller) GetConfig(source ...string) *Config {
	if len(source) > 0 {
		if conf, ok := ctl.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return ctl.Config
}
