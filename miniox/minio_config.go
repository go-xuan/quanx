package miniox

import (
	"fmt"
	"path"
	"path/filepath"
	"time"

	"github.com/go-xuan/quanx/utilx/timex"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Host       string `yaml:"host" json:"host"`             // 主机
	Port       int    `yaml:"port" json:"port"`             // 端口
	UserName   string `yaml:"userName" json:"userName"`     // 用户名
	Password   string `yaml:"password" json:"password"`     // 密码
	Token      string `yaml:"token" json:"token"`           // token
	Secure     bool   `yaml:"secure" json:"secure"`         // 是否使用https
	BucketName string `yaml:"bucketName" json:"bucketName"` // 桶名
	RootDir    string `yaml:"rootDir" json:"rootDir"`       // 对象存储根路径
	ExpireHour int    `yaml:"expireHour" json:"expireHour"` // 下载链接有效时长
}

// 配置信息格式化
func (config *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d userName=%s bucketName=%s", config.Host, config.Port, config.UserName, config.BucketName)
}

func (config *Config) Init() {
	InitMinioX(config)
}

// 初始化minioX
func InitMinioX(conf *Config) {
	client, err := conf.NewClient()
	if err == nil {
		instance = &Handler{Config: conf, Client: client}
		log.Info("Minio连接成功!", conf.Format())
	} else {
		log.Error("Minio连接失败!", conf.Format())
		log.Error("error : ", err)
	}
}

// 配置信息格式化
func (config *Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// 初始化minio客户端
func (config *Config) NewClient() (client *minio.Client, err error) {
	client, err = minio.New(config.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(config.UserName, config.UserName, config.Token),
		Secure: config.Secure,
	})
	if err != nil {
		log.Warn("初始化Mino客户端失败 : ", err)
		return
	}
	return
}

func (config *Config) MinioPath(fileName string) (minioPath string) {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	fileName = uuid.NewString() + fileSuffix
	minioPath = path.Join(time.Now().Format(timex.TimestampFmt), fileName)
	if config.RootDir != "" {
		minioPath = path.Join(config.RootDir, minioPath)
	}
	return
}
