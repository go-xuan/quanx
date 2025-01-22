package logx

import (
	"io"
	"os"

	"github.com/go-xuan/quanx/os/filex"
)

func DefaultWriter() io.Writer {
	return &ConsoleWriter{
		std: os.Stdout, // 标准输出
	}
}

// ConsoleWriter 控制台输出
type ConsoleWriter struct {
	std io.Writer
}

func (w *ConsoleWriter) Write(bytes []byte) (int, error) {
	return w.std.Write(bytes)
}

// FileWriter 文件输出
type FileWriter struct {
	path string
	file *os.File
}

func (w *FileWriter) Write(bytes []byte) (int, error) {
	if w.file == nil {
		filex.CreateIfNotExist(w.path)
		if file, err := os.OpenFile(w.path, os.O_APPEND|os.O_WRONLY, 0644); err != nil {
			return 0, err
		} else {
			w.file = file
		}
	}
	if n, err := w.file.Write(bytes); err != nil {
		return n, err
	} else {
		return n, nil
	}
}
