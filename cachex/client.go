package cachex

import (
	"context"
	"time"

	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/constx"
)

var (
	pool     *typex.Enum[string, Client]        // 客户端池
	builders *typex.Enum[string, ClientBuilder] // 客户端构造函数池
)

func init() {
	RegisterClientBuilder("local", LocalClientBuilder) // 注册本地缓存客户端构建器
	RegisterClientBuilder("redis", RedisClientBuilder) // 注册redis缓存客户端构建器
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
		return LocalClientBuilder(config)
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

// Client 缓存客户端接口
type Client interface {
	GetInstance() interface{} // 获取实例
	GetConfig() *Config       // 获取配置
	Close() error             // 关闭客户端,

	GetKey(key string) string                                                       // 获取缓存key
	Set(ctx context.Context, key string, value any, expiration time.Duration) error // 更新缓存
	Get(ctx context.Context, key string, value any) bool                            // 获取缓存（指针，任意类型）
	GetString(ctx context.Context, key string) string                               // 获取缓存（字符串类型）
	Expire(ctx context.Context, key string, expiration time.Duration) error         // 续期缓存
	Delete(ctx context.Context, key string) bool                                    // 删除缓存
	Exist(ctx context.Context, key string) bool                                     // 是否存在缓存
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

// GetRedisUniversalClient 获取客户端
func GetRedisUniversalClient(source ...string) redis.UniversalClient {
	return GetInstance[redis.UniversalClient](source...)
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
