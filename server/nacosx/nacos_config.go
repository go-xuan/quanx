package nacosx

import (
	"fmt"
	"github.com/go-xuan/quanx/common/constx"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/server/confx"
)

const (
	OnlyConfig      = iota // 仅用配置中心
	OnlyNaming             // 仅用服务发现
	ConfigAndNaming        // 配置中心和服务发现都使用
)

// nacos访问配置
type NacosConfig struct {
	Address   string `yaml:"address" json:"address" default:"127.0.0.1"`  // nacos服务地址,多个以英文逗号分割
	Username  string `yaml:"username" json:"username" default:"nacos"`    // 用户名
	Password  string `yaml:"password" json:"password" default:"nacos"`    // 密码
	NameSpace string `yaml:"nameSpace" json:"nameSpace" default:"public"` // 命名空间
	Mode      int    `yaml:"mode" json:"mode" default:"2"`                // 模式（0-仅用配置中心；1-仅用服务发现；2-配置中心和服务发现都用）
}

// 配置信息格式化
func (conf *NacosConfig) ToString(title string) string {
	return fmt.Sprintf("%s => address=%s username=%s password=%s nameSpace=%s mode=%d",
		title, conf.AddressUrl(), conf.Username, conf.Password, conf.NameSpace, conf.Mode)
}

// 配置器名称
func (conf *NacosConfig) Theme() string {
	return "Nacos"
}

// 配置文件读取
func (*NacosConfig) Reader() *confx.Reader {
	return nil
}

// 配置器运行
func (conf *NacosConfig) Run() (err error) {
	if handler == nil {
		handler = &Handler{Config: conf}
		var clientParam = vo.NacosClientParam{
			ClientConfig:  conf.ClientConfig(),
			ServerConfigs: conf.ServerConfigs(),
		}
		switch conf.Mode {
		case OnlyConfig:
			if handler.ConfigClient, err = conf.ConfigClient(clientParam); err != nil {
				return
			}
		case OnlyNaming:
			if handler.NamingClient, err = conf.NamingClient(clientParam); err != nil {
				return
			}
		case ConfigAndNaming:
			if handler.ConfigClient, err = conf.ConfigClient(clientParam); err != nil {
				return
			}
			if handler.NamingClient, err = conf.NamingClient(clientParam); err != nil {
				return
			}
		}
	}
	log.Info(conf.ToString("Nacos connect successful!"))
	return
}

// nacos访问地址
func (conf *NacosConfig) AddressUrl() string {
	return conf.Address + "/nacos"
}

// 开启服务注册
func (conf *NacosConfig) EnableNaming() bool {
	return conf.Mode == OnlyNaming || conf.Mode == ConfigAndNaming
}

// nacos客户端配置
func (conf *NacosConfig) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		Username:            conf.Username,
		Password:            conf.Password,
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         conf.NameSpace,
		LogDir:              filepath.Join(constx.ResourceDir, ".nacos/log"),
		CacheDir:            filepath.Join(constx.ResourceDir, ".nacos/cache"),
	}
}

// nacos服务中间件配置
func (conf *NacosConfig) ServerConfigs() (serverConfigs []constant.ServerConfig) {
	var adds = strings.Split(conf.Address, ",")
	if len(adds) == 0 {
		log.Error("the address of nacos cannot be empty!")
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
func (conf *NacosConfig) ConfigClient(param vo.NacosClientParam) (client config_client.IConfigClient, err error) {
	if client, err = clients.NewConfigClient(param); err != nil {
		log.Error(conf.ToString("init nacos config client failed !"))
		log.Error("error : ", err)
		return
	}
	return
}

// 初始化Nacos服务发现客户端
func (conf *NacosConfig) NamingClient(param vo.NacosClientParam) (client naming_client.INamingClient, err error) {
	if client, err = clients.NewNamingClient(param); err != nil {
		log.Error(conf.ToString("init nacos naming client failed!"))
		log.Error("error : ", err)
		return
	}
	return
}
