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
func (m *Mongo) ToString() string {
	return fmt.Sprintf("host=%s port=%d database=%s", m.Host, m.Port, m.Database)
}

// 运行器名称
func (m *Mongo) Name() string {
	return "连接Mongo"
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
func (m *Mongo) Run() error {
	var client = m.NewClient()
	err := client.Ping(context.Background(), nil)
	if err != nil {
		log.Error("MongoDB连接失败！", m.ToString())
		log.Error("error : ", err)
		return err
	}
	handler = &Handler{Config: m, Client: client}
	log.Error("MongoDB连接成功！", m.ToString())
	return nil
}

func (m *Mongo) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, m.Host)
}

// 配置信息格式化
func (m *Mongo) NewClient() *mongo.Client {
	// 设置连接选项
	clientOptions := options.Client().ApplyURI(m.Uri())
	// 建立连接
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Errorf("Mongo Connect Failed!")
	}
	return client
}
