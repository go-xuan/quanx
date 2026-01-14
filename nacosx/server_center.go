package nacosx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"

	"github.com/go-xuan/quanx/serverx"
)

// ServerInstance 服务实例
type ServerInstance struct {
	Id   string `json:"id" yaml:"id"`     // 实例id
	Name string `json:"name" yaml:"name"` // 服务名
	Host string `json:"host" yaml:"host"` // 服务host
	Port int    `json:"port" yaml:"port"` // 服务端口
}

func (s *ServerInstance) GetID() string {
	return s.Id
}

func (s *ServerInstance) GetDomain() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}

func (s *ServerInstance) GetName() string {
	return s.Name
}

func (s *ServerInstance) GetHost() string {
	return s.Host
}

func (s *ServerInstance) GetPort() int {
	return s.Port
}

// ServerCenter 服务中心
type ServerCenter struct{}

func (c *ServerCenter) Register(instance serverx.Instance) error {
	if err := GetClient().RegisterInstance(&ServerInstance{
		Name: instance.GetName(),
		Host: instance.GetHost(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "register nacos server instance failed")
	}
	return nil
}

func (c *ServerCenter) Deregister(instance serverx.Instance) error {
	if err := GetClient().DeregisterInstance(&ServerInstance{
		Name: instance.GetName(),
		Host: instance.GetHost(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "deregister nacos server instance failed")
	}
	return nil
}

func (c *ServerCenter) SelectOne(server string) (serverx.Instance, error) {
	instance, err := GetClient().SelectOneHealthyInstance(server)
	if err != nil {
		return nil, errorx.Wrap(err, "select one healthy nacos server instance failed")
	}
	return &ServerInstance{
		Id:   instance.InstanceId,
		Name: instance.ServiceName,
		Host: instance.Ip,
		Port: int(instance.Port),
	}, nil
}

func (c *ServerCenter) SelectAll(server string) ([]serverx.Instance, error) {
	instances, err := GetClient().SelectInstances(server)
	if err != nil {
		return nil, errorx.Wrap(err, "select nacos server instances failed")
	}
	var result []serverx.Instance
	for _, instance := range instances {
		result = append(result, &ServerInstance{
			Id:   instance.InstanceId,
			Name: instance.ServiceName,
			Host: instance.Ip,
			Port: int(instance.Port),
		})
	}
	return result, nil
}
