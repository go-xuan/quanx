package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 数据源配置
type Config struct {
	Type      string `json:"type" yaml:"type"`           // 数据库类型
	Host      string `json:"host" yaml:"host"`           // 数据库Host
	Port      string `json:"port" yaml:"port"`           // 数据库端口
	DbName    string `json:"dbName" yaml:"dbName"`       // 数据库名
	UserName  string `json:"userName" yaml:"userName"`   // 用户名
	Password  string `json:"password" yaml:"password"`   // 密码
	Debug     bool   `json:"debug" yaml:"debug"`         // 是否开启debug
	InitTable bool   `json:"initTable" yaml:"initTable"` // 是否初始化表结构
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf(
		"type=%s host=%s port=%s DbName=%s showSql=%v",
		conf.Type, conf.Host, conf.Port, conf.DbName, conf.Debug)
}

// 创建数据库连接
func (conf *Config) NewGormCTL() (ctl *Control, err error) {
	var gormDB *gorm.DB
	gormDB, err = conf.NewGormDB()
	if err != nil {
		return
	}
	ctl = &Control{Config: conf, DB: gormDB}
	return
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
	Mysql      = "mysql"
	Postgres   = "postgres"
	ClickHouse = "clickhouse"
	SqlServer  = "sqlserver"
)

// 获取数据库连接DSN
func (conf *Config) GetDSN() (dsn string) {
	switch strings.ToLower(conf.Type) {
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			conf.UserName,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DbName) +
			"?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local"
	case Postgres:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			conf.Host,
			conf.Port,
			conf.UserName,
			conf.Password,
			conf.DbName)
	case ClickHouse:
		dsn = fmt.Sprintf("tcp://%s:%s?username=%s&password=%s",
			conf.Host,
			conf.Port,
			conf.UserName,
			conf.Password)
	case SqlServer:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			conf.UserName,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DbName)
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
	case ClickHouse:
		return gorm.Open(clickhouse.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	case SqlServer:
		return gorm.Open(sqlserver.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	}
	return
}
