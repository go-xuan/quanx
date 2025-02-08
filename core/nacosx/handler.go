package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/anyx"
)

var handler *Handler

func this() *Handler {
	if handler == nil {
		panic("the nacos handler has not been initialized, please check the relevant config")
	}
	return handler
}

type Handler struct {
	config       *Config                     // nacos配置
	configClient config_client.IConfigClient // nacos配置中心客户端
	namingClient naming_client.INamingClient // nacos服务发现客户端
}

// ScanConfig 从nacos获取配置并扫描
func ScanConfig(v any, group, dataId string, listen ...bool) error {
	var scanner = &Scanner{
		Group:  group,
		DataId: dataId,
		Type:   vo.ConfigType(filex.GetSuffix(dataId)),
		Listen: anyx.Default(false, listen...),
	}
	if err := scanner.Scan(v); err != nil {
		return errorx.Wrap(err, "nacos config scan failed")
	}
	return nil
}

// GetNacosConfigClient 获取配置中心客户端
func GetNacosConfigClient() config_client.IConfigClient {
	if client := this().configClient; client == nil {
		panic("the nacos config client has not been initialized")
	} else {
		return client
	}
}

// GetNacosNamingClient 获取服务中心客户端
func GetNacosNamingClient() naming_client.INamingClient {
	if client := this().namingClient; client == nil {
		panic("the nacos naming client has not been initialized")
	} else {
		return client
	}
}

// Register 注册Nacos服务
func Register(server ServerInstance) error {
	if ok, err := GetNacosNamingClient().RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(server.Port),
		GroupName:   server.Group,
		ServiceName: server.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}); !ok || err != nil {
		log.Error("nacos server register failed: ", server.Info())
		return errorx.Wrap(err, "nacos server register failed")
	} else {
		log.Info("nacos server register successfully: ", server.Info())
		return nil
	}
}

// Deregister 注销nacos服务
func Deregister(server ServerInstance) error {
	if ok, err := GetNacosNamingClient().DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(server.Port),
		GroupName:   server.Group,
		ServiceName: server.Name,
		Ephemeral:   true,
	}); !ok || err != nil {
		log.Error("nacos server deregister failed: ", server.Info())
		return errorx.Wrap(err, "nacos server deregister failed")
	} else {
		log.Info("nacos server deregister successfully: ", server.Info())
		return nil
	}
}
