package miniox

import (
	"fmt"
	"github.com/quanxiaoxuan/quanx/common/constx"
	"path"
	"path/filepath"
	"time"

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
func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d userName=%s bucketName=%s", conf.Host, conf.Port, conf.UserName, conf.BucketName)
}

// 配置信息格式化
func (conf *Config) Endpoint() string {
	return fmt.Sprintf("%s:%d", conf.Host, conf.Port)
}

// 初始化minio客户端
func (conf *Config) NewClient() (client *minio.Client, err error) {
	client, err = minio.New(conf.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(conf.UserName, conf.UserName, conf.Token),
		Secure: conf.Secure,
	})
	if err != nil {
		log.Warn("初始化Mino客户端失败 : ", err)
		return
	}
	return
}

func (conf *Config) MinioPath(fileName string) (minioPath string) {
	fileSuffix := filepath.Ext(filepath.Base(fileName))
	fileName = uuid.NewString() + fileSuffix
	minioPath = path.Join(time.Now().Format(constx.TimestampFmt), fileName)
	if conf.RootDir != "" {
		minioPath = path.Join(conf.RootDir, minioPath)
	}
	return
}
