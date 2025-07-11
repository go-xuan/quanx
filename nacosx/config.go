package nacosx

import (
	"path/filepath"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
)

const (
	OnlyConfig      = iota // 仅启用配置中心
	OnlyNaming             // 仅启用服务发现
	ConfigAndNaming        // 启用配置中心和服务发现
)

// Config nacos连接配置
type Config struct {
	Address   string `yaml:"address" json:"address" default:"localhost:8848"` // nacos服务地址,多个以英文逗号分割
	Username  string `yaml:"username" json:"username"`                        // 用户名
	Password  string `yaml:"password" json:"password"`                        // 密码
	Namespace string `yaml:"namespace" json:"namespace" default:"public"`     // 命名空间
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                    // 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
}

func (c *Config) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"address":   c.AddressUrl(),
		"username":  c.Username,
		"password":  c.Password,
		"namespace": c.Namespace,
		"mode":      c.Mode,
	})
}

func (c *Config) NeedRead() bool {
	if c.Address == "" && c.Username == "" {
		return true
	}
	return false
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FromEnv:
		return &configx.EnvReader{}
	case configx.FromFile:
		return &configx.FileReader{
			Name: "nacos.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if _client == nil {
		_client = &Client{config: c}
		var param = c.ClientParam()
		switch c.Mode {
		case OnlyConfig: // 仅初始化配置中心
			if configClient, err := c.ConfigClient(param); err != nil {
				return errorx.Wrap(err, "init nacos config client error")
			} else {
				_client.configClient = configClient
			}
		case OnlyNaming: // 仅初始化服务发现
			if namingClient, err := c.NamingClient(param); err != nil {
				return errorx.Wrap(err, "init nacos naming client error")
			} else {
				_client.namingClient = namingClient
			}
		case ConfigAndNaming: // 初始化配置中心和服务发现
			if configClient, err := c.ConfigClient(param); err != nil {
				return errorx.Wrap(err, "init nacos config client error")
			} else {
				_client.configClient = configClient
			}
			if namingClient, err := c.NamingClient(param); err != nil {
				return errorx.Wrap(err, "init nacos naming client error")
			} else {
				_client.namingClient = namingClient
			}
		}
	}
	c.LogEntry().Info("nacos init success")
	return nil
}

// AddressUrl nacos访问地址
func (c *Config) AddressUrl() string {
	return c.Address + "/nacos"
}

// EnableNaming 开启服务注册
func (c *Config) EnableNaming() bool {
	return c.Mode == OnlyNaming || c.Mode == ConfigAndNaming
}

// ClientConfig nacos客户端配置
func (c *Config) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		Username:            c.Username,
		Password:            c.Password,
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         c.Namespace,
		LogDir:              filepath.Join(constx.DefaultResourceDir, ".nacos/log"),
		CacheDir:            filepath.Join(constx.DefaultResourceDir, ".nacos/cache"),
	}
}

// ServerConfigs nacos服务中间件配置
func (c *Config) ServerConfigs() []constant.ServerConfig {
	var addrs = strings.Split(c.Address, ",")
	if len(addrs) == 0 {
		log.Error("the address of nacos cannot be empty")
		return nil
	}
	var configs []constant.ServerConfig
	for _, addr := range addrs {
		host, port, _ := strings.Cut(addr, ":")
		configs = append(configs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      host,
			Port:        uint64(stringx.ParseInt64(port)),
		})
	}
	return configs
}

func (c *Config) ClientParam() vo.NacosClientParam {
	return vo.NacosClientParam{
		ClientConfig:  c.ClientConfig(),
		ServerConfigs: c.ServerConfigs(),
	}
}

// ConfigClient 初始化Nacos配置中心客户端
func (c *Config) ConfigClient(param vo.NacosClientParam) (config_client.IConfigClient, error) {
	if client, err := clients.NewConfigClient(param); err != nil {
		c.LogEntry().WithError(err).Error("nacos config client init failed")
		return nil, errorx.Wrap(err, "nacos config client init failed")
	} else {
		return client, nil
	}
}

// NamingClient 初始化Nacos服务发现客户端
func (c *Config) NamingClient(param vo.NacosClientParam) (naming_client.INamingClient, error) {
	if client, err := clients.NewNamingClient(param); err != nil {
		c.LogEntry().WithError(err).Error("nacos naming client init failed")
		return nil, errorx.Wrap(err, "nacos naming client init failed")
	} else {
		return client, nil
	}
}
