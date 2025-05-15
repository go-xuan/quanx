package captchax

import (
	"context"
)

// ImageCaptchaService 图片验证码
type ImageCaptchaService interface {
	New(ctx context.Context) (id, image, captcha string, err error) // 生成验证码
	Verify(ctx context.Context, id, captcha string) bool            // 校验验证码
}

// CodeCaptchaService code验证码
type CodeCaptchaService interface {
	Send(ctx context.Context, receiver ...string) (captcha string, expired int, err error)
	Verify(ctx context.Context, receiver, captcha string) bool
}
