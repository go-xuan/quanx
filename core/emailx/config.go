package emailx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/go-xuan/quanx/core/configx"
)

// Config 邮件服务器配置
type Config struct {
	Host     string `json:"host" yaml:"host"`         // 邮件发送服务器地址
	Port     int    `json:"port" yaml:"port"`         // 邮件发送服务器端口
	Username string `json:"username" yaml:"username"` // 账户名
	Password string `json:"password" yaml:"password"` // 账号授权码
}

func (c *Config) Format() string {
	return fmt.Sprintf("host=%s port=%v username=%s password=%s", c.Host, c.Port, c.Username, c.Password)
}

func (c *Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mail.yaml",
		NacosDataId: "mail.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	handler = &Handler{
		config: c,
		dialer: gomail.NewDialer(c.Host, c.Port, c.Username, c.Password),
	}
	log.Info("email-server init success: ", c.Format())
	return nil
}
