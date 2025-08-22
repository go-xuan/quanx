package elasticx

import (
	"context"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/taskx"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	config *Config
	client *elastic.Client
}

func (c *Client) Instance() *elastic.Client {
	return c.client
}

func (c *Client) Config() *Config {
	return c.config
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, index string) (bool, error) {
	logger := log.WithField("index", index)
	if resp, err := c.Instance().CreateIndex(index).Do(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("create index failed")
		return false, errorx.Wrap(err, "create index failed")
	} else {
		logger.Info("create index success")
		return resp.Acknowledged, nil
	}
}

// AllIndices 查询所有索引
func (c *Client) AllIndices(ctx context.Context) ([]string, error) {
	if resp, err := c.Instance().CatIndices().Do(ctx); err != nil {
		return nil, errorx.Wrap(err, "cat indices failed")
	} else {
		var indices []string
		for _, row := range resp {
			indices = append(indices, row.Index)
		}
		return indices, nil
	}
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, index string) (bool, error) {
	logger := log.WithField("index", index)
	if resp, err := c.Instance().DeleteIndex(index).Do(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("delete index failed")
		return false, errorx.Wrap(err, "delete index failed")
	} else {
		logger.Info("delete index success")
		return resp.Acknowledged, nil
	}
}

// DeleteIndicesInBatches 批量删除索引
func (c *Client) DeleteIndicesInBatches(ctx context.Context, indices []string, limit int) error {
	if err := taskx.NewSplitter(limit).Execute(ctx, len(indices), func(ctx context.Context, start int, end, batch int) error {
		indices_ := indices[start:end]
		logger := log.WithField("start", start).WithField("end", end).WithField("batch", batch)
		if _, err := c.DeleteIndices(ctx, indices_); err != nil {
			logger.WithField("error", err.Error()).Error("delete indices in batches failed")
			return errorx.Wrap(err, "delete indices in batches failed")
		}
		logger.Info("delete indices in batches success")
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// DeleteIndices 删除索引
func (c *Client) DeleteIndices(ctx context.Context, indices []string) (bool, error) {
	logger := log.WithField("indices", indices)
	resp, err := c.Instance().DeleteIndex(indices...).Do(ctx)
	if err != nil {
		logger.WithField("error", err.Error()).Error("delete indices failed")
		return false, errorx.Wrap(err, "delete indices failed")
	}
	logger.Info("delete indices success")
	return resp.Acknowledged, nil
}

func (c *Client) Create(ctx context.Context, index, id string, body any) error {
	logger := log.WithField("index", index).WithField("id", id)
	if resp, err := c.Instance().Index().Index(index).Id(id).BodyJson(body).Do(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("create failed")
		return errorx.Wrap(err, "create failed")
	} else {
		logger.WithField("type", resp.Type).Info("create success")
		return nil
	}
}

func (c *Client) Update(ctx context.Context, index, id string, body any) error {
	logger := log.WithField("index", index).WithField("id", id)
	if resp, err := c.Instance().Update().Index(index).Id(id).Doc(body).Do(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("update failed")
		return errorx.Wrap(err, "update failed")
	} else {
		logger.WithField("type", resp.Type).Info("update success")
		return nil
	}
}

func (c *Client) Delete(ctx context.Context, index, id string) error {
	logger := log.WithField("index", index).WithField("id", id)
	if resp, err := c.Instance().Delete().Index(index).Id(id).Do(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("delete failed: ", err)
		return errorx.Wrap(err, "delete index failed")
	} else {
		logger.WithField("type", resp.Type).Info("delete success")
		return nil
	}
}

func (c *Client) Get(ctx context.Context, index, id string) (*elastic.GetResult, error) {
	return c.Instance().Get().Index(index).Id(id).Do(ctx)
}

func (c *Client) Search(ctx context.Context, index string, query elastic.Query) (*elastic.SearchResult, error) {
	return c.Instance().Search().Index(index).Query(query).Do(ctx)
}

// AllDocId 获取索引中全部文档ID，sortField字段必须支持排序
func (c *Client) AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) ([]string, error) {
	var total, offset int64
	var sortValue float64
	var ids []string
	for offset <= total {
		var server *elastic.SearchService
		server = c.Instance().Search().Index(index).
			Query(query).
			TrackTotalHits(true).
			Sort(sortField, true).
			Size(10000)
		if sortValue != 0 {
			server = server.SearchAfter(sortValue)
		}
		if result, err := server.Do(ctx); result == nil || err != nil {
			total, offset = 0, 10000
			return nil, err
		} else {
			for _, hit := range result.Hits.Hits {
				ids = append(ids, hit.Id)
				sortValue = hit.Sort[0].(float64)
			}
			total = result.Hits.TotalHits.Value
			offset += 10000
		}
	}
	return ids, nil
}
