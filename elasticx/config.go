package elasticx

import (
	"context"
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/httpx"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// Config ES配置
type Config struct {
	Source   string   `json:"source" yaml:"source" default:"default"` // 数据源名称
	Enable   bool     `json:"enable" yaml:"enable"`                   // 数据源启用
	Host     string   `json:"host" yaml:"host"`                       // 主机
	Port     int      `json:"port" yaml:"port"`                       // 端口
	Username string   `json:"username" yaml:"username"`               // 用户名
	Password string   `json:"password" yaml:"password"`               // 密码
	Index    []string `json:"index" yaml:"index"`                     // 索引
}

// LogEntry 日志打印实体类
func (c *Config) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"source": c.Source,
		"host":   c.Host,
		"port":   c.Port,
	})
}

func (c *Config) Valid() bool {
	return c.Host != "" && c.Port != 0
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("elastic.yaml"),
		configx.NewFileReader("elastic.yaml"),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if client, err := c.NewClient(); err != nil {
			c.LogEntry().WithError(err).Error("elastic-search init failed")
			return errorx.Wrap(err, "init elasticx client error")
		} else {
			c.LogEntry().Info("elastic-search init success")
			AddClient(c, client)
		}
	}
	return nil
}

// GetUrl 构建ES连接URL
func (c *Config) GetUrl() string {
	protocol, host, port := httpx.ParseHost(c.Host)
	if port == 0 && c.Port > 0 {
		port = c.Port
	}
	if protocol != "" {
		return fmt.Sprintf("%s://%s:%d", protocol, host, port)
	} else {
		return fmt.Sprintf("%s:%d", host, port)
	}
}

func (c *Config) NewClient() (*elastic.Client, error) {
	url := c.GetUrl()
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(c.Username, c.Password),
		elastic.SetHttpClient(httpx.NewClient()),
	)
	if err != nil {
		return nil, errorx.Wrap(err, "new elastic client failed")
	}
	result, code, err := client.Ping(url).Do(context.Background())
	if err != nil || code != 200 {
		return nil, errorx.Wrap(err, "elastic-search ping failed")
	}
	log.Info("elastic-search version: ", result.Version.Number)
	return client, nil
}

type Configs []*Config

func (s Configs) Valid() bool {
	return len(s) > 0
}

func (s Configs) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("elastic.yaml"),
		configx.NewFileReader("elastic.yaml"),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "elastic config execute error")
		}
	}
	if !Initialized() {
		log.Error("elastic-search not initialized because no enabled source")
	}
	return nil
}
