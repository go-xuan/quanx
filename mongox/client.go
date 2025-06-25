package mongox

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	config *Config
	client *mongo.Client
}

func (c *Client) Instance() *mongo.Client {
	return c.client
}

func (c *Client) Config() *Config {
	return c.config
}
