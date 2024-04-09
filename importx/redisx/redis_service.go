package redisx

import (
	"context"
	"github.com/go-xuan/quanx/utilx/slicex"
	"time"

	"github.com/go-xuan/quanx/importx/marshalx"
)

type RedisService struct {
	Source      string
	Prefix      string
	MarshalType string
}

func (s *RedisService) Get(ctx context.Context, key string, v any) error {
	if value, err := DB(s.Source).Get(ctx, s.Prefix+key).Bytes(); err != nil {
		return err
	} else {
		return marshalx.Unmarshal(s.MarshalType, value, v)
	}
}

func (s *RedisService) Set(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := marshalx.Marshal(s.MarshalType, v); err != nil {
		return err
	} else {
		return DB(s.Source).Set(ctx, s.Prefix+key, value, expiration).Err()
	}
}

func (s *RedisService) SetNX(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := marshalx.Marshal(s.MarshalType, v); err != nil {
		return err
	} else {
		return DB(s.Source).SetNX(ctx, s.Prefix+key, value, expiration).Err()
	}
}

func (s *RedisService) Exists(ctx context.Context, keys ...string) (total int64, err error) {
	if err = slicex.ExecInBatches(len(keys), 100, func(x int, y int) (err error) {
		var batches = s.AddPrefix(keys[x:y]...)
		var n int64
		if n, err = DB(s.Source).Exists(ctx, batches...).Result(); err != nil {
			return
		}
		total += n
		return
	}); err != nil {
		return
	}
	return
}

func (s *RedisService) Delete(ctx context.Context, keys ...string) (total int64, err error) {
	if err = slicex.ExecInBatches(len(keys), 100, func(x int, y int) (err error) {
		var batches = s.AddPrefix(keys[x:y]...)
		var n int64
		if n, err = DB(s.Source).Del(ctx, batches...).Result(); err != nil {
			return
		}
		total += n
		return
	}); err != nil {
		return
	}
	return
}

func (s *RedisService) AddPrefix(keys ...string) []string {
	var newKeys []string
	if len(keys) > 0 {
		for _, key := range keys {
			newKeys = append(newKeys, s.Prefix+key)
		}
	}
	return newKeys
}
