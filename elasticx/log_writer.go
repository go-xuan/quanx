package elasticx

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
)

// NewLogWriter 创建ES日志写入器
func NewLogWriter[T any](source, index string) io.Writer {
	if Initialized() {
		client := GetESClient(source)
		ctx := context.Background()
		if exist, err := client.IndexExists(index).Do(ctx); err != nil || !exist {
			_, _ = client.CreateIndex(index).Do(ctx)
		}
		return &LogWriter[T]{
			index:  index,
			client: client,
		}
	}
	return nil
}

// LogWriter 日志写入
type LogWriter[T any] struct {
	index  string
	client *elastic.Client
}

func (w *LogWriter[T]) Write(bytes []byte) (int, error) {
	go func() {
		var log T
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		ctx := context.Background()
		_, _ = w.client.Index().Index(w.index).BodyJson(log).Do(ctx)
	}()
	return 0, nil
}
