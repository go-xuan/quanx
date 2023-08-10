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

func (ctl *Controller) Query(ctx context.Context, index string, query elastic.Query) *elastic.SearchResult {
	var result *elastic.SearchResult
	var err error
	result, err = ctl.Client.Search().Index(index).Query(query).Do(ctx)
	if err != nil {
		log.Error(err)
	}
	return result
}
