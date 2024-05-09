package logx

import (
	"os"

	"github.com/go-xuan/quanx/file/filex"
)

// 日志输出
type FileWriter struct {
	path string
}

func (w *FileWriter) Write(bytes []byte) (n int, err error) {
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
