package cachex

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const (
	ModeStandAlone = iota // 单机
	ModeCluster           // 集群
	ModeSentinel          // 哨兵
)

// RedisClientBuilder redis缓存客户端构建器
func RedisClientBuilder(config *Config) (Client, error) {
	return NewRedisClient(config)
}

// NewRedisClient 创建redis缓存客户端
func NewRedisClient(config *Config) (*RedisClient, error) {
	client, err := newRedisUniversalClient(config)
	if err != nil {
		return nil, errorx.Wrap(err, "new redis universal client failed")
	}
	return &RedisClient{
		config:  config,
		client:  client,
		marshal: marshalx.Apply(config.Marshal),
	}, nil
}

// 创建redis universal客户端
func newRedisUniversalClient(config *Config) (redis.UniversalClient, error) {
	opts := &redis.UniversalOptions{
		Addrs:      strings.Split(config.Address, ","),
		ClientName: config.Source,
		Username:   config.Username,
		Password:   config.Password,
		PoolSize:   config.PoolSize,
		DB:         config.Database,
	}
	var client redis.UniversalClient
	switch config.Mode {
	case ModeStandAlone:
		client = redis.NewClient(opts.Simple())
	case ModeCluster:
		client = redis.NewClusterClient(opts.Cluster())
	case ModeSentinel:
		opts.MasterName = config.Master
		client = redis.NewFailoverClient(opts.Failover())
	default:
		return nil, errors.New("redis mode is invalid")
	}
	if result, err := client.Ping(context.TODO()).Result(); err != nil || result != "PONG" {
		log.WithFields(config.LogFields()).
			WithField("ping_result", result).
			WithError(err).
			Error("redis client ping failed")
	}
	return client, nil
}

// RedisClient redis缓存客户端
type RedisClient struct {
	config  *Config
	marshal marshalx.Marshal
	client  redis.UniversalClient
}

func (c *RedisClient) GetClient() redis.UniversalClient {
	return c.client
}

func (c *RedisClient) GetConfig() *Config {
	return c.config
}

func (c *RedisClient) GetInstance() interface{} {
	return c.client
}

// GetKey 获取缓存键
func (c *RedisClient) GetKey(key string) string {
	return c.GetConfig().GetKey(key)
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if bytes, err := c.marshal.Marshal(value); err != nil {
		return errorx.Wrap(err, "marshal value error")
	} else if err = c.GetClient().Set(ctx, c.GetKey(key), bytes, expiration).Err(); err != nil {
		return errorx.Wrap(err, "set value error")
	}
	return nil
}

func (c *RedisClient) Get(ctx context.Context, key string, value any) bool {
	if result := c.GetString(ctx, key); result != "" {
		if err := c.marshal.Unmarshal([]byte(result), value); err == nil {
			return true
		}
	}
	return false

}

func (c *RedisClient) GetString(ctx context.Context, key string) string {
	if result, err := c.GetClient().Get(ctx, c.GetKey(key)).Result(); err == nil {
		return result
	}
	return ""
}

func (c *RedisClient) Delete(ctx context.Context, key string) bool {
	if result, err := c.GetClient().Del(ctx, c.GetKey(key)).Result(); err == nil {
		return result > 0
	}
	return false
}

func (c *RedisClient) Exist(ctx context.Context, key string) bool {
	if result, err := c.GetClient().Exists(ctx, c.GetKey(key)).Result(); err == nil {
		return result > 0
	}
	return false
}

func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if _, err := c.GetClient().Expire(ctx, c.GetKey(key), expiration).Result(); err != nil {
		return errorx.Wrap(err, "redis expire error")
	}
	return nil
}

func (c *RedisClient) Copy(source string, database int) (Client, error) {
	config := c.config.Copy()
	config.Source = source
	config.Database = database
	logger := log.WithFields(config.LogFields())
	client, err := NewRedisClient(config)
	if err != nil {
		logger.WithError(err).Error("new redis client failed")
		return nil, err
	}
	logger.Info("redis client copy success")
	return client, nil
}
