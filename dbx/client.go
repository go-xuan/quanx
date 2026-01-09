package dbx

import (
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/constx"
)

var (
	pool     *typex.Enum[string, Client]        // 客户端池
	builders *typex.Enum[string, ClientBuilder] // 客户端构造函数池
)

func init() {
	RegisterClientBuilder("gorm", GormClientBuilder) // 注册gorm客户端构建器
}

// RegisterClientBuilder 注册客户端构造函数
func RegisterClientBuilder(name string, builder ClientBuilder) {
	if builders == nil {
		builders = typex.NewStringEnum[ClientBuilder]()
	}
	builders.Add(name, builder)
}

// NewClient 创建客户端
func NewClient(config *Config) (Client, error) {
	if config.Builder == "" {
		return GormClientBuilder(config)
	}
	if builders != nil {
		if builder, ok := builders.Exist(config.Builder); ok && builder != nil {
			return builder(config)
		}
	}
	return nil, errorx.Sprintf("client builder is not registered: %s", config.Builder)
}

// ClientBuilder 客户端构造函数
type ClientBuilder func(*Config) (Client, error)

// Client 数据库客户端接口
type Client interface {
	GetInstance() interface{} // 获取实例
	GetConfig() *Config       // 获取配置
	Close() error             // 关闭客户端, 释放资源

	Raw(sql string, dest interface{}) error // 查询SQL, 将结果存储到dest中
	Exec(sql string) error                  // 执行SQL, 不返回结果
}

// Pool 获取客户端池
func Pool() *typex.Enum[string, Client] {
	if !Initialized() {
		panic("client pool not initialized")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

// AddClient 添加客户端
func AddClient(source string, client Client) {
	if client == nil {
		return
	}
	if !Initialized() {
		pool = typex.NewStringEnum[Client]()
		pool.Add(constx.DefaultSource, client)
	}
	pool.Add(source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != "" {
		if client := Pool().Get(source[0]); client != nil {
			return client
		}
	}
	return Pool().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取实例
func GetInstance[T any](source ...string) T {
	instance, ok := GetClient(source...).GetInstance().(T)
	if !ok {
		panic("unexpected instance type")
	}
	return instance
}

// GetGormDB 获取数据库连接
func GetGormDB(source ...string) *gorm.DB {
	return GetInstance[*gorm.DB](source...)
}

// Close 关闭所有数据库客户端
func Close() error {
	var err error
	Pool().Range(func(source string, client Client) bool {
		if err = client.Close(); err != nil {
			err = errorx.Wrap(err, "close database client failed")
			return true
		}
		return false
	})
	return err
}
