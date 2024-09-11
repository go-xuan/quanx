package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/os/filex"
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
	FileOutput       = "file"
	LumberjackOutput = "lumberjack"
)

func New(app string) *Log {
	return &Log{FileName: app + ".log"}
}

// Log 日志配置
type Log struct {
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

// Info 配置信息格式化
func (l *Log) Info() string {
	return fmt.Sprintf("logPath=%s level=%s output=%s maxSize=%d maxAge=%d backups=%d",
		l.LogPath(), l.Level, l.Output, l.MaxSize, l.MaxAge, l.Backups)
}

// Title 配置器标题
func (*Log) Title() string {
	return "Log"
}

// Reader 配置文件读取
func (*Log) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "log.yaml",
		NacosDataId: "log.yaml",
	}
}

// Run 配置器运行
func (l *Log) Run() error {
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

func (l *Log) LogPath() string {
	return filepath.Join(l.Dir, l.FileName)
}

func (l *Log) Formatter() logrus.Formatter {
	host, _ := os.Hostname()
	return &LogFormatter{timeFormat: l.TimeFormat, host: host, Output: l.Output, useColor: l.UseColor}
}

func (l *Log) Writer() io.Writer {
	switch l.Output {
	case ConsoleOutput:
		return &ConsoleWriter{std: os.Stdout}
	case FileOutput:
		return &FileWriter{path: l.LogPath()}
	case LumberjackOutput:
		return &lumberjack.Logger{
			Filename:   l.LogPath(),
			MaxSize:    intx.IfZero(l.MaxSize, 100),
			MaxAge:     intx.IfZero(l.MaxAge, 7),
			MaxBackups: intx.IfZero(l.Backups, 10),
			Compress:   true,
		}
	default:
		return &FileWriter{path: l.LogPath()}
	}
}

// GetLevel 日志级别映射，默认debug
func (l *Log) GetLevel() logrus.Level {
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
