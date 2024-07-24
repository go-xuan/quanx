package logx

import (
	"io"
	"os"

	"github.com/go-xuan/quanx/os/filex"
)

func DefaultWriter() io.Writer {
	return &ConsoleWriter{std: os.Stdout}
}

// ConsoleWriter 控制台输出
type ConsoleWriter struct {
	std io.Writer
}

func (w *ConsoleWriter) Write(bytes []byte) (n int, err error) {
	return w.std.Write(bytes)
}

// FileWriter 文件输出
type FileWriter struct {
	path string
}

func (w *FileWriter) Write(bytes []byte) (n int, err error) {
	filex.CreateIfNotExist(w.path)
	var file *os.File
	if file, err = os.OpenFile(w.path, filex.AppendOnly, 0666); err != nil {
		return
	}
	defer file.Close()
	if _, err = file.Write(bytes); err != nil {
		return
	}
	return
}
