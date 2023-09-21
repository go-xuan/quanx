package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var CTL *Controller

// ES控制器
type Controller struct {
	Config *Config
	Url    string
	Client *elastic.Client
}

func Init(conf *Config) {
	var url = conf.Url()
	client, err := conf.NewClient(url)
	if err != nil {
		log.Error("redis连接失败！", conf.Format())
		log.Error("error : ", err)
		return
	}
	var ok bool
	if ok, err = Ping(client, url); ok && err == nil {
		CTL = &Controller{Config: conf, Url: url, Client: client}
		log.Error("ElasticSearch连接成功！", conf.Format())
	} else {
		log.Error("ElasticSearch连接失败！", conf.Format())
		log.Error("error : ", err)
	}
}

func Ping(client *elastic.Client, url string) (bool, error) {
	info, code, err := client.Ping(url).Do(context.Background())
	if err != nil && code != 200 {
		return false, err
	}
	log.Info("ElasticSearch 版本 : ", info.Version.Number)
	return true, nil
}

func (ctl *Controller) Create(ctx context.Context, index, id string, body interface{}) {
	var result *elastic.IndexResponse
	var err error
	result, err = ctl.Client.Index().Index(index).Id(id).BodyJson(body).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Printf("create success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
}

func (ctl *Controller) Update(ctx context.Context, index, id string, body interface{}) {
	var result *elastic.UpdateResponse
	var err error
	result, err = ctl.Client.Update().Index(index).Id(id).Doc(body).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Printf("update success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
}

func (ctl *Controller) Delete(ctx context.Context, index, id string) {
	var result *elastic.DeleteResponse
	var err error
	result, err = ctl.Client.Delete().Index(index).Id(id).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	log.Printf("delete success: %s\n", result.Result)
}

func (ctl *Controller) Get(ctx context.Context, index, id string) *elastic.GetResult {
	var result *elastic.GetResult
	var err error
	result, err = ctl.Client.Get().Index(index).Id(id).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	return result
}

func (ctl *Controller) Search(ctx context.Context, index string, query elastic.Query) *elastic.SearchResult {
	var result *elastic.SearchResult
	var err error
	result, err = ctl.Client.Search().Index(index).Query(query).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	return result
}

// 获取索引中全部文档ID，sortField字段必须支持排序
func (ctl *Controller) AllDocId(ctx context.Context, index string, query elastic.Query, sortField string) (ids []string, err error) {
	var total, offset int64
	var sortValue float64
	for total >= offset {
		var server *elastic.SearchService
		server = ctl.Client.Search().Index(index).
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
