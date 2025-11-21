package redisx

import (
	"github.com/redis/go-redis/v9"
)

type Client struct {
	config *Config
	client redis.UniversalClient
}

func (c *Client) GetInstance() redis.UniversalClient {
	return c.client
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) Copy(target string, database int) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database
	if client, err := config.NewRedisClient(); err != nil {
		config.LogEntry().WithError(err).Error("new redis client failed")
		return nil, err
	} else {
		config.LogEntry().Info("new redis client success")
		return &Client{
			config: config,
			client: client,
		}, nil
	}
}
