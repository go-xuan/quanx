package nacosx

import (
	"fmt"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// 服务注册
type Server struct {
	Name  string `yaml:"name"`  // 实例名称
	Host  string `yaml:"host"`  // 服务实例host
	Port  string `yaml:"port"`  // 实例端口
	Group string `yaml:"group"` // 实例分组
}

func (s Server) Format() string {
	return fmt.Sprintf("group=%s name=%s", s.Group, s.Name)
}

// 注册Nacos服务实例
func RegisterInstance(server Server) {
	port, _ := strconv.Atoi(server.Port)
	_, err := CTL.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(port),
		GroupName:   server.Group,
		ServiceName: server.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    nil,
	})
	msg := server.Format()
	if err != nil {
		log.Error("注册nacos服务-失败! ", msg)
		log.Error("error : ", err)
	} else {
		log.Info("注册nacos服务-成功! ", msg)
	}
}

// 随机获取一个健康的服务实例
func SelectOneHealthyInstance(serviceName, groupName string) (addr string, err error) {
	var instance *model.Instance
	if instance, err = CTL.NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName,
	}); err == nil {
		addr = fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	}
	return
}

// 获取所有健康服务实例
func SelectInstances(serviceName, groupName string) (addrs []string, err error) {
	var instances []model.Instance
	if instances, err = CTL.NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		HealthyOnly: true,
	}); err == nil {
		for _, instance := range instances {
			addrs = append(addrs, fmt.Sprintf("%s:%d", instance.Ip, instance.Port))
		}
	}
	return
}
