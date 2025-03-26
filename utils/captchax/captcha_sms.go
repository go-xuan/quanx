package captchax

import (
	"context"
	"fmt"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/randx"
)

// NewSmsCaptcha 初始化短信验证码发送
func NewSmsCaptcha() CodeCaptchaService {
	return &SmsCaptcha{
		template: "",
		store:    DefaultStore(),
	}
}

type SmsCaptcha struct {
	template string
	store    *CaptchaStore
}

func (c *SmsCaptcha) Send(ctx context.Context, phone string) (string, int, error) {
	// 根据模板生成消息体
	captcha := randx.NumberCode(6)

	// 构建模板填充数据
	var data = make(map[string]string)
	data["captcha"] = captcha

	content, err := NewMessageByTemplate(c.template, data)
	if err != nil {
		return "", 0, errorx.Wrap(err, "new message content error")
	}
	fmt.Println(content)

	// todo 短信发送逻辑处理

	// 存储验证码
	expired := c.store.expired
	if err = c.store.set(ctx, phone, captcha); err != nil {
		return "", 0, errorx.Wrap(err, "store captcha error")
	}
	return content, expired, nil
}

func (c *SmsCaptcha) Verify(ctx context.Context, phone, captcha string) bool {
	return c.store.verify(ctx, phone, captcha)
}
