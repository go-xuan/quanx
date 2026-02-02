package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// NewWriter 创建日志写入器
func NewWriter(writer string, name string, level string) io.Writer {
	switch writer {
	case WriterConsole:
		return NewConsoleWriter()
	case WriterFile:
		if level != "" {
			name = fmt.Sprintf("%s.%s.log", name, level)
		} else {
			name = fmt.Sprintf("%s.log", name)
		}
		return NewFileWriter(filepath.Join("log", name))
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
