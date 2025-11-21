package gormx

import (
	"gorm.io/gorm"
)

type Client struct {
	config *Config
	db     *gorm.DB
}

func (c *Client) GetInstance() *gorm.DB {
	return c.db
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) Copy(target, database string) (*Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database
	if client, err := config.NewClient(); err != nil {
		config.LogEntry().WithError(err).Error("new gorm client failed")
		return nil, err
	} else {
		config.LogEntry().Info("new gorm client success")
		return client, nil
	}
}
