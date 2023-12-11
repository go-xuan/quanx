package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var instance *Handler

type Handler struct {
	Config *Config
	Client *mongo.Client
}

func This() *Handler {
	if instance == nil {
		panic("The mongo instance has not been initialized, please check the relevant config")
	}
	return instance
}

func (h *Handler) GetDatabaseNames(ctx context.Context) (dbs []string, err error) {
	dbs, err = h.Client.ListDatabaseNames(ctx, bson.M{})
	return
}

func (h *Handler) GetCollection(collection string) (mc *mongo.Collection) {
	mc = h.Client.Database(h.Config.Database).Collection(collection)
	return
}
