package emailx

import (
	"github.com/go-xuan/quanx/core/configx"
	"testing"
)

func TestEmail(t *testing.T) {
	if err := configx.Execute(&Config{
		Host:     "smtp.qq.com",
		Port:     465,
		Username: "quanchao1996@qq.com",
		Password: "zshgafvebiyfebdh",
	}); err != nil {
		panic(err)
	}

	if err := SendMail(Send{
		To:      []string{"2465836880@qq.com"},
		Cc:      []string{"quanchao@gmail.com"},
		Title:   "全超测试邮件发送",
		Content: "全超测试邮件发送！\n全超测试邮件发送!\n全超测试邮件发送!",
	}); err != nil {
		panic(err)
	}
}
