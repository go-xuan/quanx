package nacosx

import (
	"github.com/go-xuan/quanx/importx/marshalx"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var handler *Handler

type Handler struct {
	Config       *Nacos                      // nacos配置
	ConfigClient config_client.IConfigClient // nacos配置中心客户端
	NamingClient naming_client.INamingClient // nacos服务发现客户端
}

func This() *Handler {
	if handler == nil {
		panic("the nacos handler has not been initialized, please check the relevant config")
	}
	return handler
}

// 获取配置文件类型
func GetConfigType(dataId string) vo.ConfigType {
	for i := len(dataId) - 1; i >= 0; i-- {
		if dataId[i] == '.' {
			return vo.ConfigType(dataId[i+1:])
		}
	}
	return ""
}

func GetConfigContent(group, dataId string) (content string, err error) {
	var param = vo.ConfigParam{Group: group, DataId: dataId, Type: GetConfigType(dataId)}
	content, err = This().ConfigClient.GetConfig(param)
	if err != nil {
		return
	}
	return
}

func (h *Handler) BuildConfigFromFile(filePath string) (err error) {
	err = marshalx.LoadFromFile(h.Config, filePath)
	if err != nil {
		return
	}
	return
}

func (h *Handler) BuildConfigFromNacos(group, dataId string) (err error) {
	module := &Config{Group: group, DataId: dataId}
	err = module.LoadConfig(h.Config)
	if err != nil {
		return
	}
	return
}
