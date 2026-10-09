package elasticx

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
)

const writerMaxConcurrency = 100

// NewWriter 创建ES日志写入器
func NewWriter(source, index string) io.Writer {
	if Initialized() {
		client := GetInstance(source)
		ctx := context.Background()
		if exist, err := client.IndexExists(index).Do(ctx); err != nil || !exist {
			_, _ = client.CreateIndex(index).Do(ctx)
		}
		return &Writer{
			index:  index,
			client: client,
			sem:    make(chan struct{}, writerMaxConcurrency),
		}
	}
	return nil
}

// Writer 日志写入
type Writer struct {
	index  string
	client *elastic.Client
	sem    chan struct{} // 信号量，限制最大并发 goroutine 数
}

func (w *Writer) Write(p []byte) (int, error) {
	select {
	case w.sem <- struct{}{}:
		data := make([]byte, len(p))
		copy(data, p)
		go func() {
			defer func() { <-w.sem }()
			var body interface{}
			if err := json.Unmarshal(data, &body); err == nil {
				_, _ = w.client.Index().Index(w.index).BodyJson(body).Do(context.Background())
			}
		}()
	default:
		// 并发已达上限，丢弃本次写入，避免阻塞调用方
	}
	return len(p), nil
}