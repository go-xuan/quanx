package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/slicex"
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

// 查询所有索引
func (h *Handler) AllIndices(ctx context.Context) (indices []string, err error) {
	var resp elastic.CatIndicesResponse
	if resp, err = h.Client.CatIndices().Do(ctx); err != nil {
		return
	}
	for _, row := range resp {
		indices = append(indices, row.Index)
	}
	return
}

// 创建索引
func (h *Handler) CreateIndex(ctx context.Context, index string) (ok bool, err error) {
	var resp *elastic.IndicesDeleteResponse
	if resp, err = h.Client.DeleteIndex(index).Do(ctx); err != nil {
		return
	}
	ok = resp.Acknowledged
	return
}

// 删除索引
func (h *Handler) DeleteIndex(ctx context.Context, index string) (ok bool, err error) {
	var resp *elastic.IndicesDeleteResponse
	if resp, err = h.Client.DeleteIndex(index).Do(ctx); err != nil {
		return
	}
	ok = resp.Acknowledged
	return
}

// 批量索引
func (h *Handler) DeleteIndices(ctx context.Context, indices []string) (ok bool, err error) {
	if err = slicex.ExecInBatches(len(indices), 100, func(x int, y int) (err error) {
		var resp *elastic.IndicesDeleteResponse
		if resp, err = h.Client.DeleteIndex(indices[x:y]...).Do(ctx); err != nil {
			panic(err)
		}
		log.Printf("delete indices[%d-%d] acknowledged：%v\n", x, y, resp.Acknowledged)
		return
	}); err != nil {
		return
	}
	return
}

func (h *Handler) Create(ctx context.Context, index, id string, body any) (err error) {
	var resp *elastic.IndexResponse
	if resp, err = h.Client.Index().Index(index).Id(id).BodyJson(body).Do(ctx); err != nil {
		return
	}
	log.Printf("create success: id=%s index=%s type=%s\n", resp.Id, resp.Index, resp.Type)
	return
}

func (h *Handler) Update(ctx context.Context, index, id string, body any) (err error) {
	var resp *elastic.UpdateResponse
	if resp, err = h.Client.Update().Index(index).Id(id).Doc(body).Do(ctx); err != nil {
		return
	}
	log.Printf("update success: id=%s index=%s type=%s\n", resp.Id, resp.Index, resp.Type)
	return
}

func (h *Handler) Delete(ctx context.Context, index, id string) (err error) {
	var resp *elastic.DeleteResponse
	if resp, err = h.Client.Delete().Index(index).Id(id).Do(ctx); err != nil {
		return
	}
	log.Printf("delete success: %s\n", resp.Result)
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
