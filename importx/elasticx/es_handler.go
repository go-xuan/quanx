package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var handler *Handler

// elastic处理器
type Handler struct {
	Config *Elastic
	Url    string
	Client *elastic.Client
}

func This() *Handler {
	if handler == nil {
		panic("The elastic handler has not been initialized, please check the relevant config")
	}
	return handler
}

func (h *Handler) Create(ctx context.Context, index, id string, body interface{}) (err error) {
	var result *elastic.IndexResponse
	if result, err = h.Client.Index().Index(index).Id(id).BodyJson(body).Do(ctx); err != nil {
		return
	}
	log.Printf("create success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
	return
}

func (h *Handler) Update(ctx context.Context, index, id string, body interface{}) (err error) {
	var result *elastic.UpdateResponse
	if result, err = h.Client.Update().Index(index).Id(id).Doc(body).Do(ctx); err != nil {
		return
	}
	log.Printf("update success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
	return
}

func (h *Handler) Delete(ctx context.Context, index, id string) (err error) {
	var result *elastic.DeleteResponse
	if result, err = h.Client.Delete().Index(index).Id(id).Do(ctx); err != nil {
		return
	}
	log.Printf("delete success: %s\n", result.Result)
	return
}

func (h *Handler) Get(ctx context.Context, index, id string) (result *elastic.GetResult, err error) {
	return h.Client.Get().Index(index).Id(id).Do(ctx)
}

func (h *Handler) Search(ctx context.Context, index string, query elastic.Query) (result *elastic.SearchResult, err error) {
	return h.Client.Search().Index(index).Query(query).Do(ctx)
}

// 获取索引中全部文档ID，sortField字段必须支持排序
func (h *Handler) AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) (ids []string, err error) {
	var total, offset int64
	var sortValue float64
	for total >= offset {
		var server *elastic.SearchService
		server = h.Client.Search().Index(index).
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
