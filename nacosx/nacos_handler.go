package nacosx

import (
	"github.com/go-xuan/quanx/utilx/structx"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
)

var instance *Handler

type Handler struct {
	Config       *Config                     // nacos配置
	ConfigClient config_client.IConfigClient // nacos配置中心客户端
	NamingClient naming_client.INamingClient // nacos服务发现客户端
}

func This() *Handler {
	if instance == nil {
		panic("The nacos instance has not been initialized, please check the relevant config")
	}
	return instance
}

func (h *Handler) BuildConfigFromFile(filePath string) (err error) {
	err = structx.ReadFileToPointer(h.Config, filePath)
	if err != nil {
		return
	}
	return
}

func (h *Handler) BuildConfigFromNacos(group, dataId string) (err error) {
	module := &Item{Group: group, DataId: dataId}
	err = module.LoadConfig(h.Config)
	if err != nil {
		return
	}
	return
}
