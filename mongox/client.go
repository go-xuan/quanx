package mongox

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewClient 创建客户端
func NewClient(config *Config) (*Client, error) {
	client, err := NewMongoClient(config)
	if err != nil {
		return nil, errorx.Wrap(err, "mongodb client failed")
	}
	return &Client{config: config, client: client}, nil
}

// NewMongoClient 创建mongo客户端
func NewMongoClient(config *Config) (*mongo.Client, error) {
	ctx := context.TODO()
	// 连接mongo
	client, err := mongo.Connect(ctx, config.ClientOptions())
	if err != nil {
		return nil, errorx.Wrap(err, "mongo connect failed")
	}
	// PING
	if err = client.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, errorx.Wrap(err, "mongo ping failed")
	}
	return client, nil
}

// Client MongoDB客户端的封装
type Client struct {
	config *Config
	client *mongo.Client
}

func (c *Client) GetClient() *mongo.Client {
	return c.client
}

func (c *Client) GetInstance() interface{} {
	return c.client
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) Copy(target string, database string) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database

	logger := log.WithFields(config.LogFields())
	client, err := NewClient(config)
	if err != nil {
		logger.WithError(err).Error("mongo client copy failed")
		return nil, errorx.Wrap(err, "mongo client copy failed")
	}
	logger.Info("mongo client copy success")
	return client, nil
}

func (c *Client) Close() error {
	logger := log.WithFields(c.config.LogFields())
	if err := c.GetClient().Disconnect(context.Background()); err != nil {
		logger.Error("mongo client disconnect failed")
		return errorx.Wrap(err, "mongo client disconnect failed")
	}
	logger.Info("close mongo success")
	return nil
}
