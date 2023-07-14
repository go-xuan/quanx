package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

var CTL *Control

type Control struct {
	Config       *Config                     // nacos配置
	ConfigClient config_client.IConfigClient // nacos配置中心客户端
	NamingClient naming_client.INamingClient // nacos服务发现客户端
}

func InitNacosCTL(conf *Config) {
	if conf != nil && conf.Address != "" {
		InitNacos(conf)
	}
}

const (
	ConfigAndNaming = iota // 配置中心和服务发现都使用
	OnlyConfig             // 仅用配置中心
	OnlyNaming             // 仅用服务发现
)

func InitNacos(conf *Config) {
	if CTL == nil {
		CTL = &Control{Config: conf}
		switch conf.Mode {
		case OnlyConfig:
			InitConfigClient()
		case OnlyNaming:
			InitNamingClient()
		default:
			InitConfigClient()
			InitNamingClient()
		}
	}
}

// 初始化Nacos配置客户端
func InitConfigClient() {
	if CTL.ConfigClient == nil {
		serverConfigs := CTL.Config.ServerConfigList()
		if len(serverConfigs) == 0 {
			log.Error("nacos服务地址不能为空！")
			return
		}
		client, err := clients.NewConfigClient(vo.NacosClientParam{
			ClientConfig:  CTL.Config.ClientConfig(),
			ServerConfigs: serverConfigs,
		})
		if err != nil {
			log.Error("初始化nacos配置中心客户端-失败 : ", err)
			return
		}
		CTL.ConfigClient = client
	}
	return
}

// 初始化nacos服务发现客户端
func InitNamingClient() {
	if CTL.NamingClient == nil {
		serverConfigs := CTL.Config.ServerConfigList()
		if len(serverConfigs) == 0 {
			log.Error("nacos服务地址不能为空！")
			return
		}
		client, err := clients.NewNamingClient(vo.NacosClientParam{
			ClientConfig:  CTL.Config.ClientConfig(),
			ServerConfigs: serverConfigs,
		})
		if err != nil {
			log.Error("初始化nacos服务发现客户端-失败 : ", err)
			return
		}
		CTL.NamingClient = client
	}
	return
}
