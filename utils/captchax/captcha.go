package captchax

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/types/stringx"
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

// NewMessageByTemplate 根据模板生成消息
func NewMessageByTemplate(templateText string, data map[string]string) (string, error) {
	// 默认消息模板
	templateText = stringx.IfZero(templateText, `验证码：{{.captcha}}，请妥善保管，避免外泄。`)

	// 根据模板生成消息内容
	if msgTemplate, err := template.New("message").Parse(templateText); err != nil {
		return "", errorx.Wrap(err, "new message template error")
	} else {
		var buf bytes.Buffer
		if err = msgTemplate.Execute(&buf, data); err != nil {
			return "", errorx.Wrap(err, "template.Execute error")
		}
		return buf.String(), nil
	}
}
