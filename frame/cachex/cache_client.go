package cachex

import (
	"context"
	"time"
)

type CacheClient struct {
	Set   func(context.Context, string, any, time.Duration) // 设置缓存值
	Get   func(context.Context, string) any                 // 获取缓存值
	Del   func(context.Context, string)                     // 删除缓存值
	Exist func(context.Context, string) bool                // 是否存在
}

func GetValue(ctx context.Context, client CacheClient, key string, get func() any) any {
	if client.Exist(ctx, key) {
		return client.Get(ctx, key)
	} else {
		v := get()
		client.Set(ctx, key, v, time.Duration(0))
		return v
	}
}
