package nacosx

import (
	"fmt"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// 服务注册
type ServerConfig struct {
	Ip    string `yaml:"ip"`    // IP
	Port  string `yaml:"port"`  // 端口
	Group string `yaml:"group"` // 分组
	Name  string `yaml:"name"`  // 名称
}

func (s ServerConfig) Format() string {
	return fmt.Sprintf("group=%s name=%s", s.Group, s.Name)
}

// 注册Nacos服务实例
func RegisterInstance(server ServerConfig) {
	port, _ := strconv.Atoi(server.Port)
	_, err := CTL.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.Ip,
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
