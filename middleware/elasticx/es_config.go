package elasticx

import (
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Config struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d", conf.Host, conf.Port)
}

func (conf *Config) Url() string {
	return fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
}

func (conf *Config) NewClient(url string) (client *elastic.Client, err error) {
	client, err = elastic.NewClient(
		elastic.SetURL(url),     // 设置url
		elastic.SetSniff(false), // 关闭嗅探
	)
	return
}
