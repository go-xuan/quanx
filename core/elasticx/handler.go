package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/taskx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the elastic handler has not been initialized, please check the relevant config")
	}
	return _handler
}

type Handler struct {
	multi   bool
	config  *Config
	client  *elastic.Client
	configs map[string]*Config
	clients map[string]*elastic.Client
}

func (h *Handler) GetConfig(source ...string) *Config {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) GetClient(source ...string) *elastic.Client {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clients[source[0]]; ok {
			return client
		}
	}
	return h.client
}

// IsInitialized 是否初始化
func IsInitialized() bool {
	return _handler != nil
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

// GetClient 获取客户端
func GetClient(source ...string) *elastic.Client {
	return this().GetClient(source...)
}

// CreateIndex 创建索引
func CreateIndex(ctx context.Context, index string) (bool, error) {
	if resp, err := GetClient().CreateIndex(index).Do(ctx); err != nil {
		log.WithField("index", index).Error("create index failed: ", err)
		return false, errorx.Wrap(err, "create index failed")
	} else {
		log.WithField("index", index).Error("create index success")
		return resp.Acknowledged, nil
	}
}

// AllIndices 查询所有索引
func AllIndices(ctx context.Context) ([]string, error) {
	if resp, err := GetClient().CatIndices().Do(ctx); err != nil {
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
func DeleteIndex(ctx context.Context, index string) (bool, error) {
	if resp, err := GetClient().DeleteIndex(index).Do(ctx); err != nil {
		log.WithField("index", index).Error("delete index failed: ", err)
		return false, errorx.Wrap(err, "delete index failed")
	} else {
		log.WithField("index", index).Error("delete index success")
		return resp.Acknowledged, nil
	}
}

// DeleteIndices 批量索引
func DeleteIndices(ctx context.Context, indices []string) (bool, error) {
	var ok bool
	if err := taskx.ExecWithBatches(len(indices), 100, func(start int, end int) error {
		deleteIndices := indices[start:end]
		if resp, err := GetClient().DeleteIndex(deleteIndices...).Do(ctx); err != nil {
			log.WithField("deleteIndices", deleteIndices).Error("delete indices failed: ", err)
			return err
		} else {
			log.WithField("deleteIndices", deleteIndices).Info("delete indices success")
			ok = resp.Acknowledged
			return nil
		}
	}); err != nil {
		return ok, err
	}
	return ok, nil
}

func Create(ctx context.Context, index, id string, body any) error {
	if resp, err := GetClient().Index().Index(index).Id(id).BodyJson(body).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("create failed: ", err)
		return err
	} else {
		log.WithField("index", resp.Index).WithField("id", resp.Id).
			WithField("type", resp.Type).Info("create success")
		return nil
	}
}

func Update(ctx context.Context, index, id string, body any) error {
	if resp, err := GetClient().Update().Index(index).Id(id).Doc(body).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("update failed: ", err)
		return errorx.Wrap(err, "update index failed")
	} else {
		log.WithField("index", resp.Index).WithField("id", resp.Id).
			WithField("type", resp.Type).Info("update success")
		return nil
	}
}

func Delete(ctx context.Context, index, id string) error {
	if resp, err := GetClient().Delete().Index(index).Id(id).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("delete failed: ", err)
		return errorx.Wrap(err, "delete index failed")
	} else {
		log.WithField("index", resp.Index).WithField("id", resp.Id).
			WithField("type", resp.Type).Info("delete success")
		return nil
	}
}

func Get(ctx context.Context, index, id string) (*elastic.GetResult, error) {
	return GetClient().Get().Index(index).Id(id).Do(ctx)
}

func Search(ctx context.Context, index string, query elastic.Query) (*elastic.SearchResult, error) {
	return GetClient().Search().Index(index).Query(query).Do(ctx)
}

// AllDocId 获取索引中全部文档ID，sortField字段必须支持排序
func AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) ([]string, error) {
	var total, offset int64
	var sortValue float64
	var ids []string
	for offset <= total {
		var server *elastic.SearchService
		server = GetClient().Search().Index(index).
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
