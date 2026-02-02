package mongox

import (
	"context"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

// NewWriter 创建MongoDB日志写入器
func NewWriter(source, coll string) io.Writer {
	if Initialized() {
		if db := GetDatabase(source); db != nil {
			collection := db.Collection(coll)
			return &Writer{collection: collection}
		}
	}
	return nil
}

// Writer 日志写入
type Writer struct {
	collection *mongo.Collection
}

func (w *Writer) Write(bytes []byte) (int, error) {
	go func() {
		var doc interface{}
		if err := json.Unmarshal(bytes, &doc); err == nil {
			_, _ = w.collection.InsertOne(context.Background(), doc)
		}
	}()
	return 0, nil
}
