package miniox

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/timex"
)

func NewConfigurator(conf *Config) configx.Configurator {
	return conf
}

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
	return fmtx.Yellow.XSPrintf("host=%s port=%v accessId=%s bucketName=%s", m.Host, m.Port, m.AccessId, m.BucketName)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "minio.yaml",
		NacosDataId: "minio.yaml",
		Listen:      false,
	}
}

func (m *Config) Execute() error {
	if client, err := m.NewClient(); err != nil {
		log.Error("minio connect failed: ", m.Format(), err)
		return errorx.Wrap(err, "new minio client failed")
	} else {
		_handler = &Handler{config: m, client: client}
		log.Info("minio connect successfully: ", m.Format())
		return nil
	}
}

func (m *Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// NewClient 初始化minio客户端
func (m *Config) NewClient() (client *minio.Client, err error) {
	if client, err = minio.New(m.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessId, m.AccessSecret, m.SessionToken),
		Secure: m.Secure,
		Region: Region,
	}); err != nil {
		return
	}
	return
}

func (m *Config) MinioPath(fileName string) (minioPath string) {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	minioPath = filepath.Join(time.Now().Format(timex.TimestampFmt), uuid.NewString()+fileSuffix)
	if m.PrefixPath != "" {
		minioPath = filepath.Join(m.PrefixPath, minioPath)
	}
	return
}
