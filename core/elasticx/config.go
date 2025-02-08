package elasticx

import (
	"context"
	"fmt"
	"strings"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
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

func (c *Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if client, err := c.NewClient(); err != nil {
			return errorx.Wrap(err, "elasticx new client error")
		} else {
			log.Info("elastic-search connect success: ", c.Format())
			if _handler == nil {
				_handler = &Handler{
					multi: false, config: c, client: client,
					configs: make(map[string]*Config),
					clients: make(map[string]*elastic.Client),
				}
			} else {
				_handler.multi = true
				if c.Source == constx.DefaultSource {
					_handler.config = c
					_handler.client = client
				}
			}
			_handler.configs[c.Source] = c
			_handler.clients[c.Source] = client
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
		return nil, errorx.Wrap(err, "elastic ping failed")
	}
	log.Info("elastic-search version: ", result.Version.Number)
	return client, nil
}

type MultiConfig []*Config

func (m MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range m {
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

func (MultiConfig) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "elastic.yaml",
		NacosDataId: "elastic.yaml",
		Listen:      false,
	}
}

func (m MultiConfig) Execute() error {
	if len(m) == 0 {
		return errorx.New("elastic not connected! cause: elastic.yaml is invalid")
	}
	if _handler == nil {
		_handler = &Handler{
			configs: make(map[string]*Config),
			clients: make(map[string]*elastic.Client),
		}
	}
	_handler.multi = true
	for i, c := range m {
		if c.Enable {
			if err := anyx.SetDefaultValue(c); err != nil {
				return errorx.Wrap(err, "set default value error")
			}
			if client, err := c.NewClient(); err != nil {
				return errorx.Wrap(err, "new elastic client failed")
			} else {
				_handler.clients[c.Source] = client
				_handler.configs[c.Source] = c
				if i == 0 || c.Source == constx.DefaultSource {
					_handler.client = client
					_handler.config = c
				}
			}
		}
	}
	if len(_handler.configs) == 0 {
		log.Error("elastic not connected! cause: elastic.yaml is empty or no enabled source")
	}
	return nil
}
