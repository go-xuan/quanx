package mongox

import (
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	config *Config
	client *mongo.Client
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) GetInstance() *mongo.Client {
	return c.client
}

func (c *Client) Copy(target string, database string) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database

	logger := log.WithFields(config.LogFields())
	client, err := config.NewClient()
	if err != nil {
		logger.WithError(err).Error("mongo client copy failed")
		return nil, errorx.Wrap(err, "mongo client copy failed")
	}
	logger.Info("mongo client copy success")
	return &Client{
		config: config,
		client: client,
	}, nil
}
