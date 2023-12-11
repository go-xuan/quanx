package redisx

import (
	"encoding/json"
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

func (configs Configs) Format() (f string) {
	bytes, _ := json.Marshal(configs)
	return string(bytes)
}

func (configs Configs) Init() {
	InitRedisX(configs)
}

// 初始化redisX
func InitRedisX(configs Configs) {
	if len(configs) == 0 {
		log.Error("redis配置为空!")
		return
	}
	instance = &Handler{
		CmdMap:    make(map[string]redis.Cmdable),
		ConfigMap: make(map[string]*Config),
	}
	for i, conf := range configs {
		var cmd = conf.NewRedisCmdable()
		if ok, err := Ping(cmd); ok && err == nil {
			log.Info("redis连接成功! ", conf.Format())
			if i == 0 || conf.Source == "default" {
				instance.Cmd = cmd
				instance.Config = conf
			}
			instance.CmdMap[conf.Source] = cmd
			instance.ConfigMap[conf.Source] = conf
		} else {
			log.Error("redis连接失败! ", conf.Format())
			log.Error("error : ", err)
		}
	}
}

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
func (config *Config) Format() string {
	return fmt.Sprintf("source=%s mode=%d host=%s port=%d database=%d",
		config.Source, config.Mode, config.Host, config.Port, config.Database)
}

func (config *Config) Init() {
	var cmd = config.NewRedisCmdable()
	if ok, err := Ping(cmd); ok && err == nil {
		log.Info("redis连接成功! ", config.Format())
		instance.Cmd = cmd
		instance.Config = config
		instance.CmdMap[config.Source] = cmd
		instance.ConfigMap[config.Source] = config
	} else {
		log.Error("redis连接失败! ", config.Format())
		log.Error("error : ", err)
	}
}

// 配置信息格式化
func (config *Config) Address() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// 初始化redis，默认单机模式
func (config *Config) NewRedisCmdable(database ...int) (cmd redis.Cmdable) {
	var db = config.Database
	if len(database) > 0 {
		db = database[0]
	}
	switch config.Mode {
	case StandAlone:
		cmd = redis.NewClient(&redis.Options{
			Addr:     config.Host + ":" + strconv.Itoa(config.Port),
			Password: config.Password,
			PoolSize: config.PoolSize,
			DB:       db,
		})
	case Cluster:
		cmd = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(config.Host, ","),
			Password: config.Password,
			PoolSize: config.PoolSize,
		})
	default:
		log.Warn("redis模式配置错误！", config)
		return
	}
	return
}
