package dbx

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

func init() {
	RegisterBuilder("gorm", GormClientBuilder)
}

type Client interface {
	GetConfig() *Config                           // 获取配置
	GetInstance() interface{}                     // 获取数据库连接
	Copy(target, database string) (Client, error) // 复制客户端
	Close() error                                 // 关闭客户端
}

// NewClient 创建数据库客户端
func NewClient(config *Config) (Client, error) {
	client := stringx.IfZero(config.Client, "gorm")
	if clientBuilders != nil {
		if builder, ok := clientBuilders[client]; ok {
			return builder(config)
		}
	}
	return nil, errorx.Sprintf("this database client builder is not registered: %s", client)
}

// 客户端构造函数
var clientBuilders = map[string]ClientBuilder{}

// ClientBuilder 客户端构造函数
type ClientBuilder func(*Config) (Client, error)

// RegisterBuilder 注册客户端构造函数
func RegisterBuilder(client string, builder ClientBuilder) {
	if clientBuilders[client] == nil {
		clientBuilders = make(map[string]ClientBuilder)
	}
	clientBuilders[client] = builder
}
