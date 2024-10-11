package nacosx

import (
	"fmt"

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
func SelectOneHealthyInstance(name, group string) (string, error) {
	if instance, err := This().SelectOneHealthyInstance(name, group); err != nil {
		return "", errorx.Wrap(err, "select one healthy instances error")
	} else {
		return fmt.Sprintf("%s:%d", instance.Host, instance.Port), nil
	}
}

// SelectInstances 获取所有健康服务实例
func SelectInstances(name, group string) ([]string, error) {
	if instances, err := This().SelectInstances(name, group); err != nil {
		return nil, errorx.Wrap(err, "select instances error")
	} else {
		var addrs []string
		for _, instance := range instances {
			addrs = append(addrs, fmt.Sprintf("%s:%d", instance.Host, instance.Port))
		}
		return addrs, nil
	}
}
