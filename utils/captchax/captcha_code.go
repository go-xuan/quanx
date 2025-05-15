package captchax

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/randx"
)

// NewCodeCaptcha 初始化code验证码服务
func NewCodeCaptcha(sender Sender) CodeCaptchaService {
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
	len      int                // 验证码长度
	template *template.Template // 验证码模板
	store    *CaptchaStore      //	验证码存储
	sender   Sender             //	验证码发送器
}

func (c *CodeCaptcha) Send(ctx context.Context, reciver ...string) (string, int, error) {
	if len(reciver) == 0 || reciver[0] == "" {
		return "", 0, errorx.New("reciver is empty")
	}

	// 发送验证码
	var buf bytes.Buffer
	var captcha = randx.NumberCode(c.len)
	if err := c.template.Execute(&buf, map[string]string{"captcha": captcha}); err != nil {
		return "", 0, errorx.Wrap(err, "template execute error")
	} else if err = c.sender.AddReceiver(reciver...).SetContent(buf.String()).Send(); err != nil {
		return "", 0, errorx.Wrap(err, "send captcha error")
	}

	// 存储验证码
	for _, key := range reciver {
		if err := c.store.set(ctx, key, captcha); err != nil {
			return "", 0, errorx.Wrap(err, "store captcha error")
		}
	}
	return captcha, c.store.expired, nil
}

func (c *CodeCaptcha) Verify(ctx context.Context, email, captcha string) bool {
	return c.store.verify(ctx, email, captcha)
}

// Sender 验证码发送器，由业务自行实现
type Sender interface {
	AddReceiver(reciver ...string) Sender
	SetContent(content string) Sender
	Send() error
}
