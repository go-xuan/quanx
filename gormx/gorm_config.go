package gormx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据源配置
type Configs []*Config

func (configs Configs) Format() (f string) {
	bytes, _ := json.Marshal(configs)
	return string(bytes)
}

// 初始化Gorm控制器
func (configs Configs) Init() {
	InitGormX(configs)
}

// 初始化GormX
func InitGormX(configs Configs) {
	if len(configs) == 0 {
		log.Error("数据库配置为空!")
		return
	}
	instance = &Handler{
		DBMap:     make(map[string]*gorm.DB),
		ConfigMap: make(map[string]*Config),
	}
	for i, conf := range configs {
		if conf.Enable {
			if newDB, err := conf.NewGormDB(); err == nil {
				instance.DBMap[conf.Source] = newDB
				instance.ConfigMap[conf.Source] = conf
				log.Info("数据库连接成功! ", conf.Format())
				if i == 0 || conf.Source == "default" {
					instance.DB = newDB
					instance.Config = conf
				}
			} else {
				log.Error("数据库连接失败! ", conf.Format())
				log.Error("error : ", err)
			}
		}
	}
}

// 根据数据源名称获取配置
func (configs Configs) GetConfig(source string) *Config {
	for _, config := range configs {
		if config.Source == source {
			return config
		}
	}
	return nil
}

// 获取默认配置
func (configs Configs) GetDefault() *Config {
	if len(configs) > 0 {
		var conf = configs.GetConfig("default")
		if conf == nil {
			conf = configs[0]
		}
		return conf
	}
	return nil
}

type Config struct {
	Source   string `json:"source" yaml:"source"`     // 数据源名称
	Enable   bool   `json:"enable" yaml:"enable"`     // 数据源启用
	Type     string `json:"type" yaml:"type"`         // 数据库类型
	Host     string `json:"host" yaml:"host"`         // 数据库Host
	Port     int    `json:"port" yaml:"port"`         // 数据库端口
	Database string `json:"database" yaml:"database"` // 数据库名
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Debug    bool   `json:"debug" yaml:"debug"`       // 开启debug（打印SQL以及初始化模型建表）
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("source=%s type=%s host=%s port=%d database=%s debug=%v",
		conf.Source, conf.Type, conf.Host, conf.Port, conf.Database, conf.Debug)
}

// 配置信息格式化
func (conf *Config) Init() {
	if conf.Enable {
		if newDB, err := conf.NewGormDB(); err == nil {
			log.Info("数据库连接成功! ", conf.Format())
			instance.DB = newDB
			instance.Config = conf
			instance.DBMap[conf.Source] = newDB
			instance.ConfigMap[conf.Source] = conf
		} else {
			log.Error("数据库连接失败! ", conf.Format())
			log.Error("error : ", err)
		}
	}
}

// 创建数据库连接
func (conf *Config) NewGormDB() (gormDB *gorm.DB, err error) {
	gormDB, err = GetGormDB(conf.GetDSN(), conf.Type)
	if err != nil {
		return
	}
	var sqlDB *sql.DB
	sqlDB, err = gormDB.DB()
	if err != nil {
		return
	}
	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(60)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 是否打印SQL
	if conf.Debug {
		gormDB = gormDB.Debug()
	}
	return
}

// 数据库类型
const (
	Mysql    = "mysql"
	Postgres = "postgres"
)

// 获取数据库连接DSN
func (conf *Config) GetDSN() (dsn string) {
	switch strings.ToLower(conf.Type) {
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Database) +
			"?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local"
	case Postgres:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.Username, conf.Password, conf.Database)
	}
	return
}

// 根据dsn生成gormDB
func GetGormDB(dsn, dialect string) (gormDb *gorm.DB, err error) {
	switch strings.ToLower(dialect) {
	case Mysql:
		return gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	case Postgres:
		return gorm.Open(postgres.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	}
	return
}
