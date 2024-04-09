package mongox

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-xuan/quanx/configx"
)

type Mongo struct {
	Host     string `yaml:"host" json:"host"`         // 主机
	Port     int    `yaml:"port" json:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名
}

// 配置信息格式化
func (m *Mongo) ToString(title string) string {
	return fmt.Sprintf("%s => host=%s port=%d database=%s", title, m.Host, m.Port, m.Database)
}

// 配置器名称
func (m *Mongo) Title() string {
	return "init mongo"
}

// 配置文件读取
func (*Mongo) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (m *Mongo) Run() (err error) {
	var client *mongo.Client
	if client, err = m.NewClient(); err != nil {
		log.Error(m.ToString("mongo connect failed!"))
		log.Error("error : ", err)
		return
	}
	handler = &Handler{Config: m, Client: client}
	log.Error(m.ToString("mongo connect successful!"))
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
