package captchax

import (
	"context"
	"github.com/go-xuan/quanx/os/syncx"
)

type Captcha interface {
	New(ctx context.Context) (id, image, answer string, err error) // 生成验证码
	Verify(ctx context.Context, id, answer string) bool            // 校验验证码
}

var clickCaptchaImpl *goCaptchaImpl
var ordinaryCaptchaImpl *base64CaptchaImpl

// 点击式验证码
func ClickCaptcha() Captcha {
	if clickCaptchaImpl == nil {
		syncx.OnceDo(func() {
			clickCaptchaImpl = newGoCaptcha()
		})
	}
	return clickCaptchaImpl
}

// 普通验证码
func OrdinaryCaptcha() Captcha {
	if ordinaryCaptchaImpl == nil {
		syncx.OnceDo(func() {
			ordinaryCaptchaImpl = newBase64Captcha()
		})
	}
	return ordinaryCaptchaImpl
}
