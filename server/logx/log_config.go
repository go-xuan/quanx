package logx

import (
	"fmt"
	"github.com/go-xuan/quanx/types/intx"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/go-xuan/quanx/server/confx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/filex"
)

func New(app string) *LogConfig {
	return &LogConfig{FileName: app + ".log"}
}

// 日志配置
type LogConfig struct {
	FileName string `json:"fileName" yaml:"fileName" default:"app.log"` // 日志文件名
	Dir      string `json:"dir" yaml:"dir" default:"resource/log"`      // 日志保存文件夹
	Level    string `json:"level" yaml:"level" default:"debug"`         // 日志级别
	MaxSize  int    `json:"maxSize" yaml:"maxSize" default:"100"`       // 日志大小(单位：mb)
	MaxAge   int    `json:"maxAge" yaml:"maxAge" default:"1"`           // 日志保留天数(单位：天)
	Backups  int    `json:"backups" yaml:"backups" default:"10"`        // 日志备份数
}

// 配置信息格式化
func (l *LogConfig) ToString() string {
	return fmt.Sprintf("logPath=%s level=%s maxSize=%d maxAge=%d backups=%d",
		l.LogPath(), l.Level, l.MaxSize, l.MaxAge, l.Backups)
}

// 配置器名称
func (*LogConfig) Theme() string {
	return "Log"
}

// 配置文件读取
func (*LogConfig) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath: "log.yaml",
	}
}

// 配置器运行
func (l *LogConfig) Run() error {
	if err := anyx.SetDefaultValue(l); err != nil {
		return err
	}
	filex.CreateDir(l.Dir)
	var logWriter = &lumberjack.Logger{
		Filename:   l.LogPath(),
		MaxSize:    intx.IfZero(l.MaxSize, 100),
		MaxAge:     intx.IfZero(l.MaxSize, 1),
		MaxBackups: intx.IfZero(l.MaxSize, 10),
	}
	var format = &LogFormatter{}
	var hook = NewHook(WriterMap{
		logrus.TraceLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}, format)
	var logger = logrus.StandardLogger()
	logger.AddHook(hook)
	// Flag for whether to l caller info (off by default)
	logger.SetReportCaller(true)
	logger.SetFormatter(format)
	logger.SetLevel(getLogrusLevel(l.Level))
	return nil
}

func (l *LogConfig) LogPath() string {
	return filepath.Join(l.Dir, l.FileName)
}

// 日志级别
const (
	Trace = "trace"
	Debug = "debug"
	Info  = "info"
	Error = "error"
	Fatal = "fatal"
	Panic = "panic"
)

// 日志级别映射，默认debug
func getLogrusLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case Trace:
		return logrus.TraceLevel
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Error:
		return logrus.ErrorLevel
	case Fatal:
		return logrus.FatalLevel
	case Panic:
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
