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

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/timex"
	"github.com/go-xuan/quanx/utils/fmtx"
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

func (*Minio) ID() string {
	return "minio"
}

func (m *Minio) Format() string {
	return fmtx.Yellow.XSPrintf("host=%s port=%v accessId=%s bucketName=%s", m.Host, m.Port, m.AccessId, m.BucketName)
}

func (*Minio) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "minio.yaml",
		NacosDataId: "minio.yaml",
		Listen:      false,
	}
}

func (m *Minio) Execute() error {
	if client, err := m.NewClient(); err != nil {
		log.Error("minio connect failed: ", m.Format(), err)
		return errorx.Wrap(err, "new minio client failed")
	} else {
		handler = &Handler{config: m, client: client}
		log.Info("minio connect successfully: ", m.Format())
		return nil
	}
}

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
