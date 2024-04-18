package redisx

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/frame/confx"
	"github.com/go-xuan/quanx/utils/anyx"
)

const (
	StandAlone = iota // 单机
	Cluster           // 集群
)

// redis连接配置
type MultiRedis []*Redis

type Redis struct {
	Source   string `json:"source" yaml:"source"`     // 数据源名称
	Enable   bool   `json:"enable" yaml:"enable"`     // 数据源启用
	Mode     int    `json:"mode" yaml:"mode"`         // 模式（0-单机；1-集群），默认单机模式
	Host     string `json:"host" yaml:"host"`         // 主机
	Port     int    `json:"port" yaml:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database int    `json:"database" yaml:"database"` // 数据库，默认0
	PoolSize int    `json:"poolSize" yaml:"poolSize"` // 池大小
	SSH      SSH    `json:"ssh" yaml:"ssh"`           // ssh
}

type SSH struct {
	Host     string `json:"host" yaml:"host"`         // 主机
	Port     int    `json:"port" yaml:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
}

// 配置信息格式化
func (MultiRedis) Theme() string {
	return "Redis"
}

// 配置文件读取
func (MultiRedis) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (conf MultiRedis) Run() error {
	if len(conf) == 0 {
		log.Error("redis connect failed! reason: redis.yaml not found!")
		return nil
	}
	handler = &Handler{
		Multi:     true,
		clientMap: make(map[string]*redis.UniversalClient),
		ConfigMap: make(map[string]*Redis),
	}
	for i, r := range conf {
		if r.Enable {
			var client = r.NewRedisClient()
			if ok, err := Ping(*client); !ok && err != nil {
				log.Error(r.ToString("redis connect failed!"))
				log.Error("error : ", err)
				return err
			}
			handler.clientMap[r.Source] = client
			handler.ConfigMap[r.Source] = r
			if i == 0 || r.Source == constx.Default {
				handler.Client = client
				handler.Config = r
			}
			log.Info(r.ToString("redis connect successful!"))
		}
	}
	if len(handler.ConfigMap) == 0 {
		log.Error("redis connect failed! reason: redis.yaml is empty or all enable values are false!")
	}
	return nil
}

// 配置信息格式化
func (r *Redis) ToString(title string) string {
	return fmt.Sprintf("%s => source=%s mode=%d host=%s port=%d database=%d",
		title, r.Source, r.Mode, r.Host, r.Port, r.Database)
}

// 配置器名称
func (r *Redis) Theme() string {
	return "Redis"
}

// 配置文件读取
func (*Redis) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "redis.yaml",
		NacosDataId: "redis.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (r *Redis) Run() (err error) {
	if r.Enable {
		if err = anyx.SetDefaultValue(r); err != nil {
			return
		}
		var client = r.NewRedisClient()
		var ok bool
		if ok, err = Ping(*client); !ok && err != nil {
			log.Error(r.ToString("Redis connect failed"))
			log.Error("error : ", err)
			return
		}
		handler = &Handler{
			Multi:     false,
			Client:    client,
			Config:    r,
			clientMap: make(map[string]*redis.UniversalClient),
			ConfigMap: make(map[string]*Redis),
		}
		handler.clientMap[r.Source] = client
		handler.ConfigMap[r.Source] = r
		log.Info(r.ToString("Redis connect successful"))
		return
	}
	log.Error(`Redis connect failed! reason: redis.yaml is empty or the value of "enable" is false`)
	return
}

// 配置信息格式化
func (r *Redis) Address() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// 初始化redis客户端
// UniversalClient并不是一个客户端，而是对 Client、ClusterClient、FailoverClient 客户端的包装。根据不同的选项，客户端的类型如下：
// 1、如果指定了MasterName选项，则返回FailoverClient哨兵客户端。
// 2、如果Addrs是2个以上的地址，则返回ClusterClient集群客户端。
// 3、其他情况，返回Client单节点客户端。
func (r *Redis) NewRedisClient() *redis.UniversalClient {
	var client redis.UniversalClient
	switch r.Mode {
	case StandAlone:
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    []string{net.JoinHostPort(r.Host, strconv.Itoa(r.Port))},
			Username: r.Username,
			Password: r.Password,
			PoolSize: r.PoolSize,
			DB:       r.Database,
		})
	case Cluster:
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:    strings.Split(r.Host, ","),
			Username: r.Username,
			Password: r.Password,
			PoolSize: r.PoolSize,
		})
	default:
		log.Warn(r.ToString("the mode value can only be 1 or 2"))
		return nil
	}
	return &client
}
