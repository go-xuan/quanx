package redisx

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/types/enumx"
)

var pool *enumx.Enum[string, *Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

func this() *enumx.Enum[string, *Client] {
	if pool == nil {
		panic("redis client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, cli redis.UniversalClient) {
	if config == nil || cli == nil {
		return
	}
	client := &Client{config, cli}
	if pool == nil {
		pool = enumx.NewStringEnum[*Client]()
		pool.Add(constx.DefaultSource, client)
	}
	pool.Add(config.Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) redis.UniversalClient {
	return GetClient(source...).Instance()
}

// Ping 连接检查
func Ping(ctx context.Context, source ...string) (bool, error) {
	if result, err := GetInstance(source...).Ping(ctx).Result(); err != nil || result != "PONG" {
		return false, err
	} else {
		return true, nil
	}
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
