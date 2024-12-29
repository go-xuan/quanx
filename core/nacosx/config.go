package nacosx

import (
	"path/filepath"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	OnlyConfig      = iota // 仅用配置中心
	OnlyNaming             // 仅用服务发现
	ConfigAndNaming        // 配置中心和服务发现都使用
)

// Config nacos连接配置
type Config struct {
	Address   string `yaml:"address" json:"address" default:"127.0.0.1"`  // nacos服务地址,多个以英文逗号分割
	Username  string `yaml:"username" json:"username" default:"nacos"`    // 用户名
	Password  string `yaml:"password" json:"password" default:"nacos"`    // 密码
	NameSpace string `yaml:"nameSpace" json:"nameSpace" default:"public"` // 命名空间
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                // 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
}

func (c *Config) Format() string {
	return fmtx.Yellow.XSPrintf("address=%s username=%s password=%s nameSpace=%s mode=%v",
		c.AddressUrl(), c.Username, c.Password, c.NameSpace, c.Mode)
}

func (*Config) Reader() *configx.Reader {
	return nil
}

func (c *Config) Execute() (err error) {
	if handler == nil {
		handler = &Handler{config: c}
		switch c.Mode {
		case OnlyConfig:
			if handler.configClient, err = c.ConfigClient(c.ClientParam()); err != nil {
				return
			}
		case OnlyNaming:
			if handler.namingClient, err = c.NamingClient(c.ClientParam()); err != nil {
				return
			}
		case ConfigAndNaming:
			var param = c.ClientParam()
			if handler.configClient, err = c.ConfigClient(param); err != nil {
				return
			}
			if handler.namingClient, err = c.NamingClient(param); err != nil {
				return
			}
		}
	}
	log.Info("nacos connect successfully: ", c.Format())
	return
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
		NamespaceId:         c.NameSpace,
		LogDir:              filepath.Join(constx.DefaultResourceDir, ".nacos/log"),
		CacheDir:            filepath.Join(constx.DefaultResourceDir, ".nacos/cache"),
	}
}

// ServerConfigs nacos服务中间件配置
func (c *Config) ServerConfigs() []constant.ServerConfig {
	var adds = strings.Split(c.Address, ",")
	if len(adds) == 0 {
		log.Error("the address of nacos cannot be empty")
		return nil
	}
	var configs []constant.ServerConfig
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
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
func (c *Config) ConfigClient(param vo.NacosClientParam) (client config_client.IConfigClient, err error) {
	if client, err = clients.NewConfigClient(param); err != nil {
		log.Error("nacos config client init failed: ", c.Format())
		log.Error(err)
		return
	}
	return
}

// NamingClient 初始化Nacos服务发现客户端
func (c *Config) NamingClient(param vo.NacosClientParam) (client naming_client.INamingClient, err error) {
	if client, err = clients.NewNamingClient(param); err != nil {
		log.Error("nacos naming client init failed: ", c.Format(), err)
		return
	}
	return
}
