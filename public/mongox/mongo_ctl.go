package mongox

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CTL *Controller

type Controller struct {
	Config *Config
	Client *mongo.Client
}

func Init(conf *Config) {
	var client = conf.NewClient()
	if err := client.Ping(context.Background(), nil); err == nil {
		CTL = &Controller{Config: conf, Client: client}
		log.Error("MongoDB连接成功！", conf.Format())
	} else {
		log.Error("MongoDB连接失败！", conf.Format())
		log.Error("error : ", err)
	}
}

func (ctl *Controller) GetDatabaseNames(ctx context.Context) (dbs []string, err error) {
	dbs, err = ctl.Client.ListDatabaseNames(ctx, bson.M{})
	return
}

func (ctl *Controller) GetCollection(collection string) (mc *mongo.Collection) {
	mc = ctl.Client.Database(ctl.Config.Database).Collection(collection)
	return
}
