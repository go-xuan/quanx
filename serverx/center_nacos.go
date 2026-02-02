package serverx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosCenter(group string, client naming_client.INamingClient) *NacosCenter {
	return &NacosCenter{
		group:  group,
		client: client,
	}
}

// NacosCenter 服务中心
type NacosCenter struct {
	group  string                      // 服务分组
	client naming_client.INamingClient // nacos服务发现客户端
}

func (c *NacosCenter) Register(instance Instance) error {
	if _, err := c.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          instance.GetHost(),
		Port:        uint64(instance.GetPort()),
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		ServiceName: instance.GetName(),
		GroupName:   c.group,
		Ephemeral:   true,
	}); err != nil {
		return errorx.Wrap(err, "register nacos server instance failed")
	}
	return nil
}

func (c *NacosCenter) Deregister(instance Instance) error {
	if _, err := c.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          instance.GetHost(),
		Port:        uint64(instance.GetPort()),
		ServiceName: instance.GetName(),
		GroupName:   c.group,
		Ephemeral:   true,
	}); err != nil {
		return errorx.Wrap(err, "deregister nacos server instance failed")
	}
	return nil
}

func (c *NacosCenter) SelectOne(server string) (Instance, error) {
	instance, err := c.client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: server,
		GroupName:   c.group,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "select one healthy nacos server instance failed")
	}
	return &NacosInstance{
		Id:   instance.InstanceId,
		Name: instance.ServiceName,
		Host: instance.Ip,
		Port: int(instance.Port),
	}, nil
}

func (c *NacosCenter) SelectAll(server string) ([]Instance, error) {
	instances, err := c.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: server,
		GroupName:   c.group,
		HealthyOnly: true,
	})
	if err != nil {
		return nil, errorx.Wrap(err, "select nacos server instances failed")
	}
	var result []Instance
	for _, instance := range instances {
		result = append(result, &NacosInstance{
			Id:   instance.InstanceId,
			Name: instance.ServiceName,
			Host: instance.Ip,
			Port: int(instance.Port),
		})
	}
	return result, nil
}

// NacosInstance nacos服务实例
type NacosInstance struct {
	Id   string `json:"id" yaml:"id"`     // 实例id
	Name string `json:"name" yaml:"name"` // 服务名
	Host string `json:"host" yaml:"host"` // 服务host
	Port int    `json:"port" yaml:"port"` // 服务端口
}

func (i *NacosInstance) GetID() string {
	return i.Id
}

func (i *NacosInstance) GetDomain() string {
	return fmt.Sprintf("http://%s:%d", i.Host, i.Port)
}

func (i *NacosInstance) GetName() string {
	return i.Name
}

func (i *NacosInstance) GetHost() string {
	return i.Host
}

func (i *NacosInstance) GetPort() int {
	return i.Port
}
