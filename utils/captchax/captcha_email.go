package captchax

import (
	"context"

	"github.com/go-xuan/quanx/core/emailx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/utils/randx"
)

// NewEmailCaptcha 初始化邮箱验证码发送
func NewEmailCaptcha() *EmailCaptcha {
	return &EmailCaptcha{
		title:    "",
		template: "",
		handler:  emailx.This(),
		store:    DefaultStore(),
	}
}

type EmailCaptcha struct {
	title    string
	template string
	handler  *emailx.Handler
	store    *CaptchaStore
}

func (c *EmailCaptcha) Send(ctx context.Context, reciver string) (string, int, error) {
	// 根据模板生成消息体
	captcha := randx.NumberCode(6)

	// 构建模板填充数据
	var data = make(map[string]string)
	data["captcha"] = captcha

	// 生成message内容
	content, err := NewMessageByTemplate(c.template, data)
	if err != nil {
		return "", 0, errorx.Wrap(err, "new message content error")
	}
	// 发送邮箱验证码
	if err = c.handler.SendMail(emailx.Send{
		To:      []string{reciver},
		Title:   c.title,
		Content: content,
	}); err != nil {
		return "", 0, errorx.Wrap(err, "send captcha mail error")
	}

	// 存储验证码
	expired := c.store.expired
	if err = c.store.set(ctx, reciver, captcha); err != nil {
		return "", 0, errorx.Wrap(err, "store captcha error")
	}
	return captcha, expired, nil
}

func (c *EmailCaptcha) Verify(ctx context.Context, email, captcha string) bool {
	return c.store.verify(ctx, email, captcha)
}
