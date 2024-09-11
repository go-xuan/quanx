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

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/app/constx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

// MultiDB 数据源配置
type MultiDB []*DB

type DB struct {
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

// Title 配置器标题
func (MultiDB) Title() string {
	return "DB"
}

// Reader 配置文件读取
func (MultiDB) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "database.yaml",
		NacosDataId: "database.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (m MultiDB) Run() error {
	if len(m) == 0 {
		log.Error("database not connected! reason: [database.yaml] Not Found")
		return nil
	}
	if handler == nil {
		handler = &Handler{
			multi:     true,
			gormMap:   make(map[string]*gorm.DB),
			configMap: make(map[string]*DB),
		}
	} else {
		handler.multi = true
	}
	for i, db := range m {
		if db.Enable {
			if err := anyx.SetDefaultValue(db); err != nil {
				return errorx.Wrap(err, "set-default-value error")
			}
			gormDB, err := db.NewGormDB()
			if err != nil {
				log.Error("database connect failed: ", db.Info(), err)
				return errorx.Wrap(err, "NewGormDB Failed")
			}
			handler.gormMap[db.Source] = gormDB
			handler.configMap[db.Source] = db
			if i == 0 || db.Source == constx.DefaultKey {
				handler.gormDB = gormDB
				handler.config = db
			}
			log.Info("database connect successful: ", db.Info())
		}
	}
	if len(handler.configMap) == 0 {
		log.Error("database not connected! reason: [database.yaml] is empty or no enabled database configured")
	}
	return nil
}

// Info 配置信息格式化
func (d *DB) Info() string {
	return fmtx.Green.XSPrintf("source=%s type=%s host=%s port=%v database=%s debug=%v",
		d.Source, d.Type, d.Host, d.Port, d.Database, d.Debug)
}

// Title 配置器标题
func (d *DB) Title() string {
	return "DB"
}

// Reader 配置文件读取
func (d *DB) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "database.yaml",
		NacosDataId: "database.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (d *DB) Run() error {
	if d.Enable {
		if err := anyx.SetDefaultValue(d); err != nil {
			return errorx.Wrap(err, "set-default-value error")
		}
		gormDB, err := d.NewGormDB()
		if err != nil {
			log.Error("database connect failed: ", d.Info(), err)
			return errorx.Wrap(err, "NewGormDB Failed")
		}

		if handler == nil {
			handler = &Handler{
				multi:     false,
				config:    d,
				configMap: make(map[string]*DB),
				gormDB:    gormDB,
				gormMap:   map[string]*gorm.DB{},
			}
		} else {
			handler.multi = true
		}
		handler.gormMap[d.Source] = gormDB
		handler.configMap[d.Source] = d

		log.Info("database connect successful: ", d.Info())
		return nil
	}
	log.Info("database not connected! reason: database.yaml is empty or the value of enable is false")
	return nil
}

// NewGormDB 创建数据库连接
func (d *DB) NewGormDB() (gormDB *gorm.DB, err error) {
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

// CommentTableSql 生成表备注
func (d *DB) CommentTableSql(table, comment string) string {
	switch strings.ToLower(d.Type) {
	case Mysql:
		return "alter table " + table + " comment = '" + comment + "'"
	case Postgres:
		return "comment on table " + table + " is '" + comment + "'"
	}
	return ""
}

// GetGormDB 根据dsn生成gormDB
func (d *DB) GetGormDB() (gormDb *gorm.DB, err error) {
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
