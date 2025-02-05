package elasticx

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
)

// Config ES配置
type Config struct {
	Host  string   `yaml:"host" json:"host"`   // 主机
	Port  int      `yaml:"port" json:"port"`   // 端口
	Index []string `yaml:"index" json:"index"` // 索引
}

func (e *Config) Format() string {
	return fmt.Sprintf("host=%s port=%v", e.Host, e.Port)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

func (e *Config) Execute() error {
	if e.Host != "" {
		var url = e.Url()
		if client, err := e.NewClient(url); err != nil {
			log.Error("elastic-search connect failed: ", e.Format(), err)
			return errorx.Wrap(err, "elasticx new client error")
		} else {
			_handler = &Handler{config: e, client: client}
			log.Info("elastic-search connect successfully: ", e.Format())
		}
	}
	return nil
}

func (e *Config) Url() string {
	return fmt.Sprintf("http://%s:%d", e.Host, e.Port)
}

func (e *Config) NewClient(url string) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, errorx.Wrap(err, "new elastic client failed")
	}
	var result *elastic.PingResult
	var code int
	if result, code, err = client.Ping(url).Do(context.Background()); err != nil || code != 200 {
		return nil, errorx.Wrap(err, "elastic ping failed")
	}
	log.Info("elastic-search version: ", result.Version.Number)
	return client, nil
}
