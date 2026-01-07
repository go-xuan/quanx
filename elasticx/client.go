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

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) GetInstance() *elastic.Client {
	return c.client
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, index string) (bool, error) {
	logger := log.WithField("index", index)
	resp, err := c.GetInstance().CreateIndex(index).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("create index failed")
		return false, errorx.Wrap(err, "create index failed")
	}
	logger.Info("create index success")
	return resp.Acknowledged, nil
}

// AllIndices 查询所有索引
func (c *Client) AllIndices(ctx context.Context) ([]string, error) {
	resp, err := c.GetInstance().CatIndices().Do(ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "cat indices failed")
	}
	var indices []string
	for _, row := range resp {
		indices = append(indices, row.Index)
	}
	return indices, nil
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, index string) (bool, error) {
	logger := log.WithField("index", index)
	resp, err := c.GetInstance().DeleteIndex(index).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("delete index failed")
		return false, errorx.Wrap(err, "delete index failed")
	}
	logger.Info("delete index success")
	return resp.Acknowledged, nil
}

// DeleteIndicesInBatches 批量删除索引
func (c *Client) DeleteIndicesInBatches(ctx context.Context, indices []string, limit int) error {
	if err := taskx.NewSplitterStrategy(limit).Execute(ctx, len(indices), func(ctx context.Context, start, end, batch int) error {
		indices_ := indices[start:end]
		logger := log.WithField("start", start).WithField("end", end).WithField("batch", batch)
		if _, err := c.DeleteIndices(ctx, indices_); err != nil {
			logger.WithError(err).Error("delete indices in batches failed")
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
	resp, err := c.GetInstance().DeleteIndex(indices...).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("delete indices failed")
		return false, errorx.Wrap(err, "delete indices failed")
	}
	logger.Info("delete indices success")
	return resp.Acknowledged, nil
}

func (c *Client) Create(ctx context.Context, index, id string, body any) error {
	logger := log.WithField("index", index).WithField("id", id)
	resp, err := c.GetInstance().Index().Index(index).Id(id).BodyJson(body).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("create failed")
		return errorx.Wrap(err, "create failed")
	}
	logger.WithField("type", resp.Type).Info("create success")
	return nil
}

func (c *Client) Update(ctx context.Context, index, id string, body any) error {
	logger := log.WithField("index", index).WithField("id", id)
	resp, err := c.GetInstance().Update().Index(index).Id(id).Doc(body).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("update failed")
		return errorx.Wrap(err, "update failed")
	}
	logger.WithField("type", resp.Type).Info("update success")
	return nil
}

func (c *Client) Delete(ctx context.Context, index, id string) error {
	logger := log.WithField("index", index).WithField("id", id)
	resp, err := c.GetInstance().Delete().Index(index).Id(id).Do(ctx)
	if err != nil {
		logger.WithError(err).Error("delete failed")
		return errorx.Wrap(err, "delete index failed")
	}
	logger.WithField("type", resp.Type).Info("delete success")
	return nil
}

func (c *Client) Get(ctx context.Context, index, id string) (*elastic.GetResult, error) {
	return c.GetInstance().Get().Index(index).Id(id).Do(ctx)
}

func (c *Client) Search(ctx context.Context, index string, query elastic.Query) (*elastic.SearchResult, error) {
	return c.GetInstance().Search().Index(index).Query(query).Do(ctx)
}

// AllDocId 获取索引中全部文档ID，sortField字段必须支持排序
func (c *Client) AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) ([]string, error) {
	var total, offset int64
	var sortValue float64
	var ids []string
	for offset <= total {
		var server *elastic.SearchService
		server = c.GetInstance().Search().Index(index).
			Query(query).
			TrackTotalHits(true).
			Sort(sortField, true).
			Size(10000)
		if sortValue != 0 {
			server = server.SearchAfter(sortValue)
		}
		result, err := server.Do(ctx)
		if result == nil || err != nil {
			total, offset = 0, 10000
			return nil, err
		}
		for _, hit := range result.Hits.Hits {
			ids = append(ids, hit.Id)
			sortValue = hit.Sort[0].(float64)
		}
		total = result.Hits.TotalHits.Value
		offset += 10000
	}
	return ids, nil
}
