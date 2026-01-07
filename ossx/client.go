package ossx

import (
	"context"
	"io"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
)

func init() {
	RegisterBuilder("minio", MinioClientBuilder)
}

// Client oss客户端
type Client interface {
	GetConfig() *Config                                                                                    // 获取oss配置
	GetInstance() interface{}                                                                              // 获取数据库连接
	CreateBucket(ctx context.Context, name string, options ...interface{}) error                           // 创建bucket
	Upload(ctx context.Context, key string, reader io.Reader, options ...interface{}) error                // 上传文件
	Get(ctx context.Context, key string, options ...interface{}) (io.ReadCloser, error)                    // 获取文件
	Exist(ctx context.Context, key string, options ...interface{}) (bool, error)                           // 判断文件是否存在
	Remove(ctx context.Context, key string, options ...interface{}) error                                  // 删除文件
	GetUrl(ctx context.Context, key string, expires time.Duration, options ...interface{}) (string, error) // 获取文件url
}

// NewClient 创建oss客户端
func NewClient(config *Config) (Client, error) {
	client := stringx.IfZero(config.Client, "minio")
	if clientBuilders != nil {
		if builder, ok := clientBuilders[client]; ok {
			return builder(config)
		}
	}
	return nil, errorx.Sprintf("this oss client builder is not registered: %s", client)
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
