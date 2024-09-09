package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var handler *Handler

type Handler struct {
	config *Mongo
	client *mongo.Client
}

func GetConfig() *Mongo {
	return This().GetConfig()
}

func GetClient() *mongo.Client {
	return This().GetClient()
}

func This() *Handler {
	if handler == nil {
		panic("the mongo handler has not been initialized, please check the relevant config")
	}
	return handler
}

func (h *Handler) GetConfig() *Mongo {
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
