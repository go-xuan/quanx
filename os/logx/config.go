package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/file/filex"
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
)

// 日志输出类型
const (
	ConsoleOutput = "console"
	FileOutput    = "file"
)

func New(app string) *LogConfig {
	return &LogConfig{FileName: app + ".log"}
}

// 日志配置
type LogConfig struct {
	FileName   string `json:"fileName" yaml:"fileName" default:"app.log"`                     // 日志文件名
	Dir        string `json:"dir" yaml:"dir" default:"resource/log"`                          // 日志保存文件夹
	Level      string `json:"level" yaml:"level" default:"info"`                              // 日志级别
	TimeFormat string `json:"timeFormat" yaml:"timeFormat" default:"2006-01-02 15:04:05.999"` // 时间格式化
	UseColor   bool   `json:"useColor" yaml:"useColor" default:"false"`                       // 使用颜色
	Output     string `json:"output" yaml:"output" default:"default"`                         // 日志输出
	Caller     bool   `json:"caller" yaml:"caller" default:"false"`                           // Flag for whether to caller
	MaxSize    int    `json:"maxSize" yaml:"maxSize" default:"100"`                           // 日志大小(单位：mb)
	MaxAge     int    `json:"maxAge" yaml:"maxAge" default:"1"`                               // 日志保留天数(单位：天)
	Backups    int    `json:"backups" yaml:"backups" default:"10"`                            // 日志备份数
}

// 配置信息格式化
func (l *LogConfig) Info() string {
	return fmt.Sprintf("logPath=%s level=%s output=%s maxSize=%d maxAge=%d backups=%d",
		l.LogPath(), l.Level, l.Output, l.MaxSize, l.MaxAge, l.Backups)
}

// 配置器标题
func (*LogConfig) Title() string {
	return "Log"
}

// 配置文件读取
func (*LogConfig) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "log.yaml",
		NacosDataId: "log.yaml",
	}
}

// 配置器运行
func (l *LogConfig) Run() error {
	if err := anyx.SetDefaultValue(l); err != nil {
		return err
	}
	filex.CreateDir(l.Dir)
	var writer, formatter = l.Writer(), l.Formatter()
	logrus.AddHook(NewHook(writer, formatter))
	logrus.SetFormatter(formatter)
	logrus.SetLevel(l.GetLevel())
	logrus.SetReportCaller(l.Caller)
	logrus.Info("Log Init Successful: ", l.Info())
	return nil
}

func (l *LogConfig) LogPath() string {
	return filepath.Join(l.Dir, l.FileName)
}

func (l *LogConfig) Formatter() *Formatter {
	return &Formatter{l.TimeFormat, l.UseColor}
}

func (l *LogConfig) Writer() io.Writer {
	switch l.Output {
	case ConsoleOutput:
		return os.Stdout
	case FileOutput:
		return &FileWriter{path: l.LogPath()}
	default:
		return &lumberjack.Logger{
			Filename:   l.LogPath(),
			MaxSize:    intx.IfZero(l.MaxSize, 100),
			MaxAge:     intx.IfZero(l.MaxAge, 7),
			MaxBackups: intx.IfZero(l.Backups, 10),
			Compress:   true,
		}
	}
}

// 日志级别映射，默认debug
func (l *LogConfig) GetLevel() logrus.Level {
	switch strings.ToLower(l.Level) {
	case TraceLevel:
		return logrus.TraceLevel
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}

func AllLevels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
