package logx

import (
	"github.com/go-xuan/quanx/file/filex"
	"io"
	"os"
)

// 日志输出
type Output struct {
	Writer io.Writer
}

func NewOutput(writer io.Writer) Output {
	return Output{Writer: writer}
}

type FileWriter struct {
	path string
}

func (w FileWriter) Write(bytes []byte) (n int, err error) {
	filex.CreateDirNotExist(w.path)
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
