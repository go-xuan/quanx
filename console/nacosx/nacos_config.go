package nacosx

import (
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// nacos访问配置
type Config struct {
	Address    string      `yaml:"address"`                         // NACOS服务地址,多个以英文逗号分割
	Path       string      `yaml:"path" default:"/nacos"`           // web访问路径
	Username   string      `yaml:"username" default:"nacos"`        // 用户名
	Password   string      `yaml:"password" default:"nacos"`        // 密码
	Mode       int         `yaml:"mode" default:"0"`                // 模式（0-配置中心和服务发现都用；1-仅用配置中心；2-仅用服务发现）
	NameSpace  string      `yaml:"nameSpace" default:"public"`      // 命名空间
	LogDir     string      `yaml:"logDir" default:".nacos/log"`     // 日志文件夹
	CacheDir   string      `yaml:"cacheDir" default:".nacos/cache"` // 缓存文件夹
	LoadConfig *LoadConfig `yaml:"loadConfig"`                      // 加载配置项
}

// nacos加载配置项
type LoadConfig struct {
	Init   string `yaml:"init"`   // 初始化配置项
	Custom string `yaml:"custom"` // 自定义配置项
	Listen string `yaml:"listen"` // 监听配置项
}

// nacos客户端配置
func (conf *Config) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         conf.NameSpace,
		LogDir:              conf.LogDir,
		CacheDir:            conf.CacheDir,
		Username:            conf.Username,
		Password:            conf.Password,
	}
}

// nacos服务中间件配置
func (conf *Config) ServerConfigs() (serverConfigs []constant.ServerConfig) {
	var adds = strings.Split(conf.Address, ",")
	if len(adds) == 0 {
		log.Error("Nacos服务地址不能为空!")
		return
	}
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
		portInt, _ := strconv.ParseInt(port, 10, 64)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: conf.Path,
			IpAddr:      host,
			Port:        uint64(portInt),
		})
	}
	return
}

// 初始化Nacos配置中心客户端
func (conf *Config) ConfigClient() (client config_client.IConfigClient) {
	var err error
	client, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  conf.ClientConfig(),
		ServerConfigs: conf.ServerConfigs(),
	})
	if err != nil {
		log.Error("初始化Nacos配置中心客户端-失败 : ", err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (conf *Config) NamingClient() (client naming_client.INamingClient) {
	var err error
	client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  conf.ClientConfig(),
		ServerConfigs: conf.ServerConfigs(),
	})
	if err != nil {
		log.Error("初始化Nacos服务发现客户端-失败 : ", err)
		return
	}
	return
}
