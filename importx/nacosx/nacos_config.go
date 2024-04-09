package nacosx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
)

const (
	OnlyConfig      = iota // 仅用配置中心
	OnlyNaming             // 仅用服务发现
	ConfigAndNaming        // 配置中心和服务发现都使用
)

// nacos访问配置
type Nacos struct {
	Address   string `yaml:"address" json:"address"`                          // nacos服务地址,多个以英文逗号分割
	Username  string `yaml:"username" json:"username" default:"nacos"`        // 用户名
	Password  string `yaml:"password" json:"password" default:"nacos"`        // 密码
	NameSpace string `yaml:"nameSpace" json:"nameSpace" default:"public"`     // 命名空间
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                    // 模式
	LogDir    string `yaml:"logDir" json:"logDir" default:".nacos/log"`       // 日志文件夹
	CacheDir  string `yaml:"cacheDir" json:"cacheDir" default:".nacos/cache"` // 缓存文件夹
}

// 配置信息格式化
func (n *Nacos) ToString(title string) string {
	return fmt.Sprintf("%s => address=%s username=%s password=%s nameSpace=%s mode=%d",
		title, n.AddressUrl(), n.Username, n.Password, n.NameSpace, n.Mode)
}

// 配置器名称
func (n *Nacos) Title() string {
	return "init nacos"
}

// 配置文件读取
func (*Nacos) Reader() *configx.Reader {
	return nil
}

// 配置器运行
func (n *Nacos) Run() (err error) {
	if handler == nil {
		handler = &Handler{Config: n}
		switch n.Mode {
		case OnlyConfig:
			if handler.ConfigClient, err = n.ConfigClient(); err != nil {
				return
			}
		case OnlyNaming:
			if handler.NamingClient, err = n.NamingClient(); err != nil {
				return
			}
		case ConfigAndNaming:
			if handler.ConfigClient, err = n.ConfigClient(); err != nil {
				return
			}
			if handler.NamingClient, err = n.NamingClient(); err != nil {
				return
			}
		}
	}
	log.Info(n.ToString("connect nacos successful!"))
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
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         n.NameSpace,
		LogDir:              n.LogDir,
		CacheDir:            n.CacheDir,
		Username:            n.Username,
		Password:            n.Password,
	}
}

// nacos服务中间件配置
func (n *Nacos) ServerConfigs() (serverConfigs []constant.ServerConfig) {
	var adds = strings.Split(n.Address, ",")
	if len(adds) == 0 {
		log.Error("nacos.address cannot be empty!")
		return
	}
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
		portInt, _ := strconv.ParseInt(port, 10, 64)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: "/nacos",
			IpAddr:      host,
			Port:        uint64(portInt),
		})
	}
	return
}

// 初始化Nacos配置中心客户端
func (n *Nacos) ConfigClient() (client config_client.IConfigClient, err error) {
	if client, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  n.ClientConfig(),
		ServerConfigs: n.ServerConfigs(),
	}); err != nil {
		log.Error(n.ToString("init nacos config client failed !"))
		log.Error("error : ", err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (n *Nacos) NamingClient() (client naming_client.INamingClient, err error) {
	if client, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  n.ClientConfig(),
		ServerConfigs: n.ServerConfigs(),
	}); err != nil {
		log.Error(n.ToString("init nacos naming client failed!"))
		log.Error("error : ", err)
		return
	}
	return
}
