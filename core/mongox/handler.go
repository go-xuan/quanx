package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/os/errorx"
)

var _handler *Handler

func Initialized() bool {
	return _handler != nil
}

func this() *Handler {
	if _handler == nil {
		panic("the mongo handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

func GetClient(source ...string) *mongo.Client {
	return this().GetClient(source...)
}

type Handler struct {
	multi   bool // 是否多数据源连接
	config  *Config
	configs map[string]*Config
	client  *mongo.Client
	clients map[string]*mongo.Client
}

func (h *Handler) GetConfig(source ...string) *Config {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) GetClient(source ...string) *mongo.Client {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if db, ok := h.clients[source[0]]; ok {
			return db
		}
	}
	return h.client
}

func (h *Handler) GetDatabaseNames(ctx context.Context) ([]string, error) {
	if dbs, err := h.client.ListDatabaseNames(ctx, bson.M{}); err != nil {
		return nil, errorx.Wrap(err, "get mongo db names failed")
	} else {
		return dbs, nil
	}
}

func (h *Handler) GetCollection(collection string) *mongo.Collection {
	return h.client.Database(h.config.Database).Collection(collection)
}

func (h *Handler) InsertOne(ctx context.Context, collection string, document any) (*mongo.InsertOneResult, error) {
	if res, err := h.GetCollection(collection).InsertOne(ctx, document); err != nil {
		return nil, errorx.Wrap(err, "insert mongo collection failed")
	} else {
		return res, nil
	}
}
