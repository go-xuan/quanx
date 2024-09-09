package miniox

import (
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/confx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/timex"
)

type Minio struct {
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

// Info 配置信息格式化
func (m *Minio) Info() string {
	return fmt.Sprintf("host=%s port=%d accessId=%s bucketName=%s", m.Host, m.Port, m.AccessId, m.BucketName)
}

// Title 配置器标题
func (*Minio) Title() string {
	return "Minio"
}

// Reader 配置文件读取
func (*Minio) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "minio.yaml",
		NacosDataId: "minio.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (m *Minio) Run() error {
	if client, err := m.NewClient(); err != nil {
		log.Error("minio connect failed: ", m.Info(), err)
		return errorx.Wrap(err, "NewClient failed")
	} else {
		handler = &Handler{config: m, client: client}
		log.Info("minio connect successful: ", m.Info())
		return nil
	}
}

// Endpoint 配置信息格式化
func (m *Minio) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// NewClient 初始化minio客户端
func (m *Minio) NewClient() (client *minio.Client, err error) {
	if client, err = minio.New(m.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessId, m.AccessSecret, m.SessionToken),
		Secure: m.Secure,
		Region: Region,
	}); err != nil {
		return
	}
	return
}

func (m *Minio) MinioPath(fileName string) (minioPath string) {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	minioPath = path.Join(time.Now().Format(timex.TimestampFmt), uuid.NewString()+fileSuffix)
	if m.PrefixPath != "" {
		minioPath = path.Join(m.PrefixPath, minioPath)
	}
	return
}
