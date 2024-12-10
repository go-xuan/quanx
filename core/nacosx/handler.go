package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
)

var handler *Handler

type Handler struct {
	config       *Config                     // nacos配置
	configClient config_client.IConfigClient // nacos配置中心客户端
	namingClient naming_client.INamingClient // nacos服务发现客户端
}

func This() *Handler {
	if handler == nil {
		panic("the nacos handler has not been initialized, please check the relevant config")
	}
	return handler
}

func NacosConfigClient() config_client.IConfigClient {
	if This().configClient == nil {
		panic("the nacos config client has not been initialized")
	}
	return This().configClient
}

func NacosNamingClient() naming_client.INamingClient {
	if This().namingClient == nil {
		panic("the nacos naming client has not been initialized")
	}
	return This().namingClient
}

// Register 注册Nacos服务
func Register(server ServerInstance) error {
	if ok, err := NacosNamingClient().RegisterInstance(vo.RegisterInstanceParam{
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
	if ok, err := NacosNamingClient().DeregisterInstance(vo.DeregisterInstanceParam{
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

func (h *Handler) SelectInstances(name string, group ...string) ([]*ServerInstance, error) {
	var groupName = stringx.Default(h.config.NameSpace, group...)
	if instances, err := NacosNamingClient().SelectInstances(vo.SelectInstancesParam{
		ServiceName: name,
		GroupName:   groupName,
		HealthyOnly: true,
	}); err != nil {
		return nil, err
	} else {
		var servers []*ServerInstance
		for _, instance := range instances {
			servers = append(servers, &ServerInstance{
				Name: instance.ServiceName,
				Host: instance.Ip,
				Port: int(instance.Port),
			})
		}
		return servers, nil
	}
}

func (h *Handler) SelectOneHealthyInstance(name string, group ...string) (*ServerInstance, error) {
	var groupName = stringx.Default(h.config.NameSpace, group...)
	if instance, err := NacosNamingClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: name,
		GroupName:   groupName,
	}); err != nil {
		return nil, err
	} else {
		return &ServerInstance{
			Name: instance.ServiceName,
			Host: instance.Ip,
			Port: int(instance.Port),
		}, nil
	}
}

// ScanConfig 从nacos获取配置并扫描
func (h *Handler) ScanConfig(v any, group, dataId string, listen ...bool) error {
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
