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
				Fonts:           []string{"chromohv.ttf"},
			},
			&CaptchaStore{
				cache:    cachex.GetClient("captcha"),
				duration: 2 * time.Minute,
			},
		),
		clear: true,
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

type CaptchaStore struct {
	cache    cachex.Client
	duration time.Duration
}

func (store *CaptchaStore) Set(id string, value string) error {
	return store.cache.Set(context.TODO(), id, value, store.duration)
}

func (store *CaptchaStore) Get(id string, clear bool) string {
	ctx := context.TODO()
	value := store.cache.GetString(ctx, id)
	if clear {
		store.cache.Delete(ctx, id)
	}
	return value
}

func (store *CaptchaStore) Verify(id, answer string, clear bool) bool {
	return store.Get(id, clear) == answer
}
