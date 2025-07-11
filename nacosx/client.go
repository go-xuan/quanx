package nacosx

import (
	"reflect"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

var _client *Client

func this() *Client {
	if _client == nil {
		panic("nacos client not initialized, please check the relevant config")
	}
	return _client
}

// Initialized 是否初始化
func Initialized() bool {
	return _client != nil
}

// ReadConfig 读取nacos配置
func ReadConfig(config any, param vo.ConfigParam, listen ...bool) error {
	if _, err := this().ReadConfig(config, param); err != nil {
		return errorx.Wrap(err, "read nacos config error")
	}
	if len(listen) > 0 && listen[0] {
		if err := this().ListenConfig(config, param); err != nil {
			return errorx.Wrap(err, "listen nacos config error")
		}
	}
	return nil
}

// Client nacos客户端
type Client struct {
	config       *Config                     // nacos配置
	configClient config_client.IConfigClient // nacos配置中心客户端
	namingClient naming_client.INamingClient // nacos服务发现客户端
}

// GetConfigClient 获取配置中心客户端
func (c *Client) GetConfigClient() config_client.IConfigClient {
	if client := this().configClient; client == nil {
		panic("the nacos config client has not been initialized")
	} else {
		return client
	}
}

// GetNamingClient 获取服务发现客户端
func (c *Client) GetNamingClient() naming_client.INamingClient {
	if client := this().namingClient; client == nil {
		panic("the nacos naming client has not been initialized")
	} else {
		return client
	}
}

// GetConfig 获取nacos配置
func (c *Client) GetConfig(param vo.ConfigParam) (string, error) {
	if content, err := c.GetConfigClient().GetConfig(param); err != nil {
		return "", errorx.Wrap(err, "get nacos config error")
	} else {
		return content, nil
	}
}

// PublishConfig 发布nacos配置
func (c *Client) PublishConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().PublishConfig(param); err != nil {
		return errorx.Wrap(err, "publish nacos config error")
	}
	return nil
}

// DeleteConfig 删除nacos配置
func (c *Client) DeleteConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().DeleteConfig(param); err != nil {
		return errorx.Wrap(err, "delete nacos config error")
	}
	return nil
}

func (c *Client) ReadConfig(config any, param vo.ConfigParam) ([]byte, error) {
	// 配置值必须是指针类型，否则不允许读取
	if typeOf := reflect.TypeOf(config); typeOf.Kind() != reflect.Pointer {
		return nil, errorx.New("the scanned object must be of pointer type")
	}
	// 读取配置文件内容
	content, err := c.GetConfig(param)
	if err != nil {
		return nil, errorx.Wrap(err, "read nacos config error")
	}
	data := []byte(content)
	if err = marshalx.Apply(param.DataId).Unmarshal(data, config); err != nil {
		return nil, errorx.Wrap(err, "unmarshal nacos config data error")
	}
	return data, nil
}

// ListenConfig 监听nacos配置
func (c *Client) ListenConfig(config any, param vo.ConfigParam) error {
	// 配置监听响应方法
	param.OnChange = func(namespace, group, dataId, data string) {
		log.WithField("dataId", dataId).
			WithField("group", group).
			WithField("namespace", namespace).
			WithField("data", data).
			Info("the nacos config content has changed !!!")
		if err := marshalx.Apply(dataId).Unmarshal([]byte(data), config); err != nil {
			log.Errorf("update config error, group: %s; dataId: %s; data: %s", group, dataId, data)
		}
	}
	if err := c.GetConfigClient().ListenConfig(param); err != nil {
		return errorx.Wrap(err, "listen nacos config failed")
	}
	return nil
}

// RegisterInstance 注册nacos服务实例
func (c *Client) RegisterInstance(instance *ServerInstance) error {
	if ok, err := c.GetNamingClient().RegisterInstance(vo.RegisterInstanceParam{
		Ip:          instance.GetIP(),
		Port:        uint64(instance.GetPort()),
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		ServiceName: instance.GetName(),
		GroupName:   instance.Group,
		Ephemeral:   true,
	}); !ok || err != nil {
		return errorx.Wrap(err, "nacos register instance failed")
	}
	return nil
}

// DeregisterInstance 注销nacos服务实例
func (c *Client) DeregisterInstance(instance *ServerInstance) error {
	if ok, err := c.GetNamingClient().DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          instance.GetIP(),
		Port:        uint64(instance.GetPort()),
		ServiceName: instance.GetName(),
		GroupName:   instance.Group,
		Ephemeral:   true,
	}); !ok || err != nil {
		return errorx.Wrap(err, "nacos deregister instance failed")
	}
	return nil
}

// SelectOneHealthyInstance 选择一个健康实例
func (c *Client) SelectOneHealthyInstance(server, group string) (*model.Instance, error) {
	if instance, err := c.GetNamingClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: server,
		GroupName:   group,
	}); err != nil {
		return nil, errorx.Wrap(err, "select one healthy instance failed")
	} else {
		return instance, nil
	}
}

// SelectInstances 选择实例列表
func (c *Client) SelectInstances(server, group string) ([]model.Instance, error) {
	if instances, err := c.GetNamingClient().SelectInstances(vo.SelectInstancesParam{
		ServiceName: server,
		GroupName:   group,
		HealthyOnly: true,
	}); err != nil {
		return nil, errorx.Wrap(err, "select instances failed")
	} else {
		return instances, nil
	}
}
