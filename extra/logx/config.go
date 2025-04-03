package logx

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/anyx"
)

// 日志级别
const (
	TraceLevel = "trace"
	DebugLevel = "debug"
	InfoLevel  = "info"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"

	defaultWriterType = "default" // Console
	fileWriterType    = "file"    // file
	mongoWriterType   = "mongo"   // Mongo
	eSOutWriterType   = "es"      // Elasticsearch

	logWriterSource = "log"
)

// Config 日志配置
type Config struct {
	Name       string            `json:"name" yaml:"name" default:"app"`                                 // 日志文件名
	Level      string            `json:"level" yaml:"level" default:"info"`                              // 默认日志级别
	Formatter  string            `json:"formatter" yaml:"formatter" default:"json"`                      // 日志格式
	Writer     string            `json:"writer" yaml:"writer" default:"file"`                            // 默认日志输出
	Writers    map[string]string `json:"writers" yaml:"writers"`                                         // 日志级别日志输出
	TimeFormat string            `json:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05.999"` // 时间格式化
	Color      bool              `json:"color" yaml:"color" default:"false"`                             // 使用颜色
	Caller     bool              `json:"caller" yaml:"caller" default:"false"`                           // caller开关
	File       *FileWriterConfig `json:"file" yaml:"file"`                                               // 日志输出到文件
}

func (c *Config) Format() string {
	return fmt.Sprintf("level=%s formatter=%s writer=%s",
		c.Level, c.Formatter, c.Writer)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "log.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalReader{
			Name: "log.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	c.initFile()
	log.AddHook(c.NewHook())           // 添加hook
	log.SetFormatter(c.LogFormatter()) // 设置formatter
	log.SetLevel(c.GetLogrusLevel())   // 设置默认日志级别
	log.SetReportCaller(c.Caller)
	return nil
}

func (c *Config) initFile() {
	var needFile bool
	if c.Writer == fileWriterType {
		needFile = true
	} else if len(c.Writers) > 0 {
		for _, writer := range c.Writers {
			if writer == fileWriterType {
				needFile = true
			}
		}
	}
	if needFile && c.File == nil {
		c.File = &FileWriterConfig{Name: c.Name}
	}
	_ = anyx.SetDefaultValue(c.File)
}

func (c *Config) LogFormatter() log.Formatter {
	host, _ := os.Hostname()
	switch c.Formatter {
	case "json":
		return &jsonFormatter{
			timeFormat: c.TimeFormat,
			hostname:   host,
		}
	default:
		return &textFormatter{
			timeFormat: c.TimeFormat,
			hostname:   host,
			color:      c.Writer == defaultWriterType && c.Color,
		}
	}
}

func (c *Config) NewWriter() io.Writer {
	switch c.Writer {
	case fileWriterType:
		return NewFileWriter(c.File)
	case mongoWriterType:
		if writer, err := NewMongoWriter(c.Name); writer != nil && err == nil {
			return writer
		}
	case eSOutWriterType:
		if writer, err := NewElasticSearchWriter(c.Name); writer != nil && err == nil {
			return writer
		}
	}
	return DefaultWriter()
}

func (c *Config) NewWriters() map[log.Level]io.Writer {
	var writers = make(map[log.Level]io.Writer)
	for lv, writerType := range c.Writers {
		if lv == c.Level {
			continue
		}
		level := ToLogrusLevel(lv)
		switch writerType {
		case fileWriterType:
			writers[level] = NewFileWriter(c.File, lv)
		case mongoWriterType:
			if writer, err := NewMongoWriter(c.Name); writer != nil && err == nil {
				writers[level] = writer
			}
		case eSOutWriterType:
			if writer, err := NewElasticSearchWriter(c.Name); writer != nil && err == nil {
				writers[level] = writer
			}
		}
		// 保底
		if _, ok := writers[level]; !ok {
			writers[level] = DefaultWriter()
		}
	}
	return writers
}

func (c *Config) NewHook() *Hook {
	hook := newHook()
	hook.InitWriter(c.NewWriter())
	hook.SetFormatter(c.LogFormatter())
	hook.SetWriters(c.NewWriters())
	return hook
}

func (c *Config) GetLogrusLevel() log.Level {
	return ToLogrusLevel(c.Level)
}

// ToLogrusLevel 日志级别映射，默认debug
func ToLogrusLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case TraceLevel:
		return log.TraceLevel
	case DebugLevel:
		return log.DebugLevel
	case InfoLevel:
		return log.InfoLevel
	case ErrorLevel:
		return log.ErrorLevel
	case FatalLevel:
		return log.FatalLevel
	case PanicLevel:
		return log.PanicLevel
	default:
		return log.DebugLevel
	}
}

// AllLogrusLevels 所有日志级别
func AllLogrusLevels() []log.Level {
	return []log.Level{
		log.TraceLevel,
		log.DebugLevel,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
		log.PanicLevel,
	}
}
