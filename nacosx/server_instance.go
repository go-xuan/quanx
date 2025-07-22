package nacosx

import (
	"fmt"
)

// ServerInstance 服务实例
type ServerInstance struct {
	Id    string `json:"id" yaml:"id"`       // 实例ID
	Name  string `json:"name" yaml:"name"`   // 服务名
	Group string `json:"group" yaml:"group"` // 服务分组
	IP    string `json:"ip" yaml:"ip"`       // 服务IP
	Port  int    `json:"port" yaml:"port"`   // 服务端口
}

func (s *ServerInstance) GetID() string {
	return s.Id
}

func (s *ServerInstance) GetDomain() string {
	return fmt.Sprintf("http://%s:%d", s.IP, s.Port)
}

func (s *ServerInstance) GetName() string {
	return s.Name
}

func (s *ServerInstance) GetIP() string {
	return s.IP
}

func (s *ServerInstance) GetPort() int {
	return s.Port
}

func (s *ServerInstance) GetStatus() string {
	return "UP"
}
