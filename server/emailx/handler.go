package emailx

import (
	"time"

	"gopkg.in/gomail.v2"
)

var handler *Handler

type Handler struct {
	config *Email
	dialer *gomail.Dialer
}

func This() *Handler {
	if handler == nil {
		panic("the mail handler has not been initialized, please check the relevant config")
	}
	return handler
}

func (h *Handler) GetConfig() *Email {
	return h.config
}

// SendMail 发送邮件
func (h *Handler) SendMail(send Send) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", h.config.Username)
	msg.SetHeader("To", send.To...)
	msg.SetHeader("Cc", send.Cc...)
	msg.SetHeader("Subject", send.Title)
	msg.SetDateHeader("X-Date", time.Now())
	msg.SetBody("text/plain", send.Content)
	// 请注意 DialAndSend() 方法是一次性的，也就是连接邮件服务器，发送邮件，然后关闭连接
	if err := h.dialer.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

// Send 邮件服务器发送配置
type Send struct {
	To      []string `json:"to"`      // 收件人
	Cc      []string `json:"cc"`      // 抄送人
	Title   string   `json:"title"`   // 标题
	Content string   `json:"content"` // 内容
}
