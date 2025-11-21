package miniox

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/timex"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

type Config struct {
	Host         string `yaml:"host" json:"host"`                 // 主机
	Port         int    `yaml:"port" json:"port"`                 // 端口
	AccessId     string `yaml:"accessId" json:"accessId"`         // 访问id
	AccessSecret string `yaml:"accessSecret" json:"accessSecret"` // 访问秘钥
	SessionToken string `yaml:"sessionToken" json:"sessionToken"` // sessionToken
	Secure       bool   `yaml:"secure" json:"secure"`             // 是否使用https
	BucketName   string `yaml:"bucketName" json:"bucketName"`     // 桶名
	PrefixPath   string `yaml:"prefixPath" json:"prefixPath"`     // 前缀路径
	Expire       int64  `yaml:"expire" json:"expire"`             // 下载链接有效时长(分钟)
}

// LogEntry 日志打印实体类
func (c *Config) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"host":       c.Host,
		"port":       c.Port,
		"accessId":   c.AccessId,
		"bucketName": c.BucketName,
	})
}

func (c *Config) Valid() bool {
	return c.Host != "" && c.Port != 0
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("minio.yaml"),
		configx.NewFileReader("minio.yaml"),
	}
}

func (c *Config) Execute() error {
	if client, err := c.NewClient(); err != nil {
		c.LogEntry().WithError(err).Error("minio init failed")
		return errorx.Wrap(err, "new minio client failed")
	} else {
		_client = &Client{config: c, client: client}
		c.LogEntry().Info("minio init success")
		return nil
	}
}

func (c *Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// NewClient 初始化minio客户端
func (c *Config) NewClient() (*minio.Client, error) {
	if client, err := minio.New(c.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessId, c.AccessSecret, c.SessionToken),
		Secure: c.Secure,
		Region: Region,
	}); err != nil {
		return nil, errorx.Wrap(err, "new minio client failed")
	} else {
		return client, nil
	}
}

func (c *Config) MinioPath(fileName string) string {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	minioPath := filepath.Join(time.Now().Format(timex.TimestampFmt), uuid.NewString()+fileSuffix)
	if c.PrefixPath != "" {
		minioPath = filepath.Join(c.PrefixPath, minioPath)
	}
	return minioPath
}
