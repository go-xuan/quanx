package nacosx

import (
	"fmt"

	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/errorx"
)

// ServerInstance 服务实例
type ServerInstance struct {
	Id   string `json:"id" yaml:"id"`     // 实例ID
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
	if err := this().RegisterInstance(&ServerInstance{
		Name: instance.GetName(),
		Host: instance.GetHost(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos server center register instance failed")
	}
	return nil
}

func (c *ServerCenter) Deregister(instance serverx.Instance) error {
	if err := this().DeregisterInstance(&ServerInstance{
		Name: instance.GetName(),
		Host: instance.GetHost(),
		Port: instance.GetPort(),
	}); err != nil {
		return errorx.Wrap(err, "nacos server center deregister instance failed")
	}
	return nil
}

func (c *ServerCenter) SelectOne(server string) (serverx.Instance, error) {
	if instance, err := this().SelectOneHealthyInstance(server); err != nil {
		return nil, errorx.Wrap(err, "nacos server center select one healthy instance failed")
	} else {
		return &ServerInstance{
			Id:   instance.InstanceId,
			Name: instance.ServiceName,
			Host: instance.Ip,
			Port: int(instance.Port),
		}, nil
	}
}

func (c *ServerCenter) SelectAll(server string) ([]serverx.Instance, error) {
	if instances, err := this().SelectInstances(server); err != nil {
		return nil, errorx.Wrap(err, "nacos server center select server instances failed")
	} else {
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
}
