package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d", conf.Host, conf.Port)
}

func (conf *Config) Init() {
	var url = conf.Url()
	client, err := conf.NewClient(url)
	if err != nil {
		log.Error("ElasticSearch连接失败！", conf.Format())
		log.Error("error : ", err)
	} else {
		instance = &Handler{Config: conf, Url: url, Client: client}
		log.Error("ElasticSearch连接成功！", conf.Format())
	}
}

func (conf *Config) Url() string {
	return fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
}

func (conf *Config) NewClient(url string) (client *elastic.Client, err error) {
	client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return
	}
	var result *elastic.PingResult
	var code int
	if result, code, err = client.Ping(url).Do(context.Background()); err != nil && code != 200 {
		return
	}
	log.Info("ElasticSearch 版本 : ", result.Version.Number)
	return
}
