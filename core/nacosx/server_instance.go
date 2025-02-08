package nacosx

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/go-xuan/quanx/os/errorx"
)

// ServerInstance 服务实例
type ServerInstance struct {
	Group string `yaml:"group"` // 实例分组
	Name  string `yaml:"name"`  // 实例名称
	Host  string `yaml:"host"`  // 服务实例host
	Port  int    `yaml:"port"`  // 实例端口
}

func (s *ServerInstance) Info() string {
	return fmt.Sprintf("group=%s name=%s", s.Group, s.Name)
}

// SelectOneHealthyInstance 随机获取一个健康的服务实例
func SelectOneHealthyInstance(name, group string) (*ServerInstance, error) {
	if instance, err := GetNacosNamingClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: name,
		GroupName:   group,
	}); err != nil {
		return nil, errorx.Wrap(err, "select one healthy instance failed")
	} else {
		return &ServerInstance{
			Name: instance.ServiceName,
			Host: instance.Ip,
			Port: int(instance.Port),
		}, nil
	}
}

// SelectInstances 获取所有健康服务实例
func SelectInstances(name, group string) ([]*ServerInstance, error) {
	if instances, err := GetNacosNamingClient().SelectInstances(vo.SelectInstancesParam{
		ServiceName: name,
		GroupName:   group,
		HealthyOnly: true,
	}); err != nil {
		return nil, errorx.Wrap(err, "select instances failed")
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
