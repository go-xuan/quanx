package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/confx"
	"github.com/go-xuan/quanx/os/errorx"
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
func (e *Elastic) Run() error {
	if e.Host != "" {
		var url = e.Url()
		if client, err := e.NewClient(url); err != nil {
			log.Error("elastic-search connect failed: ", e.Info(), err)
			return errorx.Wrap(err, "elasticx.Elastic.NewClient error")
		} else {
			handler = &Handler{config: e, client: client}
			log.Info("elastic-search connect successful: ", e.Info())
		}
	}
	return nil
}

func (e *Elastic) Url() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Port)
}

func (e *Elastic) NewClient(url string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, errorx.Wrap(err, "elastic.NewClient Failed")
	}
	var result *elastic.PingResult
	var code int
	result, code, err = client.Ping(url).Do(context.Background())
	if err != nil || code != 200 {
		return nil, errorx.Wrap(err, "elastic.Client.Ping Failed")
	}
	log.Info("elastic-search Version: ", result.Version.Number)
	return client, nil

}
