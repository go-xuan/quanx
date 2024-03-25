package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/utilx/anyx"
)

// 数据源配置
type MultiDatabase []*Database

// 配置器名称
func (MultiDatabase) Title() string {
	return "init multi-database"
}

func (MultiDatabase) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "multi_database.yaml",
		NacosDataId: "multi_database.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (md MultiDatabase) Run() (err error) {
	if len(md) == 0 {
		log.Info("database not connected! reason: multi-database.yaml not found!")
		return
	}
	handler = &Handler{
		Multi:     true,
		DBMap:     make(map[string]*gorm.DB),
		ConfigMap: make(map[string]*Database),
	}
	for i, d := range md {
		if d.Enable {
			if err = anyx.SetDefaultValue(d); err != nil {
				return
			}
			var db *gorm.DB
			if db, err = d.NewGormDB(); err != nil {
				log.Error(d.ToString("database connect failed!"))
				return err
			}
			handler.DBMap[d.Source] = db
			handler.ConfigMap[d.Source] = d
			if i == 0 || d.Source == "default" {
				handler.DB = db
				handler.Config = d
			}
			log.Info(d.ToString("database connect successful!"))
		}
	}
	if len(handler.ConfigMap) == 0 {
		log.Info("database not connected! reason: multi-database.yaml is empty or {database.enable} is false")
	}
	return
}

// 根据数据源名称获取配置
func (md MultiDatabase) GetConfig(source string) *Database {
	for _, config := range md {
		if config.Source == source {
			return config
		}
	}
	return nil
}

// 获取默认配置
func (md MultiDatabase) GetDefault() *Database {
	if len(md) > 0 {
		var conf = md.GetConfig("default")
		if conf == nil {
			conf = md[0]
		}
		return conf
	}
	return nil
}

type Database struct {
	Source          string `json:"source" yaml:"source" default:"default"`              // 数据源名称
	Enable          bool   `json:"enable" yaml:"enable"`                                // 数据源启用
	Type            string `json:"type" yaml:"type"`                                    // 数据库类型
	Host            string `json:"host" yaml:"host" default:"localhost"`                // 数据库Host
	Port            int    `json:"port" yaml:"port"`                                    // 数据库端口
	Username        string `json:"username" yaml:"username"`                            // 用户名
	Password        string `json:"password" yaml:"password"`                            // 密码
	Database        string `json:"database" yaml:"database"`                            // 数据库名
	Schema          string `json:"schema" yaml:"schema"`                                // schema模式名
	Debug           bool   `json:"debug" yaml:"debug" default:"false"`                  // 开启debug（打印SQL以及初始化模型建表）
	MaxIdleConns    int    `json:"maxIdleConns" yaml:"maxIdleConns" default:"10"`       // 最大空闲连接
	MaxOpenConns    int    `json:"maxOpenConns" yaml:"maxOpenConns" default:"10"`       // 最大打开连接
	ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime" default:"10"` // 连接存活时间(分钟)
}

// 配置信息格式化
func (d *Database) ToString(title string) string {
	return fmt.Sprintf("%s => source=%s type=%s host=%s port=%d database=%s debug=%v",
		title, d.Source, d.Type, d.Host, d.Port, d.Database, d.Debug)
}

// 配置器名称
func (d *Database) Title() string {
	return "init database"
}

func (d *Database) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "database.yaml",
		NacosDataId: "database.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (d *Database) Run() (err error) {
	if d.Enable {
		if err = anyx.SetDefaultValue(d); err != nil {
			return
		}
		var db *gorm.DB
		if db, err = d.NewGormDB(); err != nil {
			log.Error(d.ToString("database connect failed!"))
			log.Error("error : ", err)
			return
		}
		handler = &Handler{
			Multi:     false,
			DB:        db,
			Config:    d,
			DBMap:     make(map[string]*gorm.DB),
			ConfigMap: make(map[string]*Database),
		}
		handler.DBMap[d.Source] = db
		handler.ConfigMap[d.Source] = d
		log.Info(d.ToString("database connect successful!"))
		return
	}
	log.Info("database not connected! reason: database.yaml is empty or {database.enable} is false")
	return
}

// 创建数据库连接
func (d *Database) NewGormDB() (gormDB *gorm.DB, err error) {
	if gormDB, err = d.GetGormDB(); err != nil {
		return
	}
	var sqlDB *sql.DB
	if sqlDB, err = gormDB.DB(); err != nil {
		return
	}
	sqlDB.SetMaxIdleConns(d.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(d.ConnMaxLifetime) * time.Minute)
	// 是否打印SQL
	if d.Debug {
		gormDB = gormDB.Debug()
	}
	return
}

// 数据库类型
const (
	Mysql    = "mysql"
	Postgres = "postgres"
)

// 生成表备注
func (d *Database) CommentTableSql(table, comment string) string {
	switch strings.ToLower(d.Type) {
	case Mysql:
		return "ALTER TABLE " + table + " COMMENT = '" + comment + "'"
	case Postgres:
		return "COMMENT ON TABLE " + table + " IS '" + comment + "'"
	}
	return ""
}

// 根据dsn生成gormDB
func (d *Database) GetGormDB() (gormDb *gorm.DB, err error) {
	var dial gorm.Dialector
	switch strings.ToLower(d.Type) {
	case Mysql:
		dial = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Database))
	case Postgres:
		dial = postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Password, d.Database))
	}
	return gorm.Open(dial, &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
}
