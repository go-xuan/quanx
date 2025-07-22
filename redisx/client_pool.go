package redisx

import (
	"errors"

	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"github.com/redis/go-redis/v9"
)

var pool *typex.Enum[string, *Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil
}

func this() *typex.Enum[string, *Client] {
	if !Initialized() {
		panic("redis client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, cli redis.UniversalClient) {
	if config == nil || cli == nil {
		return
	}
	client := &Client{config, cli}
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		this().Add("default", client)
	}
	this().Add(config.Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get("default")
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) redis.UniversalClient {
	return GetClient(source...).Instance()
}

// CopyDatabase 复制redis数据库
func CopyDatabase(source, target string, database int) error {
	if source != "" && target != "" {
		if sourceClient := this().Get(source); sourceClient != nil {
			if targetClient, err := sourceClient.Copy(target, database); err == nil {
				this().Add(target, targetClient)
			} else {
				return errorx.Wrap(err, "copy redis client failed")
			}
		} else {
			return errors.New("redis source client not exist")
		}
	}
	return nil
}
