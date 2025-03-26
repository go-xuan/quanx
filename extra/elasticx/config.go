package elasticx

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/anyx"
)

// Config ES配置
type Config struct {
	Source string   `json:"source" yaml:"source" default:"default"` // 数据源名称
	Enable bool     `json:"enable" yaml:"enable"`                   // 数据源启用
	Host   string   `yaml:"host" json:"host"`                       // 主机
	Port   int      `yaml:"port" json:"port"`                       // 端口
	Index  []string `yaml:"index" json:"index"`                     // 索引
}

func (c *Config) Format() string {
	return fmt.Sprintf("host=%s port=%v", c.Host, c.Port)
}

func (c *Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "elastic.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "elastic.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if client, err := c.NewClient(); err != nil {
			log.Error("elastic-search connect failed: ", c.Format())
			return errorx.Wrap(err, "new elasticx client error")
		} else {
			log.Info("elastic-search connect success:", c.Format())
			AddClient(c, client)
		}
	}
	return nil
}

func (c *Config) Url() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func (c *Config) NewClient() (*elastic.Client, error) {
	url := c.Url()
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, errorx.Wrap(err, "new elastic client failed")
	}
	var result *elastic.PingResult
	var code int
	if result, code, err = client.Ping(url).Do(context.Background()); err != nil || code != 200 {
		return nil, errorx.Wrap(err, "elastic-search ping failed")
	}
	log.Info("elastic-search version: ", result.Version.Number)
	return client, nil
}

type MultiConfig []*Config

func (list MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range list {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("{")
		sb.WriteString(config.Format())
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}

func (MultiConfig) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "elastic.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "elastic.yaml",
		}
	default:
		return nil
	}
}

func (list MultiConfig) Execute() error {
	if len(list) == 0 {
		return errorx.New("elastic-search not connected! cause: elastic.yaml is invalid")
	}
	for _, config := range list {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "elastic config execute error")
		}
	}
	if len(_handler.configs) == 0 {
		log.Error("elastic-search not connected! cause: no enabled source")
	}
	return nil
}
