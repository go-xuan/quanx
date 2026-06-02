package ossx

import (
	"context"
	"io"
	"time"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
)

var (
	pool     = configx.NewPool[Client]()        // 客户端池
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
	GetInstance() any // 获取数据库连接
	GetConfig() *Config       // 获取oss配置
	Close() error             // 关闭客户端

	CreateBucket(ctx context.Context, name string, options ...any) error                           // 创建bucket
	Upload(ctx context.Context, key string, reader io.Reader, options ...any) error                // 上传文件
	Get(ctx context.Context, key string, options ...any) (io.ReadCloser, error)                    // 获取文件
	Download(ctx context.Context, key string, options ...any) error                                // 获取文件
	Exist(ctx context.Context, key string, options ...any) (bool, error)                           // 判断文件是否存在
	Remove(ctx context.Context, key string, options ...any) error                                  // 删除文件
	GetUrl(ctx context.Context, key string, expires time.Duration, options ...any) (string, error) // 获取文件url
}

// Pool 获取客户端池
func Pool() *configx.Pool[Client] {
	if !Initialized() {
		panic("client pool not initialized")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool.Initialized()
}

// AddClient 添加客户端
func AddClient(source string, client Client) {
	if any(client) == nil {
		return
	}
	pool.Add(source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	return pool.Get(source...)
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
	return pool.Close(func(client Client) error {
		return client.Close()
	})
}
