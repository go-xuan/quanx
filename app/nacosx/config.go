package nacosx

import (
	"fmt"
	"github.com/go-xuan/quanx/app/constx"
	"path/filepath"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/confx"
	"github.com/go-xuan/quanx/app/serverx"
	"github.com/go-xuan/quanx/types/stringx"
)

const (
	OnlyConfig      = iota // 仅用配置中心
	OnlyNaming             // 仅用服务发现
	ConfigAndNaming        // 配置中心和服务发现都使用
)

// nacos访问配置
type Nacos struct {
	Address   string `yaml:"address" json:"address" default:"127.0.0.1"`  // nacos服务地址,多个以英文逗号分割
	Username  string `yaml:"username" json:"username" default:"nacos"`    // 用户名
	Password  string `yaml:"password" json:"password" default:"nacos"`    // 密码
	NameSpace string `yaml:"nameSpace" json:"nameSpace" default:"public"` // 命名空间
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                // 模式（0-仅配置中心；1-仅服务发现；2-配置中心和服务发现）
}

// 配置信息格式化
func (n *Nacos) Info() string {
	return fmt.Sprintf("address=%s username=%s password=%s nameSpace=%s mode=%d",
		n.AddressUrl(), n.Username, n.Password, n.NameSpace, n.Mode)
}

// 配置器标题
func (n *Nacos) Title() string {
	return "Nacos"
}

// 配置文件读取
func (*Nacos) Reader() *confx.Reader {
	return nil
}

// 配置器运行
func (n *Nacos) Run() (err error) {
	if handler == nil {
		handler = &Handler{Config: n}
		switch n.Mode {
		case OnlyConfig:
			if handler.ConfigClient, err = n.ConfigClient(n.ClientParam()); err != nil {
				return
			}
		case OnlyNaming:
			if handler.NamingClient, err = n.NamingClient(n.ClientParam()); err != nil {
				return
			}
		case ConfigAndNaming:
			var param = n.ClientParam()
			if handler.ConfigClient, err = n.ConfigClient(param); err != nil {
				return
			}
			if handler.NamingClient, err = n.NamingClient(param); err != nil {
				return
			}
		}
	}
	log.Info("Nacos Connect Successful: ", n.Info())
	return
}

// nacos访问地址
func (n *Nacos) AddressUrl() string {
	return n.Address + "/nacos"
}

// 开启服务注册
func (n *Nacos) EnableNaming() bool {
	return n.Mode == OnlyNaming || n.Mode == ConfigAndNaming
}

// nacos客户端配置
func (n *Nacos) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		Username:            n.Username,
		Password:            n.Password,
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         n.NameSpace,
		LogDir:              filepath.Join(constx.DefaultResourceDir, ".nacos/log"),
		CacheDir:            filepath.Join(constx.DefaultResourceDir, ".nacos/cache"),
	}
}

// nacos服务中间件配置
func (n *Nacos) ServerConfigs() []constant.ServerConfig {
	var adds = strings.Split(n.Address, ",")
	if len(adds) == 0 {
		log.Error("the address of nacos cannot be empty")
		return nil
	}
	var configs []constant.ServerConfig
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
		configs = append(configs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      host,
			Port:        uint64(stringx.ToInt64(port)),
		})
	}
	return configs
}

func (n *Nacos) ClientParam() vo.NacosClientParam {
	return vo.NacosClientParam{
		ClientConfig:  n.ClientConfig(),
		ServerConfigs: n.ServerConfigs(),
	}
}

// 初始化Nacos配置中心客户端
func (n *Nacos) ConfigClient(param vo.NacosClientParam) (client config_client.IConfigClient, err error) {
	if client, err = clients.NewConfigClient(param); err != nil {
		log.Error("Nacos Config Client Init Failed: ", n.Info())
		log.Error(err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (n *Nacos) NamingClient(param vo.NacosClientParam) (client naming_client.INamingClient, err error) {
	if client, err = clients.NewNamingClient(param); err != nil {
		log.Error("Nacos Naming Client Init Failed: ", n.Info(), err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (n *Nacos) Register(server serverx.Instance) (err error) {
	var client naming_client.INamingClient
	if client = This().NamingClient; client == nil {
		if client, err = n.NamingClient(n.ClientParam()); err != nil {
			return
		}
	}
	if _, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(server.Port),
		GroupName:   n.NameSpace,
		ServiceName: server.Name,
		Weight:      1,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    nil,
	}); err != nil {
		log.Error("Nacos Server Register Failed: ", server.Info(), err)
	} else {
		log.Info("Nacos Server Register Successful: ", server.Info())
	}
	return
}

func (n *Nacos) Deregister(server serverx.Instance) (err error) {
	var client naming_client.INamingClient
	if client = This().NamingClient; client == nil {
		if client, err = n.NamingClient(n.ClientParam()); err != nil {
			return
		}
	}
	if _, err = client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          server.Host,
		Port:        uint64(server.Port),
		GroupName:   n.NameSpace,
		ServiceName: server.Name,
		Ephemeral:   true,
	}); err != nil {
		log.Error("Nacos Server Deregister Failed: ", server.Info(), err)
	} else {
		log.Info("Nacos Server Deregister Successful: ", server.Info())
	}
	return
}

func (n *Nacos) AllInstances(name string) (instances []*serverx.Instance, err error) {
	var servers []model.Instance
	if servers, err = This().NamingClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: name,
		GroupName:   n.NameSpace,
		HealthyOnly: true,
	}); err != nil {
		return
	}
	for _, server := range servers {
		instances = append(instances, &serverx.Instance{
			Name: server.ServiceName,
			Host: server.Ip,
			Port: int(server.Port),
		})
	}
	return
}

func (n *Nacos) GetInstance(name string) (instance *serverx.Instance, err error) {
	var servers *model.Instance
	if servers, err = This().NamingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: name,
		GroupName:   n.NameSpace,
	}); err != nil {
		return
	}
	instance = &serverx.Instance{
		Name: servers.ServiceName,
		Host: servers.Ip,
		Port: int(servers.Port),
	}
	return
}
