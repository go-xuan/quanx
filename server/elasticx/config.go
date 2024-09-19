package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

type Elastic struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

func (*Elastic) ID() string {
	return "elastic-search"
}

func (e *Elastic) Format() string {
	return fmtx.Yellow.XSPrintf("host=%s port=%v", e.Host, e.Port)
}

func (*Elastic) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

func (e *Elastic) Execute() error {
	if e.Host != "" {
		var url = e.Url()
		if client, err := e.NewClient(url); err != nil {
			log.Error("elastic-search connect failed: ", e.Format(), err)
			return errorx.Wrap(err, "elasticx.Elastic.NewClient error")
		} else {
			handler = &Handler{config: e, client: client}
			log.Info("elastic-search connect successfully: ", e.Format())
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
	log.Info("elastic-search version: ", result.Version.Number)
	return client, nil

}
