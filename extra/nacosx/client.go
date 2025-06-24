package nacosx

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/utils/anyx"
)

var _client *Client

// Initialized 是否初始化
func Initialized() bool {
	return _client != nil
}

func this() *Client {
	if _client == nil {
		panic("nacos client not initialized, please check the relevant config")
	}
	return _client
}

type Client struct {
	config       *Config                     // nacos配置
	configClient config_client.IConfigClient // nacos配置中心客户端
	namingClient naming_client.INamingClient // nacos服务发现客户端
}

func (c *Client) Config() *Config {
	return c.config
}

func (c *Client) ConfigClient() config_client.IConfigClient {
	return c.configClient
}

func (c *Client) NamingClient() naming_client.INamingClient {
	return c.namingClient
}

// ReadConfig 从nacos读取配置
func ReadConfig(config any, group, dataId string, listen ...bool) error {
	var reader = &Reader{
		Group:  group,
		DataId: dataId,
		Type:   filex.GetSuffix(dataId),
		Listen: anyx.Default(false, listen...),
	}
	if err := reader.Read(config); err != nil {
		return errorx.Wrap(err, "read nacos config error")
	}
	return nil
}

// GetNacosConfigClient 获取配置中心客户端
func GetNacosConfigClient() config_client.IConfigClient {
	if client := this().ConfigClient(); client == nil {
		panic("the nacos config client has not been initialized")
	} else {
		return client
	}
}

// GetNacosNamingClient 获取服务中心客户端
func GetNacosNamingClient() naming_client.INamingClient {
	if client := this().NamingClient(); client == nil {
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
		log.Info("nacos server register success: ", server.Info())
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
		log.Info("nacos server deregister success: ", server.Info())
		return nil
	}
}
