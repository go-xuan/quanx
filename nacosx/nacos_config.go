package nacosx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

const (
	ConfigAndNaming = iota // 配置中心和服务发现都使用，默认项
	OnlyConfig             // 仅用配置中心
	OnlyNaming             // 仅用服务发现
)

// nacos访问配置
type Config struct {
	Address    string      `yaml:"address" json:"address"`                          // nacos服务地址,多个以英文逗号分割
	Path       string      `yaml:"path" json:"path" default:"/nacos"`               // web访问路径
	Username   string      `yaml:"username" json:"username" default:"nacos"`        // 用户名
	Password   string      `yaml:"password" json:"password" default:"nacos"`        // 密码
	NameSpace  string      `yaml:"nameSpace" json:"nameSpace" default:"public"`     // 命名空间
	Mode       int         `yaml:"mode" json:"mode" default:"0"`                    // 模式（0-配置中心和服务发现都用；1-仅用配置中心；2-仅用服务发现）
	LogDir     string      `yaml:"logDir" json:"logDir" default:".nacos/log"`       // 日志文件夹
	CacheDir   string      `yaml:"cacheDir" json:"cacheDir" default:".nacos/cache"` // 缓存文件夹
	LoadConfig *LoadConfig `yaml:"loadConfig" json:"loadConfig"`                    // 加载配置项
}

// nacos加载配置项
type LoadConfig struct {
	Basic  string `yaml:"basic"`  // basic配置项
	Custom string `yaml:"custom"` // 自定义配置项
	Listen string `yaml:"listen"` // 监听配置项
}

func (config *Config) Format() string {
	return fmt.Sprintf("address=%s path=%s username=%s password=%s nameSpace=%s mode=%d",
		config.Address, config.Path, config.Username, config.Password, config.NameSpace, config.Mode)
}

func (config *Config) Init() {
	InitNacosX(config)
}

func InitNacosX(conf *Config) {
	instance = &Handler{Config: conf}
	switch conf.Mode {
	case OnlyConfig:
		instance.ConfigClient = conf.ConfigClient()
	case OnlyNaming:
		instance.NamingClient = conf.NamingClient()
	default:
		instance.ConfigClient = conf.ConfigClient()
		instance.NamingClient = conf.NamingClient()
	}
	return
}

// nacos客户端配置
func (config *Config) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         config.NameSpace,
		LogDir:              config.LogDir,
		CacheDir:            config.CacheDir,
		Username:            config.Username,
		Password:            config.Password,
	}
}

// nacos服务中间件配置
func (config *Config) ServerConfigs() (serverConfigs []constant.ServerConfig) {
	var adds = strings.Split(config.Address, ",")
	if len(adds) == 0 {
		log.Error("Nacos服务地址不能为空!")
		return
	}
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
		portInt, _ := strconv.ParseInt(port, 10, 64)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: config.Path,
			IpAddr:      host,
			Port:        uint64(portInt),
		})
	}
	return
}

// 初始化Nacos配置中心客户端
func (config *Config) ConfigClient() (client config_client.IConfigClient) {
	var err error
	client, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  config.ClientConfig(),
		ServerConfigs: config.ServerConfigs(),
	})
	if err != nil {
		log.Error("初始化Nacos配置中心客户端-失败! ", config.Format())
		log.Error("error : ", err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (config *Config) NamingClient() (client naming_client.INamingClient) {
	var err error
	client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  config.ClientConfig(),
		ServerConfigs: config.ServerConfigs(),
	})
	if err != nil {
		log.Error("初始化Nacos服务发现客户端-失败! ", config.Format())
		log.Error("error : ", err)
		return
	}
	return
}
