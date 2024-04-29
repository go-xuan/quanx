package nacosx

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/server/confx"
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
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                // 模式（0-仅用配置中心；1-仅用服务发现；2-配置中心和服务发现都用）
}

// 配置信息格式化
func (n *Nacos) ToString() string {
	return fmt.Sprintf("address=%s username=%s password=%s nameSpace=%s mode=%d",
		n.AddressUrl(), n.Username, n.Password, n.NameSpace, n.Mode)
}

// 配置器名称
func (n *Nacos) Theme() string {
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
		var clientParam = vo.NacosClientParam{
			ClientConfig:  n.ClientConfig(),
			ServerConfigs: n.ServerConfigs(),
		}
		switch n.Mode {
		case OnlyConfig:
			if handler.ConfigClient, err = n.ConfigClient(clientParam); err != nil {
				return
			}
		case OnlyNaming:
			if handler.NamingClient, err = n.NamingClient(clientParam); err != nil {
				return
			}
		case ConfigAndNaming:
			if handler.ConfigClient, err = n.ConfigClient(clientParam); err != nil {
				return
			}
			if handler.NamingClient, err = n.NamingClient(clientParam); err != nil {
				return
			}
		}
	}
	log.Info("Nacos Connect Successful", n.ToString())
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
		LogDir:              filepath.Join(constx.ResourceDir, ".nacos/log"),
		CacheDir:            filepath.Join(constx.ResourceDir, ".nacos/cache"),
	}
}

// nacos服务中间件配置
func (n *Nacos) ServerConfigs() (serverConfigs []constant.ServerConfig) {
	var adds = strings.Split(n.Address, ",")
	if len(adds) == 0 {
		log.Error("the address of nacos cannot be empty")
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
func (n *Nacos) ConfigClient(param vo.NacosClientParam) (client config_client.IConfigClient, err error) {
	if client, err = clients.NewConfigClient(param); err != nil {
		log.Error("Nacos Config Client Init Failed : ", n.ToString(), err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (n *Nacos) NamingClient(param vo.NacosClientParam) (client naming_client.INamingClient, err error) {
	if client, err = clients.NewNamingClient(param); err != nil {
		log.Error("Nacos Naming Client Init Failed : ", n.ToString(), err)
		return
	}
	return
}
