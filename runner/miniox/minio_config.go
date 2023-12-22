package miniox

import (
	"fmt"
	"github.com/go-xuan/quanx/runner/nacosx"
	"path"
	"path/filepath"
	"time"

	"github.com/go-xuan/quanx/utilx/timex"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
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
func (m *Minio) ToString(title string) string {
	return fmt.Sprintf("%s => host=%s port=%d accessId=%s bucketName=%s", title, m.Host, m.Port, m.AccessId, m.BucketName)
}

// 运行器名称
func (m *Minio) Name() string {
	return "init minio"
}

// nacos配置文件
func (*Minio) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "minio.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (*Minio) LocalConfig() string {
	return "conf/minio.yaml"
}

// 运行器运行
func (m *Minio) Run() error {
	client, err := m.NewClient()
	if err != nil {
		log.Error(m.ToString("minio connect failed!"))
		log.Error("error : ", err)
		return err
	}
	handler = &Handler{Config: m, Client: client}
	log.Info(m.ToString("minio connect successful!"))
	return nil
}

// 配置信息格式化
func (m *Minio) Endpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

// 初始化minio客户端
func (m *Minio) NewClient() (client *minio.Client, err error) {
	client, err = minio.New(m.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessId, m.AccessSecret, m.SessionToken),
		Secure: m.Secure,
		Region: "us-east-1",
	})
	if err != nil {
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
