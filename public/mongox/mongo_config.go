package mongox

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string `yaml:"host" json:"host"`         // 主机
	Port     int    `yaml:"port" json:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名
}

// 配置信息格式化
func (conf *Config) Format() string {
	return fmt.Sprintf("host=%s port=%d database=%s", conf.Host, conf.Port, conf.Database)
}

func (conf *Config) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", conf.Username, conf.Password, conf.Host)
}

// 配置信息格式化
func (conf *Config) NewClient() *mongo.Client {
	// 设置连接选项
	clientOptions := options.Client().ApplyURI(conf.Uri())
	// 建立连接
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Errorf("Mongo Connect Failed!")
	}
	return client
}
