package emailx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/go-xuan/quanx/app/configx"
)

// Email 邮件服务器配置
type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件发送服务器地址
	Port     int    `json:"port" yaml:"port"`         // 邮件发送服务器端口
	Username string `json:"username" yaml:"username"` // 账户名
	Password string `json:"password" yaml:"password"` // 账号授权码
}

func (e *Email) newDialer() *gomail.Dialer {
	return gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
}

// Info 配置信息格式化
func (e *Email) Info() string {
	return fmt.Sprintf("host=%s port=%d username=%s password=%s", e.Host, e.Port, e.Username, e.Password)
}

// Title 配置器标题
func (*Email) Title() string {
	return "Email Server"
}

// Reader 配置文件读取
func (*Email) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mail.yaml",
		NacosDataId: "mail.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (e *Email) Run() error {
	handler = &Handler{config: e, dialer: e.newDialer()}
	log.Info("email-server init successful: ", e.Info())
	return nil
}
