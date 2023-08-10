package redisx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// redis连接配置
type Config struct {
	Mode     int    `json:"mode" yaml:"mode" nacos:"redis.mode"`             // 模式（0-单机；1-集群），默认单机模式
	Host     string `json:"host" yaml:"host" nacos:"redis.host"`             // 主机
	Port     int    `json:"port" yaml:"port" nacos:"redis.port"`             // 端口
	Password string `json:"password" yaml:"password" nacos:"redis.password"` // 密码
	Database int    `json:"database" yaml:"database" nacos:"redis.database"` // 数据库，默认0
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d Database=%d", conf.Host, conf.Port, conf.Database)
}

const (
	StandAlone = iota // 单机
	Cluster           // 集群
)

// 初始化redis，默认单机模式
func (conf *Config) NewRedisCmdable() (cmd redis.Cmdable) {
	switch conf.Mode {
	case StandAlone:
		cmd = redis.NewClient(&redis.Options{
			Addr:     conf.Host + ":" + strconv.Itoa(conf.Port),
			Password: conf.Password,
			DB:       conf.Database,
		})
	case Cluster:
		cmd = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(conf.Host, ","),
			Password: conf.Password,
		})
	default:
		log.Warn("redis模式配置错误！", conf)
		return
	}
	return
}
