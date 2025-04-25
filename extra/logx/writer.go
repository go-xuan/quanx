package logx

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/extra/elasticx"
	"github.com/go-xuan/quanx/extra/mongox"
)

// NewConsoleWriter 创建控制台日志写入器
func NewConsoleWriter() io.Writer {
	return &ConsoleWriter{
		writer: os.Stdout, // 标准输出
	}
}

// NewFileWriter 创建本地文件日志写入器
func NewFileWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    100,      // 日志文件最大大小（MB）
		MaxAge:     7,        // 日志保留天数
		MaxBackups: 10,       // 日志备份数量
		Compress:   true,     // 是否压缩
	}
}

// NewMongoWriter 创建MongoDB日志写入器
func NewMongoWriter(collection string) io.Writer {
	if mongox.Initialized() {
		client := mongox.GetClient(logWriterSource)
		return &MongoWriter{
			database:   client.Config().Database,
			collection: collection,
			client:     client.Instance(),
		}
	}
	return nil
}

// NewElasticSearchWriter 创建ES日志写入器
func NewElasticSearchWriter(index string) io.Writer {
	if elasticx.Initialized() {
		client := elasticx.GetClient(logWriterSource)
		if exist, err := client.Instance().IndexExists(index).Do(context.TODO()); err != nil || !exist {
			_, _ = client.CreateIndex(context.TODO(), index)
		}
		return &ElasticSearchWriter{
			index:  index,
			client: client.Instance(),
		}
	}
	return nil
}

// ConsoleWriter 日志写入控制台
type ConsoleWriter struct {
	writer io.Writer
}

func (w *ConsoleWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(bytes)
}

// MongoWriter 日志写入mongo
type MongoWriter struct {
	database   string
	collection string
	client     *mongo.Client
}

// 异步写入mongo
func (w *MongoWriter) Write(bytes []byte) (int, error) {
	go func() {
		var log LogRecord
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		_, _ = w.client.Database(w.database).Collection(w.collection).InsertOne(context.Background(), log)
	}()
	return 0, nil
}

// ElasticSearchWriter 日志写入elastic search
type ElasticSearchWriter struct {
	index  string
	client *elastic.Client
}

// 异步写入es
func (w *ElasticSearchWriter) Write(bytes []byte) (int, error) {
	go func() {
		var log LogRecord
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		_, _ = w.client.Index().Index(w.index).BodyJson(log).Do(context.Background())
	}()
	return 0, nil
}
