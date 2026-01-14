package ossx

import (
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/nacosx"
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
	fields["builder"] = c.Builder
	fields["endpoint"] = c.Endpoint
	fields["bucket"] = c.Bucket
	return fields
}

func (c *Config) Valid() bool {
	return c.Endpoint != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader(constx.OssConfigName),
		configx.NewFileReader(constx.OssConfigName),
		configx.NewTagReader(),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		logger := log.WithFields(c.LogFields())
		client, err := NewClient(c)
		if err != nil {
			logger.WithError(err).Error("create oss client failed")
			return errorx.Wrap(err, "create oss client failed")
		}
		logger.Info("init oss success")
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
		nacosx.NewReader(constx.OssConfigName),
		configx.NewFileReader(constx.OssConfigName),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "oss config execute failed")
		}
	}
	if !Initialized() {
		err := errorx.New("no enabled oss source")
		log.WithField("error", err.Error()).Warn("init oss failed")
		return err
	}
	return nil
}
