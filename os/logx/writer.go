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
}

func (w *FileWriter) Write(bytes []byte) (int, error) {
	filex.CreateIfNotExist(w.path)
	file, err := os.OpenFile(w.path, filex.AppendOnly, 0666)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var n int
	if n, err = file.Write(bytes); err != nil {
		return n, err
	}
	return n, nil
}
