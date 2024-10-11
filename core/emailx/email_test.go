package emailx

import (
	"testing"
)

func TestEmail(t *testing.T) {
	if err := NewConfigurator(&Config{
		Host:     "smtp.qq.com",
		Port:     465,
		Username: "",
		Password: "",
	}).Execute(); err != nil {
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
