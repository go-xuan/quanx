package redisx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
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

func (c *Config) Copy() *Config {
	return &Config{
		Source:     c.Source,
		Enable:     c.Enable,
		Mode:       c.Mode,
		Host:       c.Host,
		Port:       c.Port,
		Username:   c.Username,
		Password:   c.Password,
		Database:   c.Database,
		MasterName: c.MasterName,
		PoolSize:   c.PoolSize,
	}
}

func (c *Config) Format() string {
	return fmt.Sprintf("source=%s host=%s port=%d database=%d mode=%d",
		c.Source, c.Host, c.Port, c.Database, c.Mode)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "redis.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalReader{
			Name: "redis.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if client, err := c.NewRedisClient(); err != nil {
			log.Error("redis connect failed:", c.Format())
			return errorx.Wrap(err, "redis init error")
		} else {
			AddClient(c, client)
			log.Info("redis connect success: ", c.Format())
		}
	}
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
func (c *Config) NewRedisClient() (redis.UniversalClient, error) {
	var opts = &redis.UniversalOptions{
		ClientName: c.Source,
		Username:   c.Username,
		Password:   c.Password,
		PoolSize:   c.PoolSize,
		DB:         c.Database,
	}
	var client redis.UniversalClient
	switch c.Mode {
	case StandAlone:
		opts.Addrs = []string{net.JoinHostPort(c.Host, strconv.Itoa(c.Port))}
		client = redis.NewClient(opts.Simple())
	case Cluster:
		opts.Addrs = strings.Split(c.Host, ",")
		client = redis.NewClusterClient(opts.Cluster())
	case Sentinel:
		opts.Addrs = []string{net.JoinHostPort(c.Host, strconv.Itoa(c.Port))}
		opts.MasterName = c.MasterName
		client = redis.NewFailoverClient(opts.Failover())
	default:
		log.Warn("redis mode is invalid: ")
		return nil, errors.New("redis mode is invalid")
	}
	if result, err := client.Ping(context.TODO()).Result(); err != nil || result != "PONG" {
		log.Error("client ping failed:", c.Format())
		return client, errorx.Wrap(err, "client ping error")
	}
	return client, nil
}

// MultiConfig redis多连接配置
type MultiConfig []*Config

func (list MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range list {
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

func (MultiConfig) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "redis.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalReader{
			Name: "redis.yaml",
		}
	default:
		return nil
	}
}

func (list MultiConfig) Execute() error {
	if len(list) == 0 {
		return errorx.New("redis not initialized! cause: redis.yaml is invalid")
	}
	for _, config := range list {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "redis config execute error")
		}
	}
	if !Initialized() {
		log.Error("redis not initialized! cause: no enabled source")
	}
	return nil
}
