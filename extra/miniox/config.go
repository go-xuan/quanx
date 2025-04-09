package miniox

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/timex"
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

func (m *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d accessId=%s bucketName=%s", m.Host, m.Port, m.AccessId, m.BucketName)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "minio.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalReader{
			Name: "minio.yaml",
		}
	default:
		return nil
	}
}

func (m *Config) Execute() error {
	if client, err := m.NewClient(); err != nil {
		log.Error("minio connect failed: ", m.Format(), err)
		return errorx.Wrap(err, "new minio client failed")
	} else {
		_client = &Client{config: m, client: client}
		log.Info("minio connect success: ", m.Format())
		return nil
	}
}

func (m *Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// NewClient 初始化minio客户端
func (m *Config) NewClient() (*minio.Client, error) {
	if client, err := minio.New(m.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessId, m.AccessSecret, m.SessionToken),
		Secure: m.Secure,
		Region: Region,
	}); err != nil {
		return nil, errorx.Wrap(err, "new minio client failed")
	} else {
		return client, nil
	}
}

func (m *Config) MinioPath(fileName string) string {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	minioPath := filepath.Join(time.Now().Format(timex.TimestampFmt), uuid.NewString()+fileSuffix)
	if m.PrefixPath != "" {
		minioPath = filepath.Join(m.PrefixPath, minioPath)
	}
	return minioPath
}
