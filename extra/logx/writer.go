package logx

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/extra/elasticx"
	"github.com/go-xuan/quanx/extra/mongox"
)

func DefaultWriter() io.Writer {
	return &ConsoleWriter{
		writer: os.Stdout, // 标准输出
	}
}

func GetWriter(writerTo string, name, level string) io.Writer {
	switch writerTo {
	case WriterToFile:
		return NewFileWriter(name, level)
	case WriterToMongo:
		return NewMongoWriter(name)
	case WriterToES:
		return NewElasticSearchWriter(name)
	}
	return nil
}

func NewFileWriter(name string, level string) io.Writer {
	if level != "" {
		name = name + "_" + level
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(constx.DefaultResourceDir, "log", name+".log"),
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 10,
		Compress:   true,
	}
}

// NewMongoWriter 初始化mongo写入
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
