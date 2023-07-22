package elasticx

import (
	"context"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var CTL *Control

func InitEsCTL(conf *Config) {
	if conf.Host == "" {
		return
	}
	var err error
	msg := conf.Format()
	if CTL == nil {
		CTL, err = conf.NewEsCtl()
		if err == nil && CTL.Ping() {
			log.Info("初始化ElasticSearch连接-成功! ", msg)
		} else {
			log.Error("初始化ElasticSearch连接-失败! ", msg)
			log.Error("error : ", err)
		}
	} else {
		var client *elastic.Client
		client, err = conf.NewClient()
		if err == nil && CTL.Ping() {
			CTL.Client = client
			CTL.Config = conf
			log.Error("更新ElasticSearch连接-成功! ", msg)
		} else {
			log.Error("更新ElasticSearch连接-失败! ", msg)
			log.Error("error : ", err)
		}
	}
}

// ES控制器
type Control struct {
	Config *Config
	Url    string
	Client *elastic.Client
	Ctx    context.Context
}

func (ctl *Control) Ping() bool {
	info, code, err := ctl.Client.Ping(ctl.Url).Do(ctl.Ctx)
	if err != nil || code == 200 {
		log.Info("elastic-search version is ", info.Version.Number)
		return true
	}
	return false
}

func (ctl *Control) Create(index, id string, body interface{}) {
	var result *elastic.IndexResponse
	var err error
	result, err = ctl.Client.Index().Index(index).Id(id).BodyJson(body).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	log.Printf("create success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
}

func (ctl *Control) Update(index, id string, body interface{}) {
	var result *elastic.UpdateResponse
	var err error
	result, err = ctl.Client.Update().Index(index).Id(id).Doc(body).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	log.Printf("update success: id=%s index=%s type=%s\n", result.Id, result.Index, result.Type)
}

func (ctl *Control) Delete(index, id string) {
	var result *elastic.DeleteResponse
	var err error
	result, err = ctl.Client.Delete().Index(index).Id(id).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	log.Printf("delete success: %s\n", result.Result)
}

func (ctl *Control) Get(index, id string) *elastic.GetResult {
	var result *elastic.GetResult
	var err error
	result, err = ctl.Client.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	return result
}

func (ctl *Control) Query(index string, query elastic.Query) *elastic.SearchResult {
	var result *elastic.SearchResult
	var err error
	result, err = ctl.Client.Search().Index(index).Query(query).Do(context.Background())
	if err != nil {
		log.Error(err)
	}
	return result
}
