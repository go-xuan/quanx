package captchax

import (
	"context"
	"time"

	"github.com/mojocn/base64Captcha"

	"github.com/go-xuan/quanx/os/cachex"
	"github.com/go-xuan/quanx/types/stringx"
)

// 初始化
func newBase64Captcha() *base64CaptchaImpl {
	return &base64CaptchaImpl{
		capt: base64Captcha.NewCaptcha(
			&base64Captcha.DriverString{
				Height:          100,
				Width:           200,
				NoiseCount:      50,
				ShowLineOptions: 20,
				Length:          6,
				Source:          "abcdefghjkmnpqrstuvwxyz23456789",
			},
			&base64CaptchaStore{
				cache:    cachex.GetClient("captcha"),
				duration: 2 * time.Minute,
			},
		),
		clear: false,
	}
}

type base64CaptchaImpl struct {
	capt  *base64Captcha.Captcha
	clear bool
}

func (impl *base64CaptchaImpl) New(ctx context.Context) (id, image, answer string, err error) {
	return impl.capt.Generate()
}

func (impl *base64CaptchaImpl) Verify(ctx context.Context, id, answer string) bool {
	return impl.capt.Verify(id, stringx.ToLowerCamel(answer), impl.clear)
}

type base64CaptchaStore struct {
	cache    cachex.Client
	duration time.Duration
}

func (store *base64CaptchaStore) Set(id string, value string) error {
	return store.cache.Set(context.TODO(), id, value, store.duration)
}

func (store *base64CaptchaStore) Get(id string, clear bool) string {
	ctx := context.TODO()
	var value string
	store.cache.Get(ctx, id, &value)
	if clear {
		store.cache.Delete(ctx, id)
	}
	return value
}

func (store *base64CaptchaStore) Verify(id, answer string, clear bool) bool {
	if store.Get(id, clear) == answer {
		store.cache.Delete(context.TODO(), id)
		return true
	}
	return false
}
