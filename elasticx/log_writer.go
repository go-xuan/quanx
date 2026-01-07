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
		client := GetClient(source)
		ctx, ins := context.Background(), client.GetInstance()
		if exist, err := ins.IndexExists(index).Do(ctx); err != nil || !exist {
			_, _ = client.CreateIndex(ctx, index)
		}
		return &LogWriter[T]{
			index:  index,
			client: ins,
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
		_, _ = w.client.Index().Index(w.index).BodyJson(log).Do(context.Background())
	}()
	return 0, nil
}
