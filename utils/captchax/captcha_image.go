package captchax

import (
	"context"
	"strings"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/mojocn/base64Captcha"
)

// NewImageCaptcha 初始化图片验证码
func NewImageCaptcha() ImageCaptchaService {
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

func (impl *ImageCaptcha) New(ctx context.Context) (string, string, string, error) {
	if id, image, answer, err := impl.capt.Generate(); err != nil {
		return "", "", "", errorx.Wrap(err, "generate captcha failed")
	} else {
		return id, image, answer, nil
	}
}

func (impl *ImageCaptcha) Verify(ctx context.Context, id, answer string) bool {
	return impl.capt.Verify(id, strings.ToLower(answer), false)
}
