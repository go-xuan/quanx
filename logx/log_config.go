package logx

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-xuan/quanx/utilx/structx"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志配置
type Config struct {
	AppName string `json:"appName" yaml:"appName" default:"app"`  // 服务名
	Dir     string `json:"dir" yaml:"dir" default:"resource/log"` // 日志保存文件夹
	Level   string `json:"level" yaml:"level" default:"debug"`    // 日志级别
	MaxSize int    `json:"maxSize" yaml:"maxSize" default:"100"`  // 日志大小(单位：mb)
	MaxAge  int    `json:"maxAge" yaml:"maxAge" default:"1"`      // 日志保留天数(单位：天)
	Backups int    `json:"backups" yaml:"backups" default:"10"`   // 日志备份数
}

func (config *Config) Format() string {
	return fmt.Sprintf("appName=%s dir=%s level=%s maxSize=%d maxAge=%d backups=%d",
		config.AppName, config.Dir, config.Level, config.MaxSize, config.MaxAge, config.Backups)
}

func (config *Config) Init() {
	InitLogX(config)
}

// 初始化日志
func InitLogX(conf *Config) {
	_ = structx.SetDefaultValue(conf)
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
	hook := NewHook(WriterMap{
		logrus.TraceLevel: &logWriter,
		logrus.DebugLevel: &logWriter,
		logrus.InfoLevel:  &logWriter,
		logrus.WarnLevel:  &logWriter,
		logrus.ErrorLevel: &logWriter,
		logrus.FatalLevel: &logWriter,
		logrus.PanicLevel: &logWriter,
	}, format)
	var logger = logrus.StandardLogger()
	logger.AddHook(hook)
	// Flag for whether to log caller info (off by default)
	logger.SetReportCaller(true)
	logger.SetFormatter(format)
	logger.SetLevel(getLogrusLevel(conf.Level))
	return
}

// 默认日志配置
func (config *Config) defaultLogger() lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   path.Join(config.Dir, config.AppName+".log"),
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
