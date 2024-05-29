package serverx

import "fmt"

type Registry interface {
	Register(Instance) error                       // 服务注册
	Deregister(Instance) error                     // 服务注销
	AllInstances(name string) ([]*Instance, error) // 获取服务的所有实例
	GetInstance(name string) (*Instance, error)    // 随机获取一个服务实例
}

// 服务实例
type Instance struct {
	Name string `yaml:"name"` // 实例名称,唯一标识
	Host string `yaml:"host"` // 服务实例host
	Port int    `yaml:"port"` // 实例端口
}

func (s *Instance) Info() string {
	return fmt.Sprintf("name=%s host=%s port=%d", s.Name, s.Host, s.Port)
}
