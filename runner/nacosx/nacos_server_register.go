package nacosx

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// 服务注册
type ServerInstance struct {
	Name  string `yaml:"name"`  // 实例名称
	Host  string `yaml:"host"`  // 服务实例host
	Port  int    `yaml:"port"`  // 实例端口
	Group string `yaml:"group"` // 实例分组
}

func (s ServerInstance) Format() string {
	return fmt.Sprintf("group=%s name=%s", s.Group, s.Name)
}

// 注册Nacos服务实例
func RegisterInstance(server ServerInstance) {
	if handler.NamingClient == nil {
		log.Error("未初始化nacos服务发现客户端!")
		return
	}
	if _, err := handler.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(server.Port),
		GroupName:   server.Group,
		ServiceName: server.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    nil,
	}); err != nil {
		log.Error("nacos服务注册失败！", server.Format(), " error : ", err)
	} else {
		log.Info("nacos服务注册成功！", server.Format())
	}
}

// 随机获取一个健康的服务实例
func SelectOneHealthyInstance(serviceName, groupName string) (addr string, err error) {
	var server *model.Instance
	if server, err = This().NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	}); err == nil {
		addr = fmt.Sprintf("%s:%d", server.Ip, server.Port)
	}
	return
}

// 获取所有健康服务实例
func SelectInstances(serviceName, groupName string) (addrs []string, err error) {
	var servers []model.Instance
	if servers, err = handler.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		HealthyOnly: true,
	}); err == nil {
		for _, server := range servers {
			addrs = append(addrs, fmt.Sprintf("%s:%d", server.Ip, server.Port))
		}
	}
	return
}
