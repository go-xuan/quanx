package logx

import (
	"io"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/osx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// 日志级别
const (
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"

	WriterConsole = "console" // 控制台打印
	WriterFile    = "file"    // 写入日志文件

	FormatterText = "text"                    // 文本格式化
	FormatterJson = "json"                    // json格式化
	TimeLayout    = "2006-01-02 15:04:05.000" // 时间格式化
)

var _config *Config // 日志配置

func init() {
	// 设置默认日志输出
	log.SetOutput(NewConsoleWriter())
	// 初始化日志配置
	errorx.Panic(configx.LoadConfigurator(&Config{}))
}

// GetConfig 获取日志配置
func GetConfig() *Config {
	return _config
}

// Config 日志配置
type Config struct {
	Name      string       `json:"name" yaml:"name" default:"app"`            // 日志文件名
	Level     string       `json:"level" yaml:"level" default:"info"`         // 默认日志级别
	Formatter string       `json:"formatter" yaml:"formatter" default:"json"` // 默认日志格式
	Writer    string       `json:"writer" yaml:"writer" default:"console"`    // 默认日志输出
	Color     bool         `json:"color" yaml:"color" default:"false"`        // 使用颜色
	Caller    bool         `json:"caller" yaml:"caller" default:"false"`      // caller开关
	Hooks     []HookConfig `json:"hooks" yaml:"hooks"`                        // 日志钩子
}

// HookConfig 日志钩子配置
type HookConfig struct {
	Writer string   `json:"writer" yaml:"writer" default:"console"`
	Levels []string `json:"levels" yaml:"levels"`
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["name"] = c.Name
	fields["level"] = c.Level
	fields["formatter"] = c.Formatter
	fields["writer"] = c.Writer
	fields["color"] = c.Color
	fields["caller"] = c.Caller
	return fields
}

func (c *Config) Valid() bool {
	return c.Level != "" && c.Formatter != "" && c.Writer != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("log.yaml"),
		configx.NewFileReader("log.yaml"),
		configx.NewTagReader(),
	}
}

func (c *Config) Execute() error {
	formatter := c.GetFormatter()
	if len(c.Hooks) > 0 {
		for _, hc := range c.Hooks {
			hook := NewHook()
			hook.SetFormatter(formatter)
			hook.SetCaller(c.Caller)
			for _, level := range hc.Levels {
				if writer := NewWriter(hc.Writer, c.Name, level); writer != nil {
					hook.AddWriter(LogrusLevel(level), writer)
				}
			}
			log.AddHook(hook)
		}
	}

	log.SetFormatter(formatter)        // 设置默认日志格式
	log.SetLevel(LogrusLevel(c.Level)) // 设置默认日志级别
	log.SetOutput(c.NewWriter())       // 设置默认日志输出
	log.SetReportCaller(c.Caller)      // 设置caller开关

	log.WithFields(c.LogFields()).Info("init log success")
	_config = c
	return nil
}

func (c *Config) GetFormatter() log.Formatter {
	return &Formatter{
		Formatter:  c.Formatter,
		TimeLayout: TimeLayout,
		Hostname:   osx.Hostname(),
		Color:      c.Color && c.Writer == WriterConsole,
	}
}

func (c *Config) NewWriter() io.Writer {
	if writer := NewWriter(c.Writer, c.Name, c.Level); writer != nil {
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
