package cachex

import (
	"context"
	"time"
)

// Client 缓存客户端
type Client interface {
	Config() *Config                                                                // 获取配置
	Set(ctx context.Context, key string, value any, expiration time.Duration) error // 更新缓存
	Get(ctx context.Context, key string, value any) bool                            // 获取缓存（指针，任意类型）
	GetString(ctx context.Context, key string) string                               // 获取缓存（字符串类型）
	Expire(ctx context.Context, key string, expiration time.Duration) error         // 续期
	Delete(ctx context.Context, key string) bool                                    // 删除
	Exist(ctx context.Context, key string) bool                                     // 是否存在
}
