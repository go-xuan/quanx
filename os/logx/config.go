package logx

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/intx"
)

// 日志级别
const (
	TraceLevel = "trace"
	DebugLevel = "debug"
	InfoLevel  = "info"
	ErrorLevel = "error"
	FatalLevel = "fatal"
	PanicLevel = "panic"

	TimeFormat = "2006-01-02 15:04:05.999"

	DefaultOutput    = "default"
	ConsoleOutput    = "console"
	LumberjackOutput = "lumberjack"
)

func NewConfigurator(conf *Config) configx.Configurator {
	return conf
}

// Config 日志配置
type Config struct {
	FileName   string            `json:"fileName" yaml:"fileName" default:"app.log"`                     // 日志文件名
	Dir        string            `json:"dir" yaml:"dir" default:"resource/log"`                          // 日志保存文件夹
	Level      string            `json:"level" yaml:"level" default:"info"`                              // 默认日志级别
	TimeFormat string            `json:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05.999"` // 时间格式化
	UseColor   bool              `json:"useColor" yaml:"useColor" default:"false"`                       // 使用颜色
	Output     string            `json:"output" yaml:"output" default:"default"`                         // 默认日志输出
	Outputs    map[string]string `json:"outputs" yaml:"outputs"`                                         // 自定义日志输出
	Caller     bool              `json:"caller" yaml:"caller" default:"false"`                           // caller开关
	MaxSize    int               `json:"maxSize" yaml:"maxSize" default:"100"`                           // 日志大小(单位：mb)
	MaxAge     int               `json:"maxAge" yaml:"maxAge" default:"1"`                               // 日志保留天数(单位：天)
	Backups    int               `json:"backups" yaml:"backups" default:"10"`                            // 日志备份数
}

func (c *Config) Format() string {
	return fmtx.Yellow.XSPrintf("logPath=%s level=%s output=%s maxSize=%v maxAge=%v backups=%v",
		c.LogPath(), c.Level, c.Output, c.MaxSize, c.MaxAge, c.Backups)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "log.yaml",
		NacosDataId: "log.yaml",
	}
}

func (c *Config) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	filex.CreateDir(c.Dir)
	// 添加hook
	log.AddHook(c.NewHook())
	// 更新formatter
	log.SetFormatter(c.Formatter())
	log.SetLevel(c.GetLogrusLevel())
	log.SetReportCaller(c.Caller)
	return nil
}

func (c *Config) LogPath(level ...string) string {
	filename := c.FileName
	if len(level) > 0 && level[0] != "" {
		filename = strings.Replace(filename, ".log", "-"+level[0]+".log", -1)
	}
	return filepath.Join(c.Dir, filename)
}

func (c *Config) Formatter() log.Formatter {
	host, _ := os.Hostname()
	return &LogFormatter{timeFormat: c.TimeFormat, host: host, Output: c.Output, useColor: c.UseColor}
}

func (c *Config) Writer(output ...string) io.Writer {
	op := c.Output
	if len(output) > 0 && output[0] != "" {
		op = output[0]
	}
	switch op {
	case ConsoleOutput:
		return &ConsoleWriter{std: os.Stdout}
	case LumberjackOutput:
		return &lumberjack.Logger{
			Filename:   c.LogPath(),
			MaxSize:    intx.IfZero(c.MaxSize, 100),
			MaxAge:     intx.IfZero(c.MaxAge, 7),
			MaxBackups: intx.IfZero(c.Backups, 10),
			Compress:   true,
		}
	case DefaultOutput:
		return &FileWriter{path: c.LogPath()}
	default:
		return &FileWriter{path: op}
	}
}

func (c *Config) NewHook() *Hook {
	hook := newHook()
	hook.InitWriter(c.Writer())
	hook.SetFormatter(c.Formatter())
	if c.Outputs != nil {
		for level, output := range c.Outputs {
			hook.SetWriter(ToLogrusLevel(level), c.Writer(output))
		}
	}
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
