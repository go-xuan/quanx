package gormx

import (
	log "github.com/sirupsen/logrus"
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
		log.Error("database connect failed: ", config.Info())
		return nil, err
	} else {
		log.Info("database connect success: ", config.Info())
		return &Client{config: config, db: db}, nil
	}
}
