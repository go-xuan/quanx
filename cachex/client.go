package cachex

import (
	"context"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

func init() {
	RegisterBuilder("redis", RedisClientBuilder) // 注册redis缓存客户端构建器
	RegisterBuilder("local", LocalClientBuilder) // 注册本地缓存客户端构建器
}

// Client 缓存客户端
type Client interface {
	GetConfig() *Config                                                             // 获取配置
	GetInstance() interface{}                                                       // 获取实例
	Copy(target string, database int) (Client, error)                               // 复制客户端
	GetKey(key string) string                                                       // 获取缓存key
	Set(ctx context.Context, key string, value any, expiration time.Duration) error // 更新缓存
	Get(ctx context.Context, key string, value any) bool                            // 获取缓存（指针，任意类型）
	GetString(ctx context.Context, key string) string                               // 获取缓存（字符串类型）
	Expire(ctx context.Context, key string, expiration time.Duration) error         // 续期缓存
	Delete(ctx context.Context, key string) bool                                    // 删除缓存
	Exist(ctx context.Context, key string) bool                                     // 是否存在缓存
}

// NewClient 创建oss客户端
func NewClient(config *Config) (Client, error) {
	client := stringx.IfZero(config.Client, "local")
	if clientBuilders != nil {
		if builder, ok := clientBuilders[client]; ok {
			return builder(config)
		}
	}
	return nil, errorx.Sprintf("the cache client builder is not registered: %s", client)
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
