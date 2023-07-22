package elasticx

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Config struct {
	Host string `yaml:"host" json:"host"` // 主机
	Port int    `yaml:"port" json:"port"` // 端口
}

func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d", conf.Host, conf.Port)
}

func (conf *Config) Url() string {
	return fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
}

func (conf *Config) NewEsCtl() (*Control, error) {
	client, err := conf.NewClient()
	if err != nil {
		return nil, err
	}
	return &Control{
		Config: conf,
		Url:    conf.Url(),
		Client: client,
		Ctx:    context.Background(),
	}, nil
}

func (conf *Config) NewClient() (client *elastic.Client, err error) {
	client, err = elastic.NewClient(
		elastic.SetURL(conf.Url()), // 设置url
		elastic.SetSniff(false),    // 关闭嗅探
	)
	return
}
