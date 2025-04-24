package captchax

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/randx"
	"github.com/go-xuan/quanx/utils/sendx"
)

// NewCodeCaptcha 初始化code验证码服务
func NewCodeCaptcha(sender sendx.Sender) CodeCaptchaService {
	return &CodeCaptcha{
		sender:   sender,
		template: defaultTemplate(),
		store:    DefaultStore(),
	}
}

func defaultTemplate() *template.Template {
	return template.Must(template.New("captcha").Parse("验证码：{{.captcha}}，请妥善保管，避免外泄。"))
}

type CodeCaptcha struct {
	len      int
	template *template.Template
	sender   sendx.Sender
	store    *CaptchaStore
}

func (c *CodeCaptcha) Send(ctx context.Context, reciver string) (string, int, error) {
	var captcha = randx.NumberCode(c.len)

	// 根据模板生成消息体
	var buf bytes.Buffer
	if err := c.template.Execute(&buf, map[string]string{
		"captcha": captcha,
	}); err != nil {
		return "", 0, errorx.Wrap(err, "template execute error")
	}

	// 发送验证码
	if err := c.sender.AddReceiver(reciver).AddContent(buf.String()).Send(); err != nil {
		return "", 0, errorx.Wrap(err, "send captcha error")
	}

	// 存储验证码
	if err := c.store.set(ctx, reciver, captcha); err != nil {
		return "", 0, errorx.Wrap(err, "store captcha error")
	}
	return captcha, c.store.expired, nil
}

func (c *CodeCaptcha) Verify(ctx context.Context, email, captcha string) bool {
	return c.store.verify(ctx, email, captcha)
}
