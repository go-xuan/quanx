package logx

import (
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/osx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/anyx"
)

// 日志级别
const (
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"

	WriterToConsole = "console" // 控制台打印
	WriterToFile    = "file"    // 写入日志文件
	WriterToMongo   = "mongo"   // 写入Mongo
	WriterToES      = "es"      // 写入ES

	FormatterText = "text" // 文本格式化
	FormatterJson = "json" // json格式化

	TimeFormat = "2006-01-02 15:04:05.999"

	logWriterSource = "log"
)

func init() {
	log.SetOutput(DefaultWriter()) // 设置默认日志输出
	if err := (&Config{}).Execute(); err != nil {
		panic(err)
	}
}

type HookWriterMapping map[string]string

// Config 日志配置
type Config struct {
	Name       string              `json:"name" yaml:"name" default:"app"`                                 // 日志文件名
	Level      string              `json:"level" yaml:"level" default:"info"`                              // 默认日志级别
	Formatter  string              `json:"formatter" yaml:"formatter" default:"json"`                      // 默认日志格式
	Writer     string              `json:"writer" yaml:"writer" default:"console"`                         // 默认日志输出
	Hook       []HookWriterMapping `json:"hook" yaml:"hook"`                                               // 日志钩子
	TimeFormat string              `json:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05.999"` // 时间格式化
	Color      bool                `json:"color" yaml:"color" default:"false"`                             // 使用颜色
	Caller     bool                `json:"caller" yaml:"caller" default:"false"`                           // caller开关
}

func (c *Config) Info() string {
	return fmt.Sprintf("level=%s formatter=%s writer=%s", c.Level, c.Formatter, c.Writer)
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
	// 添加hook钩子
	if len(c.Hook) > 0 {
		for _, mapping := range c.Hook {
			hook := NewHook()
			hook.SetFormatter(c.GetFormatter())
			for level, writerTo := range mapping {
				if writer := GetWriter(writerTo, c.Name, level); writer != nil {
					hook.SetWriter(ToLogrusLevel(level), writer)
				}
			}
			log.AddHook(hook)
		}
	}
	log.SetOutput(c.GetWriter())       // 设置默认日志输出
	log.SetFormatter(c.GetFormatter()) // 设置默认日志格式
	log.SetLevel(c.GetLogrusLevel())   // 设置默认日志级别
	log.SetReportCaller(c.Caller)      // 设置caller开关
	return nil
}

func (c *Config) GetFormatter() log.Formatter {
	switch c.Formatter {
	case FormatterJson:
		return &jsonFormatter{
			timeFormat: c.TimeFormat,
			hostname:   osx.Hostname(),
		}
	case FormatterText:
		return &textFormatter{
			timeFormat: TimeFormat,
			hostname:   osx.Hostname(),
			color:      c.Writer == WriterToConsole && c.Color,
		}
	default:
		return nil
	}
}

func (c *Config) GetWriter() io.Writer {
	switch c.Writer {
	case WriterToFile:
		return NewFileWriter(c.Name, "")
	case WriterToMongo:
		return NewMongoWriter(c.Name)
	case WriterToES:
		return NewElasticSearchWriter(c.Name)
	}
	return DefaultWriter()
}

func (c *Config) GetLogrusLevel() log.Level {
	return ToLogrusLevel(c.Level)
}

// ToLogrusLevel 日志级别映射，默认debug
func ToLogrusLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case LevelTrace:
		return log.TraceLevel
	case LevelDebug:
		return log.DebugLevel
	case LevelInfo:
		return log.InfoLevel
	case LevelError:
		return log.ErrorLevel
	case LevelFatal:
		return log.FatalLevel
	case LevelPanic:
		return log.PanicLevel
	default:
		return log.DebugLevel
	}
}
