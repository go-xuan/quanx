package mongox

import (
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"go.mongodb.org/mongo-driver/mongo"
)

// 客户端池
var pool *typex.Enum[string, *Client]

// Pool 获取客户端池
func Pool() *typex.Enum[string, *Client] {
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
func AddClient(source string, client *Client) {
	if client == nil {
		return
	}
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		pool.Add("default", client)
	}
	pool.Add(source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := Pool().Get(source[0]); client != nil {
			return client
		}
	}
	return Pool().Get("default")
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *mongo.Client {
	return GetClient(source...).GetClient()
}

// GetDatabase 获取数据库
func GetDatabase(source ...string) *mongo.Database {
	if client := GetClient(source...); client != nil {
		database := client.GetConfig().Database
		return client.GetClient().Database(database)
	}
	return nil
}

// Close 关闭所有数据库客户端
func Close() error {
	var err error
	Pool().Range(func(_ string, client *Client) bool {
		if err = client.Close(); err != nil {
			err = errorx.Wrap(err, "close mongo client failed")
			return true
		}
		return false
	})
	return err
}
