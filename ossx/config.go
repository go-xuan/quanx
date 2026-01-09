package ossx

import (
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Source          string `json:"source" yaml:"source"`                   // oss源名称
	Builder         string `json:"builder" yaml:"builder" default:"minio"` // 客户端选型
	Enable          bool   `json:"enable" yaml:"enable"`                   // 启用
	Endpoint        string `json:"endpoint" yaml:"endpoint"`               // 主机
	AccessKeyId     string `json:"accessKeyId" yaml:"accessKeyId"`         // 访问id
	AccessKeySecret string `json:"accessKeySecret" yaml:"accessKeySecret"` // 访问秘钥
	AccessToken     string `json:"accessToken" json:"accessToken"`         // 访问token
	Secure          bool   `json:"secure" yaml:"secure"`                   // 是否使用https
	Bucket          string `json:"bucket" yaml:"bucket"`                   // 桶名
}

func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["client"] = c.Builder
	fields["endpoint"] = c.Endpoint
	fields["bucket"] = c.Bucket
	return fields
}

func (c *Config) Valid() bool {
	return c.Endpoint != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("oss.yaml"),
		configx.NewFileReader("oss.yaml"),
		configx.NewTagReader(),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		client, err := NewClient(c)
		if err != nil {
			log.WithFields(c.LogFields()).WithError(err).Error("init oss client failed")
			return errorx.Wrap(err, "init oss client failed")
		}
		log.WithFields(c.LogFields()).Info("init oss client success")
		AddClient(c.Source, client)
	}
	return nil
}

type Configs []*Config

func (s Configs) Valid() bool {
	return len(s) > 0
}

func (s Configs) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("oss.yaml"),
		configx.NewFileReader("oss.yaml"),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "oss config execute error")
		}
	}
	if !Initialized() {
		log.Error("oss not initialized because no enabled source")
	}
	return nil
}
