package redisx

import "github.com/redis/go-redis/v9"

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
