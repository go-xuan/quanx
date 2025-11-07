package mongox

import (
	"context"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

// NewLogWriter 创建MongoDB日志写入器
func NewLogWriter[T any](source, collection string) io.Writer {
	if Initialized() {
		client := GetClient(source)
		return &LogWriter[T]{
			database:   client.GetConfig().Database,
			collection: collection,
			client:     client.GetInstance(),
		}
	}
	return nil
}

// LogWriter 日志写入
type LogWriter[T any] struct {
	database   string
	collection string
	client     *mongo.Client
}

func (w *LogWriter[T]) Write(bytes []byte) (int, error) {
	// 异步写入
	go func() {
		var log T
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		_, _ = w.client.Database(w.database).Collection(w.collection).InsertOne(context.Background(), log)
	}()
	return 0, nil
}
