package logx

import (
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志配置
type Config struct {
	Dir     string `json:"dir" yaml:"dir"  default:"temp/log"`   // 日志保存文件夹
	Name    string `json:"name" yaml:"name" default:"app"`       // 日志文件名
	Level   string `json:"level" yaml:"level" default:"debug"`   // 日志级别
	MaxSize int    `json:"maxSize" yaml:"maxSize" default:"100"` // 日志大小(单位：mb)
	MaxAge  int    `json:"maxAge" yaml:"maxAge" default:"1"`     // 日志保留天数(单位：天)
	Backups int    `json:"backups" yaml:"backups" default:"10"`  // 日志备份数
}

// 初始化日志
func InitLogger(conf *Config) {
	_ = os.Mkdir(conf.Dir, os.ModePerm)
	var logWriter = conf.defaultLogger()
	if conf.MaxSize != 0 {
		logWriter.MaxSize = conf.MaxSize
	}
	if conf.MaxAge != 0 {
		logWriter.MaxAge = conf.MaxAge
	}
	if conf.Backups != 0 {
		logWriter.MaxBackups = conf.Backups
	}
	format := &LogFormatter{}
	myHook := NewHook(WriterMap{
		logrus.TraceLevel: &logWriter,
		logrus.DebugLevel: &logWriter,
		logrus.InfoLevel:  &logWriter,
		logrus.WarnLevel:  &logWriter,
		logrus.ErrorLevel: &logWriter,
		logrus.FatalLevel: &logWriter,
		logrus.PanicLevel: &logWriter,
	}, format)
	logger := logrus.StandardLogger()
	logger.AddHook(myHook)
	// Flag for whether to log caller info (off by default)
	logger.SetReportCaller(true)
	logger.SetFormatter(format)
	logger.SetLevel(getLogrusLevel(conf.Level))
}

// 默认日志配置
func (conf *Config) defaultLogger() lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   path.Join(conf.Dir, conf.Name+".log"),
		MaxSize:    100,
		MaxAge:     1,
		MaxBackups: 10,
	}
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
