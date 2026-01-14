package elasticx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/httpx"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

// NewClient 创建客户端
func NewClient(config *Config) (*Client, error) {
	client, err := NewEsClient(config)
	if err != nil {
		return nil, errorx.Wrap(err, "create elastic-search client failed")
	}
	return &Client{config: config, client: client}, nil
}

// NewEsClient 创建es客户端
func NewEsClient(config *Config) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Url),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(config.Username, config.Password),
		elastic.SetHttpClient(httpx.NewClient()),
	)
	if err != nil {
		return nil, errorx.Wrap(err, "create elastic client failed")
	}
	res, code, err := client.Ping(config.Url).Do(context.Background())
	if err != nil || code != 200 {
		return nil, errorx.Wrap(err, "ping elastic-search failed")
	}
	log.Info("elastic-search version: ", res.Version.Number)
	return client, nil
}

type Client struct {
	config *Config
	client *elastic.Client
}

func (c *Client) GetClient() *elastic.Client {
	return c.client
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) GetInstance() interface{} {
	return c.client
}

func (c *Client) Close() error {
	logger := log.WithFields(c.config.LogFields())
	c.client.Stop()
	logger.Info("close elastic client success")
	return nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, index string) (bool, error) {
	res, err := c.client.CreateIndex(index).Do(ctx)
	if err != nil {
		return false, errorx.Wrap(err, "create index failed")
	}
	return res.Acknowledged, nil
}
