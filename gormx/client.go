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

func (c *Client) Copy(target, database string) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database
	if db, err := config.NewGormDB(); err != nil {
		config.LogEntry().Error("database connect failed")
		return nil, err
	} else {
		config.LogEntry().Info("database connect success")
		return &Client{config: config, db: db}, nil
	}
}
