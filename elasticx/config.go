package elasticx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// Config ES配置
type Config struct {
	Source   string   `json:"source" yaml:"source" default:"default"` // 数据源名称
	Enable   bool     `json:"enable" yaml:"enable"`                   // 数据源启用
	Url      string   `json:"url" yaml:"url"`                         // 地址
	Username string   `json:"username" yaml:"username"`               // 用户名
	Password string   `json:"password" yaml:"password"`               // 密码
	Indices  []string `json:"indices" yaml:"indices"`                 // 索引
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["address"] = c.Url
	return fields
}

func (c *Config) Valid() bool {
	return c.Url != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("elastic.yaml"),
		configx.NewFileReader("elastic.yaml"),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		logger := log.WithFields(c.LogFields())
		client, err := NewClient(c)
		if err != nil {
			logger.WithError(err).Error("init elastic-search client failed")
			return errorx.Wrap(err, "init elastic-search client failed")
		}
		logger.Info("init elastic-search client success")
		AddClient(c.Source, client)
		ctx := context.Background()
		for _, index := range c.Indices {
			if _, err = client.CreateIndex(ctx, index); err != nil {
				logger.WithError(err).Errorf("create index %s failed", index)
				return errorx.Wrap(err, "create index failed")
			}
		}
	}
	return nil
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
		log.Error("elastic not initialized because no enabled source")
	}
	return nil
}
