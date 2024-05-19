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

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/types/anyx"
)

// 数据源配置
type MultiDatabase []*Database

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

// 配置器名称
func (MultiDatabase) Theme() string {
	return "Database"
}

// 配置文件读取
func (MultiDatabase) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "database.yaml",
		NacosDataId: "database.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (c MultiDatabase) Run() (err error) {
	if len(c) == 0 {
		log.Error("Database Connect Failed! Reason: database.yaml Not Found")
		return
	}
	if handler == nil {
		handler = &Handler{
			Multi:     true,
			DBMap:     make(map[string]*gorm.DB),
			ConfigMap: make(map[string]*Database),
		}
	} else {
		handler.Multi = true
	}
	for i, d := range c {
		if d.Enable {
			if err = anyx.SetDefaultValue(d); err != nil {
				return
			}
			var db *gorm.DB
			var toString = d.ToString()
			if db, err = d.NewGormDB(); err != nil {
				log.Error("Database Connect Failed: ", toString, err)
				return err
			}
			handler.DBMap[d.Source] = db
			handler.ConfigMap[d.Source] = d
			if i == 0 || d.Source == constx.DefaultKey {
				handler.DB = db
				handler.Config = d
			}
			log.Info("Database Connect Successful: ", toString)
		}
	}
	if len(handler.ConfigMap) == 0 {
		log.Error("Database Connect Failed! reason: database.yaml is empty or all enable values are false")
	}
	return
}

// 配置信息格式化
func (d *Database) ToString() string {
	return fmt.Sprintf("source=%s type=%s host=%s port=%d database=%s debug=%v",
		d.Source, d.Type, d.Host, d.Port, d.Database, d.Debug)
}

// 配置器名称
func (d *Database) Theme() string {
	return "Database"
}

// 配置文件读取
func (d *Database) Reader() *confx.Reader {
	return &confx.Reader{
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
		var toString = d.ToString()
		if db, err = d.NewGormDB(); err != nil {
			log.Error("Database Connect Failed: ", toString, err)
			return
		}

		if handler == nil {
			handler = &Handler{
				Multi:     false,
				DB:        db,
				Config:    d,
				DBMap:     map[string]*gorm.DB{},
				ConfigMap: make(map[string]*Database),
			}
		} else {
			handler.Multi = true
		}
		handler.DBMap[d.Source] = db
		handler.ConfigMap[d.Source] = d

		log.Info("Database Connect Successful: ", toString)
		return
	}
	log.Info("Database Connect Failed! reason: database.yaml is empty or the value of enable is false")
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
		return "alter table " + table + " comment = '" + comment + "'"
	case Postgres:
		return "comment on table " + table + " is '" + comment + "'"
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
		dial = postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Port, d.Username, d.Password, d.Database))
	}
	return gorm.Open(dial, &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
}
