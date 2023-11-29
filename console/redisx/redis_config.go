package redisx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	StandAlone = iota // 单机
	Cluster           // 集群
)

// redis连接配置
type Configs []*Config
type Config struct {
	Source   string `json:"source" yaml:"source"`     // 数据源名称
	Mode     int    `json:"mode" yaml:"mode"`         // 模式（0-单机；1-集群），默认单机模式
	Host     string `json:"host" yaml:"host"`         // 主机
	Port     int    `json:"port" yaml:"port"`         // 端口
	Password string `json:"password" yaml:"password"` // 密码
	Database int    `json:"database" yaml:"database"` // 数据库，默认0
	PoolSize int    `json:"poolSize" yaml:"poolSize"` // 池大小
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf(
		"source=%s mode=%d host=%s port=%d database=%d",
		conf.Source, conf.Mode, conf.Host, conf.Port, conf.Database)
}

// 配置信息格式化
func (conf *Config) Address() string {
	return fmt.Sprintf("%s:%d", conf.Host, conf.Port)
}

// 初始化redis，默认单机模式
func (conf *Config) NewRedisCmdable(database ...int) (cmd redis.Cmdable) {
	var db = conf.Database
	if len(database) > 0 {
		db = database[0]
	}
	switch conf.Mode {
	case StandAlone:
		cmd = redis.NewClient(&redis.Options{
			Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
			Password: conf.Password,
			PoolSize: conf.PoolSize,
			DB:       db,
		})
	case Cluster:
		cmd = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(conf.Host, ","),
			Password: conf.Password,
			PoolSize: conf.PoolSize,
		})
	default:
		log.Warn("redis模式配置错误！", conf)
		return
	}
	return
}
