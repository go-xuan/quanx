package emailx

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"

	"github.com/go-xuan/quanx/frame/confx"
)

// 邮件服务器配置
type Email struct {
	Host     string `json:"host" yaml:"host"`         // 邮件发送服务器地址
	Port     int    `json:"port" yaml:"port"`         // 邮件发送服务器端口
	Username string `json:"username" yaml:"username"` // 账户名
	Password string `json:"password" yaml:"password"` // 账号授权码
}

func (e *Email) newDialer() *gomail.Dialer {
	return gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)
}

// 配置信息格式化
func (e *Email) ToString(title string) string {
	return fmt.Sprintf("%s => host=%s port=%d Account=%s password=%s", title, e.Host, e.Port, e.Username, e.Password)
}

// 配置器名称
func (*Email) Title() string {
	return "init mail-server"
}

// 配置文件读取
func (*Email) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "mail.yaml",
		NacosDataId: "mail.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (e *Email) Run() error {
	handler = &Handler{Config: e, Dialer: e.newDialer()}
	log.Info(e.ToString("mail-server init successful!"))
	return nil
}
