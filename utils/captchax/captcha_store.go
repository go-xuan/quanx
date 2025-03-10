package captchax

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/extra/cachex"
)

func DefaultStore() *CaptchaStore {
	return &CaptchaStore{
		client:  cachex.GetClient(),
		expired: 120,
		clear:   false,
	}
}

// CaptchaStore 验证码存储
type CaptchaStore struct {
	client  cachex.Client
	expired int
	clear   bool
}

// 存储
func (s *CaptchaStore) set(ctx context.Context, key, value string) error {
	return s.client.Set(ctx, key, value, time.Duration(s.expired)*time.Second)
}

// 获取
func (s *CaptchaStore) get(ctx context.Context, key string) string {
	var value string
	if ok := s.client.Get(ctx, key, &value); ok && s.clear {
		s.client.Delete(ctx, key)
	}
	return value
}

// 验证
func (s *CaptchaStore) verify(ctx context.Context, key, value string) bool {
	if s.get(ctx, key) == value {
		s.client.Delete(ctx, key)
		return true
	}
	return false
}

func (s *CaptchaStore) Set(key, value string) error {
	return s.set(context.TODO(), key, value)
}

func (s *CaptchaStore) Get(key string, clear bool) string {
	return s.get(context.TODO(), key)
}

func (s *CaptchaStore) Verify(key, value string, clear bool) bool {
	return s.verify(context.TODO(), key, value)
}
