package sendx

import (
	"time"

	"github.com/go-xuan/quanx/base/errorx"
	"gopkg.in/gomail.v2"
)

type Email struct {
	To      []string       // 收件人
	Cc      []string       // 抄送人
	Title   string         // 标题
	Content string         // 内容
	Dialer  *gomail.Dialer // 发件器
}

func (e *Email) AddReceiver(reciver ...string) Sender {
	e.To = append(e.To, reciver...)
	return e
}

func (e *Email) AddContent(content string) Sender {
	e.Content = content
	return e
}

func (e *Email) Send() error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", e.Dialer.Username)
	msg.SetHeader("To", e.To...)
	msg.SetHeader("Cc", e.Cc...)
	msg.SetHeader("Subject", e.Title)
	msg.SetDateHeader("X-Date", time.Now())
	msg.SetBody("text/plain", e.Content)
	if err := e.Dialer.DialAndSend(msg); err != nil {
		return errorx.Wrap(err, "send email error")
	}
	return nil
}
