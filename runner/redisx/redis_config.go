package redisx

import (
	"fmt"
	"github.com/go-xuan/quanx/runner/nacosx"
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
type MultiRedis []*Redis

// 配置信息格式化
func (r MultiRedis) Name() string {
	return "连接Redis"
}

// nacos配置文件
func (MultiRedis) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "multi_redis.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (MultiRedis) LocalConfig() string {
	return "conf/multi_redis.yaml"
}

func (r MultiRedis) Run() error {
	if len(r) == 0 {
		log.Info("未配置多Redis数据库! ")
		return nil
	}
	handler = &Handler{
		CmdMap:    make(map[string]redis.Cmdable),
		ConfigMap: make(map[string]*Redis),
	}
	for i, conf := range r {
		if conf.Enable {
			var cmd = conf.NewRedisCmdable()
			ok, err := Ping(cmd)
			if !ok && err != nil {
				log.Error("redis连接失败! ", conf.ToString())
				log.Error("error : ", err)
				return err
			}
			handler.CmdMap[conf.Source] = cmd
			handler.ConfigMap[conf.Source] = conf
			if i == 0 || conf.Source == "default" {
				handler.Cmd = cmd
				handler.Config = conf
			}
			log.Info("redis连接成功! ", conf.ToString())
		}
	}
	return nil
}

type Redis struct {
	Source   string `json:"source" yaml:"source"`     // 数据源名称
	Enable   bool   `json:"enable" yaml:"enable"`     // 数据源启用
	Mode     int    `json:"mode" yaml:"mode"`         // 模式（0-单机；1-集群），默认单机模式
	Host     string `json:"host" yaml:"host"`         // 主机
	Port     int    `json:"port" yaml:"port"`         // 端口
	Password string `json:"password" yaml:"password"` // 密码
	Database int    `json:"database" yaml:"database"` // 数据库，默认0
	PoolSize int    `json:"poolSize" yaml:"poolSize"` // 池大小
}

// 配置信息格式化
func (r *Redis) ToString() string {
	return fmt.Sprintf("source=%s mode=%d host=%s port=%d database=%d",
		r.Source, r.Mode, r.Host, r.Port, r.Database)
}

// 运行器名称
func (r *Redis) Name() string {
	return "连接Redis"
}

// nacos配置文件
func (*Redis) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "redis.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (*Redis) LocalConfig() string {
	return "conf/redis.yaml"
}

// 运行器运行
func (r *Redis) Run() error {
	if r.Enable {
		var cmd = r.NewRedisCmdable()
		if ok, err := Ping(cmd); !ok && err != nil {
			log.Error("redis连接失败! ", r.ToString())
			log.Error("error : ", err)
			return err
		}
		handler = &Handler{
			Cmd:       cmd,
			Config:    r,
			CmdMap:    make(map[string]redis.Cmdable),
			ConfigMap: make(map[string]*Redis),
		}
		handler.CmdMap[r.Source] = cmd
		handler.ConfigMap[r.Source] = r
		log.Info("redis连接成功! ", r.ToString())
	}
	return nil
}

// 配置信息格式化
func (r *Redis) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// 初始化redis，默认单机模式
func (r *Redis) NewRedisCmdable(database ...int) (cmd redis.Cmdable) {
	var db = r.Database
	if len(database) > 0 {
		db = database[0]
	}
	switch r.Mode {
	case StandAlone:
		cmd = redis.NewClient(&redis.Options{
			Addr:     r.Host + ":" + strconv.Itoa(r.Port),
			Password: r.Password,
			PoolSize: r.PoolSize,
			DB:       db,
		})
	case Cluster:
		cmd = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(r.Host, ","),
			Password: r.Password,
			PoolSize: r.PoolSize,
		})
	default:
		log.Warn("redis模式配置错误！", r)
		return
	}
	return
}
