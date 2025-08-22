package logx

import (
	"io"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

var _config *Config

func GetConfig() *Config {
	return _config
}

func init() {
	log.SetOutput(NewConsoleWriter()) // 设置默认日志输出
	_config = new(Config)
	if err := (&configx.TagReader{Tag: "default"}).Read(_config); err != nil {
		panic(err)
	} else if err = _config.Execute(); err != nil {
		panic(err)
	}
}

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

// Config 日志配置
type Config struct {
	Name       string       `json:"name" yaml:"name" default:"app"`                                 // 日志文件名
	Level      string       `json:"level" yaml:"level" default:"info"`                              // 默认日志级别
	Formatter  string       `json:"formatter" yaml:"formatter" default:"json"`                      // 默认日志格式
	Writer     string       `json:"writer" yaml:"writer" default:"console"`                         // 默认日志输出
	TimeFormat string       `json:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05.999"` // 时间格式化
	Color      bool         `json:"color" yaml:"color" default:"false"`                             // 使用颜色
	Caller     bool         `json:"caller" yaml:"caller" default:"false"`                           // caller开关
	Hooks      []HookConfig `json:"hooks" yaml:"hooks"`                                             // 日志钩子
}

func (c *Config) NacosReader() configx.Reader {
	return &nacosx.Reader{
		DataId: "log.yaml",
	}
}

func (c *Config) FileReader() configx.Reader {
	return &configx.FileReader{
		Name: "log.yaml",
	}
}

func (c *Config) NeedRead() bool {
	if c.Name == "" && c.Level == "" && c.Formatter == "" {
		return true
	}
	return false
}

func (c *Config) Execute() error {
	// 添加hook钩子
	formatter := c.GetFormatter()
	if len(c.Hooks) > 0 {
		for _, hc := range c.Hooks {
			log.AddHook(hc.NewHook(c.Name, formatter))
		}
	}
	log.SetFormatter(formatter)        // 设置默认日志格式
	log.SetOutput(c.GetWriter())       // 设置默认日志输出
	log.SetLevel(LogrusLevel(c.Level)) // 设置默认日志级别
	log.SetReportCaller(c.Caller)      // 设置caller开关
	return nil
}

// GetFormatter 获取日志格式化器
func (c *Config) GetFormatter() log.Formatter {
	return NewFormatter(c.Formatter, c.TimeFormat)
}

// GetWriter 获取日志写入器
func (c *Config) GetWriter() io.Writer {
	if writer := NewWriter(c.Writer, c.Name); writer != nil {
		return writer
	}
	return NewConsoleWriter()
}

// LogrusLevel 转换日志级别
func LogrusLevel(level string) log.Level {
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
		return log.PanicLevel
	}
}
