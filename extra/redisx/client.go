package redisx

import (
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	config *Config
	client redis.UniversalClient
}

func (c *Client) Instance() redis.UniversalClient {
	return c.client
}

func (c *Client) Config() *Config {
	return c.config
}

func (c *Client) Copy(target string, database int) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database
	if client, err := config.NewRedisClient(); err != nil {
		log.Error("redis connect failed:", config.Format())
		return nil, err
	} else {
		log.Info("redis connect success: ", config.Format())
		return &Client{config: config, client: client}, nil
	}
}
