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

	"github.com/go-xuan/quanx/runner/nacosx"
)

// 数据源配置
type MultiDatabase []*Database

// 运行器名称
func (d MultiDatabase) Name() string {
	return "init multi-database"
}

// nacos配置文件
func (MultiDatabase) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "multi_database.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (MultiDatabase) LocalConfig() string {
	return "conf/multi_database.yaml"
}

// 运行器运行
func (d MultiDatabase) Run() error {
	if len(d) == 0 {
		log.Info("database not connected! reason: multi-database.yaml not found!")
		return nil
	}
	handler = &Handler{
		Multi:     true,
		DBMap:     make(map[string]*gorm.DB),
		ConfigMap: make(map[string]*Database),
	}
	for i, conf := range d {
		if conf.Enable {
			db, err := conf.NewGormDB()
			if err != nil {
				log.Error(conf.ToString("database connect failed!"))
				return err
			}
			handler.DBMap[conf.Source] = db
			handler.ConfigMap[conf.Source] = conf
			if i == 0 || conf.Source == "default" {
				handler.DB = db
				handler.Config = conf
			}
			log.Info(conf.ToString("database connect successful!"))
		}
	}
	if len(handler.ConfigMap) == 0 {
		log.Info("database not connected! reason: multi-database.yaml is empty or {database.enable} is false")
	}
	return nil
}

// 根据数据源名称获取配置
func (d MultiDatabase) GetConfig(source string) *Database {
	for _, config := range d {
		if config.Source == source {
			return config
		}
	}
	return nil
}

// 获取默认配置
func (d MultiDatabase) GetDefault() *Database {
	if len(d) > 0 {
		var conf = d.GetConfig("default")
		if conf == nil {
			conf = d[0]
		}
		return conf
	}
	return nil
}

type Database struct {
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
func (d *Database) ToString(title string) string {
	return fmt.Sprintf("%s => source=%s type=%s host=%s port=%d database=%s debug=%v",
		title, d.Source, d.Type, d.Host, d.Port, d.Database, d.Debug)
}

// 运行器名称
func (d *Database) Name() string {
	return "inti database"
}

// nacos配置文件
func (*Database) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "database.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (*Database) LocalConfig() string {
	return "conf/database.yaml"
}

// 运行器运行
func (d *Database) Run() error {
	if d.Enable {
		d.Source = "default"
		db, err := d.NewGormDB()
		if err != nil {
			log.Error(d.ToString("database connect failed!"))
			log.Error("error : ", err)
			return err
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
	}
	log.Info("database not connected! reason: database.yaml is empty or {database.enable} is false")
	return nil
}

// 创建数据库连接
func (d *Database) NewGormDB() (gormDB *gorm.DB, err error) {
	gormDB, err = GetGormDB(d.GetDSN(), d.Type)
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

// 获取数据库连接DSN
func (d *Database) GetDSN() (dsn string) {
	switch strings.ToLower(d.Type) {
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			d.Username, d.Password, d.Host, d.Port, d.Database) +
			"?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local"
	case Postgres:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			d.Host, d.Port, d.Username, d.Password, d.Database)
	}
	return
}

// 名称
func (d *Database) CommentTableSql(table, comment string) string {
	switch strings.ToLower(d.Type) {
	case Mysql:
		return fmt.Sprintf(`ALTER TABLE %s COMMENT = '%s'`, table, comment)
	case Postgres:
		return fmt.Sprintf(`COMMENT ON TABLE %s IS '%s'`, table, comment)
	}
	return ""
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
