package mongox

import (
	"context"
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
)

const writerMaxConcurrency = 100

// NewWriter 创建MongoDB日志写入器
func NewWriter(source, collection string) io.Writer {
	if Initialized() {
		client := GetClient(source)
		database := client.GetConfig().Database
		return &Writer{
			collection: client.GetClient().Database(database).Collection(collection),
			sem:        make(chan struct{}, writerMaxConcurrency),
		}
	}
	return nil
}

// Writer 日志写入
type Writer struct {
	collection *mongo.Collection
	sem        chan struct{} // 信号量，限制最大并发 goroutine 数
}

func (w *Writer) Write(p []byte) (int, error) {
	select {
	case w.sem <- struct{}{}:
		data := make([]byte, len(p))
		copy(data, p)
		go func() {
			defer func() { <-w.sem }()
			var doc interface{}
			if err := json.Unmarshal(data, &doc); err == nil {
				_, _ = w.collection.InsertOne(context.Background(), doc)
			}
		}()
	default:
		// 并发已达上限，丢弃本次写入，避免阻塞调用方
	}
	return len(p), nil
}
