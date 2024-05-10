package mongox

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-xuan/quanx/core/confx"
)

type Mongo struct {
	Host     string `yaml:"host" json:"host"`         // 主机
	Port     int    `yaml:"port" json:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名
}

// 配置信息格式化
func (m *Mongo) ToString() string {
	return fmt.Sprintf("host=%s port=%d database=%s", m.Host, m.Port, m.Database)
}

// 配置器名称
func (m *Mongo) Theme() string {
	return "Mongo"
}

// 配置文件读取
func (*Mongo) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (m *Mongo) Run() (err error) {
	var client *mongo.Client
	var toString = m.ToString()
	if client, err = m.NewClient(); err != nil {
		log.Error("Mongo Connect Failed: ", toString, err)
		return
	}
	handler = &Handler{Config: m, Client: client}
	log.Info("Mongo Connect Successful: ", toString)
	return

}

func (m *Mongo) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, m.Host)
}

// 配置信息格式化
func (m *Mongo) NewClient() (client *mongo.Client, err error) {
	// 设置连接选项
	clientOptions := options.Client().ApplyURI(m.Uri())
	// 建立连接
	if client, err = mongo.Connect(context.Background(), clientOptions); err != nil {
		return
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return
	}
	return
}
