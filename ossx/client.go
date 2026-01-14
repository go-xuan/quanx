package ossx

import (
	"context"
	"io"
	"time"

	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
)

var (
	pool     *typex.Enum[string, Client]        // 客户端池
	builders *typex.Enum[string, ClientBuilder] // 客户端构造函数池
)

func init() {
	RegisterClientBuilder("minio", MinioClientBuilder) // 注册minio客户端构建器
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
		return MinioClientBuilder(config)
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

// Client oss客户端
type Client interface {
	GetInstance() interface{} // 获取数据库连接
	GetConfig() *Config       // 获取oss配置
	Close() error             // 关闭客户端

	CreateBucket(ctx context.Context, name string, options ...interface{}) error                           // 创建bucket
	Upload(ctx context.Context, key string, reader io.Reader, options ...interface{}) error                // 上传文件
	Get(ctx context.Context, key string, options ...interface{}) (io.ReadCloser, error)                    // 获取文件
	Exist(ctx context.Context, key string, options ...interface{}) (bool, error)                           // 判断文件是否存在
	Remove(ctx context.Context, key string, options ...interface{}) error                                  // 删除文件
	GetUrl(ctx context.Context, key string, expires time.Duration, options ...interface{}) (string, error) // 获取文件url
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

// Close 关闭所有客户端
func Close() error {
	var err error
	Pool().Range(func(source string, client Client) bool {
		if err = client.Close(); err != nil {
			err = errorx.Wrap(err, "close oss client failed")
			return true
		}
		return false
	})
	return err
}
