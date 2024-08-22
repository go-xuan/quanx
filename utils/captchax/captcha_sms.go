package captchax

import (
	"context"
	"fmt"

	"github.com/go-xuan/quanx/utils/randx"
)

// NewSmsCaptcha 初始化短信验证码发送
func NewSmsCaptcha() *SmsCaptcha {
	return &SmsCaptcha{
		template: "",
		store:    DefaultStore(),
	}
}

type SmsCaptcha struct {
	template string
	store    *CaptchaStore
}

func (c *SmsCaptcha) Send(ctx context.Context, phone string) (captcha string, expired int, err error) {
	// 根据模板生成消息体
	captcha = randx.NumberCode(6)
	var message string
	if message, err = GetMessageByTemplate(c.template, captcha); err != nil {
		return
	}
	fmt.Println(message)

	// todo 调用短信发送接口

	// 存储验证码
	expired = c.store.expired
	if err = c.store.set(ctx, phone, captcha); err != nil {
		return
	}
	return
}

func (c *SmsCaptcha) Verify(ctx context.Context, phone, captcha string) bool {
	return c.store.verify(ctx, phone, captcha)
}
