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

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/elasticx"
	"github.com/go-xuan/quanx/core/mongox"
	"github.com/go-xuan/quanx/types/intx"
)

func DefaultWriter() io.Writer {
	return &ConsoleWriter{
		writer: os.Stdout, // 标准输出
	}
}

// ConsoleWriter 日志写入控制台
type ConsoleWriter struct {
	writer io.Writer
}

func (w *ConsoleWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(bytes)
}

type FileWriterConfig struct {
	Name    string `json:"name" yaml:"name" default:"app"`        // 日志文件名
	Dir     string `json:"dir" yaml:"dir" default:"resource/log"` // 日志保存文件夹
	MaxSize int    `json:"maxSize" yaml:"maxSize" default:"100"`  // 日志大小(单位：MB)
	MaxAge  int    `json:"maxAge" yaml:"maxAge" default:"7"`      // 日志保留天数(单位：天)
	Backups int    `json:"backups" yaml:"backups" default:"10"`   // 日志备份数
}

func NewFileWriter(conf *FileWriterConfig, level ...string) io.Writer {
	name := conf.Name
	if len(level) > 0 && level[0] != "" {
		name = name + "_" + level[0]
	}
	return &lumberjack.Logger{
		Filename:   filepath.Join(conf.Dir, name) + ".log",
		MaxSize:    intx.IfZero(conf.MaxSize, 100),
		MaxAge:     intx.IfZero(conf.MaxAge, 7),
		MaxBackups: intx.IfZero(conf.Backups, 10),
		Compress:   true,
	}
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
		var log jsonLog
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		_, _ = w.client.Database(w.database).Collection(w.collection).InsertOne(context.Background(), log)
	}()
	return 0, nil
}

func NewMongoWriter(conf *mongox.Config, collection string) (io.Writer, error) {
	if conf != nil {
		source := "log_writer"
		conf.Source = source
		conf.Enable = true
		if err := configx.Execute(conf); err != nil {
			return nil, err
		}
		return &MongoWriter{
			database:   conf.Database,
			collection: collection,
			client:     mongox.GetClient(source),
		}, nil
	}
	return nil, nil
}

func NewElasticSearchWriter(conf *elasticx.Config, index string) (io.Writer, error) {
	if elasticx.Initialized() {
		return &ElasticSearchWriter{
			index:  index,
			client: elasticx.Client(),
		}, nil
	} else if conf != nil {
		if err := configx.Execute(conf); err != nil {
			return nil, err
		}
		return &ElasticSearchWriter{
			index:  index,
			client: elasticx.Client(),
		}, nil
	}
	return nil, nil
}

// ElasticSearchWriter 日志写入elastic search
type ElasticSearchWriter struct {
	index  string
	client *elastic.Client
}

// 异步写入es
func (w *ElasticSearchWriter) Write(bytes []byte) (int, error) {
	go func() {
		var log jsonLog
		if err := json.Unmarshal(bytes, &log); err != nil {
			return
		}
		_, _ = w.client.Index().Index(w.index).BodyJson(log).Do(context.Background())
	}()
	return 0, nil
}
