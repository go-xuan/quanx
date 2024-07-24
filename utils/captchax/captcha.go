package captchax

import (
	"context"
	"github.com/go-xuan/quanx/os/syncx"
)

type Captcha interface {
	New(ctx context.Context) (id, image, answer string, err error) // 生成验证码
	Verify(ctx context.Context, id, answer string) bool            // 校验验证码
}

var iGoCaptcha *goCaptchaImpl
var iBase64Captcha *base64CaptchaImpl

// GoCaptcha 点击式验证码
func GoCaptcha() Captcha {
	if iGoCaptcha == nil {
		syncx.OnceDo(func() {
			iGoCaptcha = newGoCaptcha()
		})
	}
	return iGoCaptcha
}

// Base64Captcha 普通验证码
func Base64Captcha() Captcha {
	if iBase64Captcha == nil {
		syncx.OnceDo(func() {
			iBase64Captcha = newBase64Captcha()
		})
	}
	return iBase64Captcha
}
