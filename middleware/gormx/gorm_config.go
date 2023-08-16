package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据源配置
type Config struct {
	Type      string `json:"type" yaml:"type"`           // 数据库类型
	Host      string `json:"host" yaml:"host"`           // 数据库Host
	Port      int    `json:"port" yaml:"port"`           // 数据库端口
	Database  string `json:"database" yaml:"database"`   // 数据库名
	UserName  string `json:"userName" yaml:"userName"`   // 用户名
	Password  string `json:"password" yaml:"password"`   // 密码
	Debug     bool   `json:"debug" yaml:"debug"`         // 是否开启debug
	InitTable bool   `json:"initTable" yaml:"initTable"` // 是否初始化表结构
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("type=%s host=%s port=%d database=%s showSql=%v",
		conf.Type, conf.Host, conf.Port, conf.Database, conf.Debug)
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
			conf.UserName, conf.Password, conf.Host, conf.Port, conf.Database) +
			"?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local"
	case Postgres:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conf.Host, conf.Port, conf.UserName, conf.Password, conf.Database)
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
