package mongox

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	config *Config
	client *mongo.Client
}

func (c *Client) GetInstance() *mongo.Client {
	return c.client
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) Copy(target string, database string) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database
	if client, err := config.NewClient(); err != nil {
		config.LogEntry().WithError(err).Error("new mongo client failed")
		return nil, err
	} else {
		config.LogEntry().Info("new mongo client success")
		return &Client{
			config: config,
			client: client,
		}, nil
	}
}
