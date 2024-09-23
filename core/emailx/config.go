package emailx

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

// Config 邮件服务器配置
type Config struct {
	Host     string `json:"host" yaml:"host"`         // 邮件发送服务器地址
	Port     int    `json:"port" yaml:"port"`         // 邮件发送服务器端口
	Username string `json:"username" yaml:"username"` // 账户名
	Password string `json:"password" yaml:"password"` // 账号授权码
}

func (e *Config) newDialer() *gomail.Dialer {
	return gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
}

func (e *Config) Format() string {
	return fmtx.Yellow.XSPrintf("host=%s port=%v username=%s password=%s", e.Host, e.Port, e.Username, e.Password)
}

func (*Config) ID() string {
	return "email-server"
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mail.yaml",
		NacosDataId: "mail.yaml",
		Listen:      false,
	}
}

func (e *Config) Execute() error {
	handler = &Handler{config: e, dialer: e.newDialer()}
	log.Info("email-server init successfully: ", e.Format())
	return nil
}
