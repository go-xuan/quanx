package elasticx

import (
	"context"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var CTL *Control

// ES控制器
type Control struct {
	Config *Config
	Url    string
	Client *elastic.Client
	Ctx    context.Context
}

func (ctl *Control) Ping() bool {
	info, code, err := ctl.Client.Ping(ctl.Url).Do(ctl.Ctx)
	if err != nil || code != 200 {
		log.Info("elastic-search version is ", info.Version.Number)
		return true
	}
	return false
}

func InitEsCTL(conf *Config) {
	if conf.Host == "" {
		return
	}
	var err error
	msg := conf.Format()
	if CTL == nil {
		CTL, err = conf.NewEsCtl()
		if err != nil {
			log.Error("初始化ElasticSearch连接-失败! ", msg)
			log.Error("error : ", err)
		} else {
			log.Info("初始化ElasticSearch连接-成功! ", msg)
		}
	} else {
		var client *elastic.Client
		client, err = conf.NewClient()
		if err != nil {
			log.Error("更新ElasticSearch连接-失败! ", msg)
			log.Error("error : ", err)
		} else {
			CTL.Client = client
			CTL.Config = conf
			log.Error("更新ElasticSearch连接-成功! ", msg)
		}
	}
}
