package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-xuan/quanx/utils/marshalx"
)

var handler *Handler

type Handler struct {
	Config       *Config                     // nacos配置
	ConfigClient config_client.IConfigClient // nacos配置中心客户端
	NamingClient naming_client.INamingClient // nacos服务发现客户端
}

func This() *Handler {
	if handler == nil {
		panic("the nacos handler has not been initialized, please check the relevant config")
	}
	return handler
}

// GetConfigType 获取配置文件类型
func GetConfigType(dataId string) vo.ConfigType {
	for i := len(dataId) - 1; i >= 0; i-- {
		if dataId[i] == '.' {
			return vo.ConfigType(dataId[i+1:])
		}
	}
	return ""
}

// ReadConfigContent 获取nacos配置文件内容
func ReadConfigContent(group, dataId string) (content string, err error) {
	var param = vo.ConfigParam{Group: group, DataId: dataId, Type: GetConfigType(dataId)}
	if content, err = This().ConfigClient.GetConfig(param); err != nil {
		return
	}
	return
}

// BuildConfigFromFile 基于文件构建配置
func (h *Handler) BuildConfigFromFile(path string) (err error) {
	if err = marshalx.UnmarshalFromFile(path, h.Config); err != nil {
		return
	}
	return
}

// BuildConfigFromNacos 基于nacos构建配置
func (h *Handler) BuildConfigFromNacos(group, dataId string) (err error) {
	module := &Scanner{Group: group, DataId: dataId}
	if err = module.Scan(h.Config); err != nil {
		return
	}
	return
}
