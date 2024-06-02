package emailx

import (
	"gopkg.in/gomail.v2"
)

var handler *Handler

type Handler struct {
	Config *Email
	Dialer *gomail.Dialer
}

func This() *Handler {
	if handler == nil {
		panic("The mail handler has not been initialized, please check the relevant config")
	}
	return handler
}

// 发送邮件
func (h *Handler) SendMail(send *Send) error {
	if err := h.Dialer.DialAndSend(send.newMessage()); err != nil {
		return err
	}
	return nil
}

// 邮件服务器发送配置
type Send struct {
	From    string   `json:"from"`    // 发件人
	To      []string `json:"to"`      // 收件人
	Cc      []string `json:"cc"`      // 抄送人
	Title   string   `json:"title"`   // 标题
	Content string   `json:"content"` // 内容
}

// 构建邮件
func (s *Send) newMessage() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("ID", s.From)
	msg.SetHeader("To", s.To...)
	msg.SetHeader("Cc", s.Cc...)
	msg.SetHeader("Subject", s.Title)
	msg.SetBody("text/plain", s.Content)
	return msg
}
