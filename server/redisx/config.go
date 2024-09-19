package redisx

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/app/constx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

const (
	StandAlone = iota // 单机
	Cluster           // 集群
	Sentinel          // 集群
)

// MultiRedis redis连接配置
type MultiRedis []*Redis

type Redis struct {
	Source     string `json:"source" yaml:"source" default:"default"` // 数据源名称
	Enable     bool   `json:"enable" yaml:"enable"`                   // 数据源启用
	Mode       int    `json:"mode" yaml:"mode" default:"0"`           // 模式（0-单机；1-集群；3-哨兵。默认单机模式）
	Host       string `json:"host" yaml:"host"`                       // 主机（单机模式使用）
	Port       int    `json:"port" yaml:"port" default:"6379"`        // 端口
	Username   string `json:"username" yaml:"username"`               // 用户名
	Password   string `json:"password" yaml:"password"`               // 密码
	Database   int    `json:"database" yaml:"database" default:"0"`   // 数据库，默认0
	MasterName string `json:"masterName" yaml:"masterName"`           // 哨兵模式主服务器名称
	PoolSize   int    `json:"poolSize" yaml:"poolSize"`               // 池大小
}

func (MultiRedis) ID() string {
	return "multi-redis"
}

func (MultiRedis) Format() string {
	return ""
}

func (MultiRedis) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

func (conf MultiRedis) Execute() error {
	if len(conf) == 0 {
		log.Error("redis not connected! cause: redis.yaml not found")
		return nil
	}
	if handler == nil {
		handler = &Handler{
			multi:     true,
			configMap: make(map[string]*Redis),
			clientMap: make(map[string]redis.UniversalClient),
		}
	} else {
		handler.multi = true
	}
	for i, r := range conf {
		if r.Enable {
			if err := anyx.SetDefaultValue(r); err != nil {
				return errorx.Wrap(err, "set default value error")
			}
			var client = r.NewRedisClient()
			if ok, err := Ping(client); !ok && err != nil {
				return errorx.Wrap(err, "redis client ping error")
			}
			handler.clientMap[r.Source] = client
			handler.configMap[r.Source] = r
			if i == 0 || r.Source == constx.DefaultKey {
				handler.client = client
				handler.config = r
			}
		}
	}
	if len(handler.configMap) == 0 {
		log.Error("redis connect failed! cause: redis.yaml is empty or no enabled redis configured")
	}
	return nil
}

func (r *Redis) ID() string {
	return "redis"
}

func (r *Redis) Format() string {
	return fmtx.Yellow.XSPrintf("source=%s mode=%v host=%s port=%v database=%v",
		r.Source, r.Mode, r.Host, r.Port, r.Database)
}

func (*Redis) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

func (r *Redis) Execute() error {
	if r.Enable {
		if err := anyx.SetDefaultValue(r); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		var client = r.NewRedisClient()
		if ok, err := Ping(client); !ok && err != nil {
			return errorx.Wrap(err, "redis client ping error")
		}
		if handler == nil {
			handler = &Handler{
				multi:     false,
				config:    r,
				configMap: make(map[string]*Redis),
				client:    client,
				clientMap: make(map[string]redis.UniversalClient),
			}
		} else {
			handler.multi = true
		}
		handler.clientMap[r.Source] = client
		handler.configMap[r.Source] = r
		return nil
	}
	log.Error(`redis connect failed! reason: redis.yaml is empty or the value of "enable" is false`)
	return nil
}

func (r *Redis) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// NewRedisClient 初始化redis客户端
// UniversalClient并不是一个客户端，而是对单节点客户端/集群客户端/哨兵客户端的接口包装。根据不同的选项，客户端的类型如下：
// 1、如果指定了MasterName选项，则返回redis.FailoverClient哨兵客户端。
// 2、如果Addrs是2个以上的地址，则返回redis.ClusterClient集群客户端。
// 3、其他情况，返回redis.Client单节点客户端。
func (r *Redis) NewRedisClient() redis.UniversalClient {
	var opts = &redis.UniversalOptions{
		ClientName: r.Source,
		Username:   r.Username,
		Password:   r.Password,
		PoolSize:   r.PoolSize,
	}
	switch r.Mode {
	case StandAlone:
		opts.Addrs = []string{net.JoinHostPort(r.Host, strconv.Itoa(r.Port))}
		return redis.NewClient(opts.Simple())
	case Cluster:
		opts.Addrs = strings.Split(r.Host, ",")
		return redis.NewClusterClient(opts.Cluster())
	case Sentinel:
		opts.Addrs = []string{net.JoinHostPort(r.Host, strconv.Itoa(r.Port))}
		opts.MasterName = r.MasterName
		return redis.NewFailoverClient(opts.Failover())
	default:
		log.Warn("redis mode is invalid: ")
		return nil
	}
}
