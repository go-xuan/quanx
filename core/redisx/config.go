package redisx

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
)

const (
	StandAlone = iota // 单机
	Cluster           // 集群
	Sentinel          // 集群
)

type Config struct {
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

func (c *Config) Format() string {
	return fmt.Sprintf("source=%s mode=%v host=%s port=%v database=%v",
		c.Source, c.Mode, c.Host, c.Port, c.Database)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		var client = c.NewRedisClient()
		if result, err := client.Ping(context.TODO()).Result(); err != nil || result != "PONG" {
			return errorx.Wrap(err, "redis client ping error")
		}
		if _handler == nil {
			_handler = &Handler{
				multi:   false,
				config:  c,
				configs: make(map[string]*Config),
				client:  client,
				clients: make(map[string]redis.UniversalClient),
			}
		} else {
			_handler.multi = true
		}
		_handler.clients[c.Source] = client
		_handler.configs[c.Source] = c
		return nil
	}
	log.Error(`redis connect failed! reason: redis.yaml is empty or the value of "enable" is false`)
	return nil
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// NewRedisClient 初始化redis客户端
// UniversalClient并不是一个客户端，而是对单节点客户端/集群客户端/哨兵客户端的接口包装。根据不同的选项，客户端的类型如下：
// 1、如果指定了MasterName选项，则返回redis.FailoverClient哨兵客户端。
// 2、如果Addrs是2个以上的地址，则返回redis.ClusterClient集群客户端。
// 3、其他情况，返回redis.Client单节点客户端。
func (c *Config) NewRedisClient() redis.UniversalClient {
	var opts = &redis.UniversalOptions{
		ClientName: c.Source,
		Username:   c.Username,
		Password:   c.Password,
		PoolSize:   c.PoolSize,
		DB:         c.Database,
	}
	switch c.Mode {
	case StandAlone:
		opts.Addrs = []string{net.JoinHostPort(c.Host, strconv.Itoa(c.Port))}
		return redis.NewClient(opts.Simple())
	case Cluster:
		opts.Addrs = strings.Split(c.Host, ",")
		return redis.NewClusterClient(opts.Cluster())
	case Sentinel:
		opts.Addrs = []string{net.JoinHostPort(c.Host, strconv.Itoa(c.Port))}
		opts.MasterName = c.MasterName
		return redis.NewFailoverClient(opts.Failover())
	default:
		log.Warn("redis mode is invalid: ")
		return nil
	}
}

// MultiConfig redis多连接配置
type MultiConfig []*Config

func (m MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range m {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("{")
		sb.WriteString(config.Format())
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}

func (MultiConfig) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

func (m MultiConfig) Execute() error {
	if len(m) == 0 {
		return errorx.New("redis not connected! cause: redis.yaml is invalid")
	}
	if _handler == nil {
		_handler = &Handler{
			multi:   true,
			configs: make(map[string]*Config),
			clients: make(map[string]redis.UniversalClient),
		}
	} else {
		_handler.multi = true
	}
	var ctx = context.Background()
	for i, c := range m {
		if c.Enable {
			if err := anyx.SetDefaultValue(c); err != nil {
				return errorx.Wrap(err, "set default value error")
			}
			var client = c.NewRedisClient()
			if result, err := client.Ping(ctx).Result(); err != nil || result != "PONG" {
				return errorx.Wrap(err, "redis client ping error")
			}
			_handler.clients[c.Source] = client
			_handler.configs[c.Source] = c
			if i == 0 || c.Source == constx.DefaultSource {
				_handler.client = client
				_handler.config = c
			}
		}
	}
	if len(_handler.configs) == 0 {
		log.Error("redis connect failed! cause: redis.yaml is empty or no enabled redis configured")
	}
	return nil
}
