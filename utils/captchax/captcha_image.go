package captchax

import (
	"context"
	"github.com/mojocn/base64Captcha"

	"github.com/go-xuan/quanx/types/stringx"
)

// NewImageCaptcha 初始化图片验证码
func NewImageCaptcha() *ImageCaptcha {
	var driver = &base64Captcha.DriverString{
		Height:          100,
		Width:           200,
		NoiseCount:      50,
		ShowLineOptions: 20,
		Length:          6,
		Source:          "abcdefghjkmnpqrstuvwxyz23456789",
	}
	return &ImageCaptcha{
		capt: &base64Captcha.Captcha{
			Driver: driver,
			Store:  DefaultStore(),
		},
	}
}

type ImageCaptcha struct {
	capt *base64Captcha.Captcha
}

func (impl *ImageCaptcha) New(ctx context.Context) (id, image, answer string, err error) {
	return impl.capt.Generate()
}

func (impl *ImageCaptcha) Verify(ctx context.Context, id, answer string) bool {
	return impl.capt.Verify(id, stringx.ToLowerCamel(answer), false)
}
