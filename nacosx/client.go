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

// PublishConfig 发布nacos配置
func (c *Client) PublishConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().PublishConfig(param); err != nil {
		return errorx.Wrap(err, "publish nacos config error")
	}
	return nil
}

// GetConfig 获取nacos配置
func (c *Client) GetConfig(param vo.ConfigParam, publishIfNotExist bool) (string, error) {
	content, err := c.GetConfigClient().GetConfig(param)
	if err != nil {
		return "", errorx.Wrap(err, "get nacos config error")
	}
	// 如果配置不存在，则发布配置
	if content == "" && publishIfNotExist {
		if err = c.PublishConfig(param); err != nil {
			return "", errorx.Wrap(err, "publish nacos config error")
		}
	}
	return content, nil
}

// DeleteConfig 删除nacos配置
func (c *Client) DeleteConfig(param vo.ConfigParam) error {
	if _, err := c.GetConfigClient().DeleteConfig(param); err != nil {
		return errorx.Wrap(err, "delete nacos config error")
	}
	return nil
}

// ReadConfig 读取nacos配置
func (c *Client) ReadConfig(config any, param vo.ConfigParam) ([]byte, error) {
	// 配置值必须是指针类型，否则不允许读取
	if typeOf := reflect.TypeOf(config); typeOf.Kind() != reflect.Pointer {
		return nil, errorx.New("the scanned object must be of pointer type")
	}
	// 读取配置文件内容
	content, err := c.GetConfig(param, false)
	if err != nil {
		return nil, errorx.Wrap(err, "read nacos config error")
	}
	if content == "" {
		return nil, errorx.New("read nacos config empty")
	}
	// 解析配置文件内容
	data := []byte(content)
	if err = marshalx.Apply(string(param.Type)).Unmarshal(data, config); err != nil {
		return nil, errorx.Wrap(err, "unmarshal nacos config error")
	}
	return data, nil
}

// ListenConfig 监听nacos配置
func (c *Client) ListenConfig(config any, param vo.ConfigParam) error {
	// 配置监听响应方法
	param.OnChange = func(namespace, group, dataId, data string) {
		logger := log.WithField("dataId", dataId).
			WithField("group", group).
			WithField("namespace", namespace).
			WithField("data", data)
		logger.Info("the nacos config data has changed !!!")
		if err := marshalx.Apply(dataId).Unmarshal([]byte(data), config); err != nil {
			logger.WithField("error", err.Error()).Error("update config error")
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
