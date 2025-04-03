package emailx

import (
	"testing"

	"github.com/go-xuan/quanx/extra/configx"
)

func TestEmail(t *testing.T) {
	if err := configx.ReadAndExecute(&Config{}, configx.FromDefault); err != nil {
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
