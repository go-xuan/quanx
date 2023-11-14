package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"

	"github.com/go-xuan/quanx/utils/structx"
)

var CTL *Controller

type Controller struct {
	Config       *Config                     // nacos配置
	ConfigClient config_client.IConfigClient // nacos配置中心客户端
	NamingClient naming_client.INamingClient // nacos服务发现客户端
}

func Init(conf *Config) {
	CTL = &Controller{Config: conf}
	switch conf.Mode {
	case OnlyConfig:
		CTL.ConfigClient = conf.ConfigClient()
	case OnlyNaming:
		CTL.NamingClient = conf.NamingClient()
	default:
		CTL.ConfigClient = conf.ConfigClient()
		CTL.NamingClient = conf.NamingClient()
	}
	return
}

const (
	ConfigAndNaming = iota // 配置中心和服务发现都使用，默认项
	OnlyConfig             // 仅用配置中心
	OnlyNaming             // 仅用服务发现
)

func (ctl *Controller) BuildConfigFromFile(filePath string) (err error) {
	err = structx.ReadFileToPointer(filePath, ctl.Config)
	if err != nil {
		return
	}
	return
}

func (ctl *Controller) BuildConfigFromNacos(group, dataId string) (err error) {
	module := &Item{Group: group, DataId: dataId}
	err = module.LoadConfig(ctl.Config)
	if err != nil {
		return
	}
	return
}
