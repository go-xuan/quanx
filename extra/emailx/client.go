package emailx

import (
	"time"

	"gopkg.in/gomail.v2"

	"github.com/go-xuan/quanx/base/errorx"
)

var _client *Client

func this() *Client {
	if _client == nil {
		panic("email client not initialized, please check the relevant config")
	}
	return _client
}

type Client struct {
	config *Config
	dialer *gomail.Dialer
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) Instance() *gomail.Dialer {
	return c.dialer
}

func (c *Client) SendMail(send Send) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", c.config.Username)
	msg.SetHeader("To", send.To...)
	msg.SetHeader("Cc", send.Cc...)
	msg.SetHeader("Subject", send.Title)
	msg.SetDateHeader("X-Date", time.Now())
	msg.SetBody("text/plain", send.Content)
	// 请注意 DialAndSend() 方法是一次性的，也就是连接邮件服务器，发送邮件，然后关闭连接
	if err := c.Instance().DialAndSend(msg); err != nil {
		return errorx.Wrap(err, "send email error")
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

// SendMail 发送邮件
func SendMail(send Send) error {
	return this().SendMail(send)
}
