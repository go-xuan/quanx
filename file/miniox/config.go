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

	"github.com/go-xuan/quanx/core/confx"
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

// 配置信息格式化
func (m *Minio) ToString() string {
	return fmt.Sprintf("host=%s port=%d accessId=%s bucketName=%s", m.Host, m.Port, m.AccessId, m.BucketName)
}

// 配置器名称
func (*Minio) Theme() string {
	return "Minio"
}

// 配置文件读取
func (*Minio) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "minio.yaml",
		NacosDataId: "minio.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (m *Minio) Run() (err error) {
	var client *minio.Client
	var toString = m.ToString()
	if client, err = m.NewClient(); err != nil {
		log.Error("Minio Connect Failed : ", toString, err)
		return
	}
	handler = &Handler{Config: m, Client: client}
	log.Info("Minio Connect Successful : ", toString)
	return
}

// 配置信息格式化
func (m *Minio) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// 初始化minio客户端
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
