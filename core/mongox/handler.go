package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the mongo handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func GetConfig() *Config {
	return this().GetConfig()
}

func Client() *mongo.Client {
	return this().GetClient()
}

type Handler struct {
	config *Config
	client *mongo.Client
}

func (h *Handler) GetConfig() *Config {
	return h.config
}

func (h *Handler) GetClient() *mongo.Client {
	return h.client
}

func (h *Handler) GetDatabaseNames(ctx context.Context) (dbs []string, err error) {
	dbs, err = h.client.ListDatabaseNames(ctx, bson.M{})
	return
}

func (h *Handler) GetCollection(collection string) (mc *mongo.Collection) {
	mc = h.client.Database(h.config.Database).Collection(collection)
	return
}
