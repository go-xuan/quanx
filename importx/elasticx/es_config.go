package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
)

type Elastic struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

// 配置信息格式化
func (e *Elastic) ToString(title string) string {
	return fmt.Sprintf("%s => host=%s port=%d", title, e.Host, e.Port)
}

// 配置器名称
func (*Elastic) Title() string {
	return "init elastic-search"
}

func (*Elastic) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (e *Elastic) Run() error {
	if e.Host == "" {
		return nil
	}
	var url = e.Url()
	client, err := e.NewClient(url)
	if err != nil {
		log.Error(e.ToString("elastic search connect failed!"))
		log.Error("error : ", err)
		return err
	}
	handler = &Handler{Config: e, Url: url, Client: client}
	log.Error(e.ToString("elastic search connect successful!"))
	return nil
}

func (e *Elastic) Url() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Port)
}

func (e *Elastic) NewClient(url string) (client *elastic.Client, err error) {
	client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return
	}
	var result *elastic.PingResult
	var code int
	if result, code, err = client.Ping(url).Do(context.Background()); err != nil && code != 200 {
		return
	}
	log.Info("elastic search version : ", result.Version.Number)
	return
}
