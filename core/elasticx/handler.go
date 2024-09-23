package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/execx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the elastic handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func GetConfig() *Config {
	return this().GetConfig()
}

func Client() *elastic.Client {
	return this().GetClient()
}

// CreateIndex 创建索引
func CreateIndex(ctx context.Context, index string) (ok bool, err error) {
	var resp *elastic.IndicesCreateResult
	if resp, err = Client().CreateIndex(index).Do(ctx); err != nil {
		log.WithField("index", index).Error("create index failed: ", err)
		return
	}
	log.WithField("index", index).Error("create index success")
	ok = resp.Acknowledged
	return
}

type Handler struct {
	config *Config
	client *elastic.Client
}

func (h *Handler) GetConfig() *Config {
	return h.config
}

func (h *Handler) GetClient() *elastic.Client {
	return h.client
}

// AllIndices 查询所有索引
func (h *Handler) AllIndices(ctx context.Context) (indices []string, err error) {
	var resp elastic.CatIndicesResponse
	if resp, err = h.client.CatIndices().Do(ctx); err != nil {
		return
	}
	for _, row := range resp {
		indices = append(indices, row.Index)
	}
	return
}

// DeleteIndex 删除索引
func (h *Handler) DeleteIndex(ctx context.Context, index string) (ok bool, err error) {
	var resp *elastic.IndicesDeleteResponse
	if resp, err = h.client.DeleteIndex(index).Do(ctx); err != nil {
		log.WithField("index", index).Error("delete index failed: ", err)
		return
	}
	log.WithField("index", index).Error("delete index success")
	ok = resp.Acknowledged
	return
}

// DeleteIndices 批量索引
func (h *Handler) DeleteIndices(ctx context.Context, indices []string) (ok bool, err error) {
	if err = execx.InBatches(len(indices), 100, func(x int, y int) (err error) {
		var resp *elastic.IndicesDeleteResponse
		deleteIndices := indices[x:y]
		if resp, err = h.client.DeleteIndex(deleteIndices...).Do(ctx); err != nil {
			log.WithField("deleteIndices", deleteIndices).Error("delete indices failed: ", err)
			return
		}
		log.WithField("deleteIndices", deleteIndices).Error("delete indices success")
		ok = resp.Acknowledged
		return
	}); err != nil {
		return
	}
	return
}

func (h *Handler) Create(ctx context.Context, index, id string, body any) (err error) {
	var resp *elastic.IndexResponse
	if resp, err = h.client.Index().Index(index).Id(id).BodyJson(body).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("create failed: ", err)
		return
	}
	log.WithField("index", resp.Index).WithField("id", resp.Id).
		WithField("type", resp.Type).Info("create success")
	return
}

func (h *Handler) Update(ctx context.Context, index, id string, body any) (err error) {
	var resp *elastic.UpdateResponse
	if resp, err = h.client.Update().Index(index).Id(id).Doc(body).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("update failed: ", err)
		return
	}
	log.WithField("index", resp.Index).WithField("id", resp.Id).
		WithField("type", resp.Type).Info("update success")
	return
}

func (h *Handler) Delete(ctx context.Context, index, id string) (err error) {
	var resp *elastic.DeleteResponse
	if resp, err = h.client.Delete().Index(index).Id(id).Do(ctx); err != nil {
		log.WithField("index", index).WithField("id", id).
			Error("delete failed: ", err)
		return
	}
	log.WithField("index", resp.Index).WithField("id", resp.Id).
		WithField("type", resp.Type).Info("delete success")
	return
}

func (h *Handler) Get(ctx context.Context, index, id string) (result *elastic.GetResult, err error) {
	return h.client.Get().Index(index).Id(id).Do(ctx)
}

func (h *Handler) Search(ctx context.Context, index string, query elastic.Query) (result *elastic.SearchResult, err error) {
	return h.client.Search().Index(index).Query(query).Do(ctx)
}

// AllDocId 获取索引中全部文档ID，sortField字段必须支持排序
func (h *Handler) AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) (ids []string, err error) {
	var total, offset int64
	var sortValue float64
	for total >= offset {
		var server *elastic.SearchService
		server = h.client.Search().Index(index).
			Query(query).
			TrackTotalHits(true).
			Sort(sortField, true).
			Size(10000)
		if sortValue != 0 {
			server = server.SearchAfter(sortValue)
		}
		var result *elastic.SearchResult
		if result, err = server.Do(ctx); result == nil || err != nil {
			total, offset = 0, 10000
			return
		} else {
			for _, hit := range result.Hits.Hits {
				ids = append(ids, hit.Id)
				sortValue = hit.Sort[0].(float64)
			}
			total = result.Hits.TotalHits.Value
			offset += 10000
		}
	}
	return
}
