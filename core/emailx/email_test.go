package emailx

import (
	"testing"

	"github.com/go-xuan/quanx/core/configx"
)

func TestEmail(t *testing.T) {
	if err := configx.Execute(&Config{
		Host:     "smtp.qq.com",
		Port:     465,
		Username: "",
		Password: "",
	}); err != nil {
		panic(err)
	}

	if err := SendMail(Send{
		To:      []string{""},
		Cc:      []string{""},
		Title:   "测试邮件发送",
		Content: "测试邮件发送！\n测试邮件发送!\n测试邮件发送!",
	}); err != nil {
		panic(err)
	}
}
