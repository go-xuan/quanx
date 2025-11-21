package logx

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
	
	"github.com/go-xuan/quanx/elasticx"
	"github.com/go-xuan/quanx/mongox"
)

// NewWriter 创建日志写入器
func NewWriter(writer string, name string, level ...string) io.Writer {
	switch writer {
	case WriterConsole:
		return NewConsoleWriter()
	case WriterFile:
		if len(level) > 0 && level[0] != "" {
			name = name + "_" + level[0]
		}
		name = filepath.Join("log", name+".log")
		return NewFileWriter(name)
	case WriterMongo:
		return mongox.NewLogWriter[LogRecord](logWriterSource, name)
	case WriterElasticSearch:
		return elasticx.NewLogWriter[LogRecord](logWriterSource, name)
	}
	return nil
}

// NewConsoleWriter 创建控制台日志写入器
func NewConsoleWriter() io.Writer {
	return &ConsoleWriter{}
}

// ConsoleWriter 日志写入控制台
type ConsoleWriter struct{}

func (w *ConsoleWriter) Write(bytes []byte) (int, error) {
	return os.Stdout.Write(bytes)
}

// NewFileWriter 创建本地文件日志写入器
func NewFileWriter(filename string) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    100,      // 日志文件最大大小（MB）
		MaxAge:     7,        // 日志保留天数
		MaxBackups: 10,       // 日志备份数量
		Compress:   true,     // 是否压缩
	}
}
