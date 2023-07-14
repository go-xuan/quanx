package nacosx

import (
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

// nacos访问配置
type Config struct {
	Address   string `yaml:"address"`                         // NACOS服务地址,多个以英文逗号分割
	Path      string `yaml:"path" default:"/nacos"`           // web访问路径
	Username  string `yaml:"username" default:"nacos"`        // 用户名
	Password  string `yaml:"password" default:"nacos"`        // 密码
	Mode      int    `yaml:"mode" default:"0"`                // 模式（0-配置中心和服务发现都用；1-仅用配置中心；2-仅用服务发现）
	NameSpace string `yaml:"nameSpace" default:"public"`      // 命名空间
	LogDir    string `yaml:"logDir" default:".nacos/log"`     // 日志文件夹
	CacheDir  string `yaml:"cacheDir" default:".nacos/cache"` // 缓存文件夹
}

// nacos客户端配置
func (conf *Config) ClientConfig() *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:           10 * 1000,
		BeatInterval:        3 * 1000,
		NotLoadCacheAtStart: true,
		NamespaceId:         conf.NameSpace,
		LogDir:              conf.LogDir,
		CacheDir:            conf.CacheDir,
		Username:            conf.Username,
		Password:            conf.Password,
	}
}

// nacos服务中间件配置
func (conf *Config) ServerConfigList() []constant.ServerConfig {
	// 至少一个ServerConfig
	serverConfigs := make([]constant.ServerConfig, 0)
	adds := strings.Split(conf.Address, ",")
	if len(adds) == 0 {
		return nil
	}
	for _, addStr := range adds {
		host, port, _ := strings.Cut(addStr, ":")
		portInt, _ := strconv.ParseInt(port, 10, 64)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			ContextPath: conf.Path,
			IpAddr:      host,
			Port:        uint64(portInt),
		})
	}
	return serverConfigs
}
