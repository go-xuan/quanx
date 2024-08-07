package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/confx"
)

type Elastic struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

// Info 配置信息格式化
func (e *Elastic) Info() string {
	return fmt.Sprintf("host=%s port=%d", e.Host, e.Port)
}

// Title 配置器标题
func (*Elastic) Title() string {
	return "ElasticSearch"
}

// Reader 配置文件读取
func (*Elastic) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (e *Elastic) Run() (err error) {
	if e.Host != "" {
		var url = e.Url()
		var client *elastic.Client
		if client, err = e.NewClient(url); err != nil {
			log.Error("Elastic Search Connect Failed: ", e.Info(), err)
			return
		}
		handler = &Handler{Config: e, Url: url, Client: client}
		log.Info("Elastic Search Connect Successful: ", e.Info())
	}
	return
}

func (e *Elastic) Url() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Port)
}

func (e *Elastic) NewClient(url string) (client *elastic.Client, err error) {
	if client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false)); err != nil {
		return
	}
	var result *elastic.PingResult
	var code int
	if result, code, err = client.Ping(url).Do(context.Background()); err != nil && code != 200 {
		return
	}
	log.Info("ElasticSearch Version: ", result.Version.Number)
	return
}
