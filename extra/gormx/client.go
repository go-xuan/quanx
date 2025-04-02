package gormx

import (
	"gorm.io/gorm"
)

type Client struct {
	config *Config
	db     *gorm.DB
}

func (c *Client) Instance() *gorm.DB {
	return c.db
}

func (c *Client) Config() *Config {
	return c.config
}
