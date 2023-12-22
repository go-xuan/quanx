package mongox

import (
	"context"
	"fmt"
	"github.com/go-xuan/quanx/runner/nacosx"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// 运行器名称
func (m *Mongo) Name() string {
	return "init mongo"
}

// nacos配置文件
func (*Mongo) NacosConfig() *nacosx.Config {
	return &nacosx.Config{
		DataId: "mongo.yaml",
		Listen: false,
	}
}

// 本地配置文件
func (*Mongo) LocalConfig() string {
	return "conf/mongo.yaml"
}

// 运行器运行
func (m *Mongo) Run() (err error) {
	var client *mongo.Client
	client, err = m.NewClient()
	if err != nil {
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
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		return
	}
	return
}
