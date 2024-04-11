package redisx

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/go-xuan/quanx/utils/slicex"
)

type CacheService struct {
	Source string
	Prefix string
	Case   *marshalx.Case
}

func (s *CacheService) DB() redis.Cmdable {
	return This().GetCmd(s.Source)
}

func (s *CacheService) Get(ctx context.Context, key string, v any) error {
	if value, err := s.DB().Get(ctx, s.Prefix+key).Bytes(); err != nil {
		return err
	} else {
		return s.Case.Unmarshal(value, v)
	}
}

func (s *CacheService) Set(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := s.Case.Marshal(v); err != nil {
		return err
	} else {
		return s.DB().Set(ctx, s.Prefix+key, value, expiration).Err()
	}
}

func (s *CacheService) SetNX(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := s.Case.Marshal(v); err != nil {
		return err
	} else {
		return s.DB().SetNX(ctx, s.Prefix+key, value, expiration).Err()
	}
}

func (s *CacheService) Exists(ctx context.Context, keys ...string) (total int64, err error) {
	if l := len(keys); l > 0 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var batches = s.AddPrefix(keys[x:y])
			var n int64
			if n, err = s.DB().Exists(ctx, batches...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return
		}
	}
	return
}

func (s *CacheService) Delete(ctx context.Context, keys ...string) (total int64, err error) {
	if l := len(keys); l > 0 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var batches = s.AddPrefix(keys[x:y])
			var n int64
			if n, err = s.DB().Del(ctx, batches...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return
		}
	}
	return
}

func (s *CacheService) AddPrefix(keys []string) []string {
	var newKeys []string
	if len(keys) > 0 {
		for _, key := range keys {
			newKeys = append(newKeys, s.Prefix+key)
		}
	}
	return newKeys
}
