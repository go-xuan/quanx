package captchax

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/utils/randx"
)

// ImageCaptchaService 图片验证码
type ImageCaptchaService interface {
	New(ctx context.Context) (id, image, captcha string, err error) // 生成验证码
	Verify(ctx context.Context, id, captcha string) bool            // 校验验证码
}

// CodeCaptchaService code验证码
type CodeCaptchaService interface {
	Send(ctx context.Context, receiver string) (captcha string, expired int, err error)
	Verify(ctx context.Context, receiver, captcha string) bool
}

func GetMessageByTemplate(text, captcha string) (string, error) {
	// 根据模板生成消息体
	tmpl, err := template.New("message").Parse(text)
	if err != nil {
		return "", errorx.Wrap(err, "template.New error")
	}
	captcha = randx.NumberCode(6)
	var data = map[string]string{"captcha": captcha}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "", errorx.Wrap(err, "template.Execute error")
	}
	return buf.String(), nil
}
