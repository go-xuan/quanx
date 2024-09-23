package mongox

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/utils/fmtx"
)

type Config struct {
	Host     string `yaml:"host" json:"host"`         // 主机
	Port     int    `yaml:"port" json:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名
}

func (m *Config) Format() string {
	return fmtx.Yellow.XSPrintf("host=%s port=%v database=%s", m.Host, m.Port, m.Database)
}

func (m *Config) ID() string {
	return "mongo"
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

func (m *Config) Execute() error {
	if client, err := m.NewClient(); err != nil {
		log.Error("mongo connect failed: ", m.Format(), err)
		return errorx.Wrap(err, "new mongo client error")
	} else {
		_handler = &Handler{config: m, client: client}
		log.Info("mongo connect successfully: ", m.Format())
		return nil
	}
}

func (m *Config) Uri() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", m.Username, m.Password, m.Host)
}

func (m *Config) NewClient() (client *mongo.Client, err error) {
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
