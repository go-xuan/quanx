package configx

import (
	"fmt"
	"path/filepath"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
)

// SetFileReaderAnchor 设置本地文件读取器锚点
func SetFileReaderAnchor(anchor string) {
	if anchor != "" && filex.Exists(anchor) {
		defaultFileAnchor = anchor
	}
}

// NewFileReader 默认本地文件读取器
func NewFileReader(name string) *FileReader {
	return &FileReader{
		Name: name,
	}
}

// FileReader 本地文件读取器
type FileReader struct {
	Dir  string `json:"dir"`  // 文件路径
	Name string `json:"name"` // 文件名称
	Data []byte `json:"data"` // 文件内容
}

func (r *FileReader) Anchor(dir string) {
	if r.Dir == "" {
		r.Dir = dir
	}
}

func (r *FileReader) Location() string {
	return fmt.Sprintf("file@%s", r.GetPath())
}

// Read 读取配置文件
func (r *FileReader) Read(v any) error {
	if r.Data == nil {
		path := r.GetPath()
		if !filex.Exists(path) {
			return errorx.Sprintf("file not exist: %s", filex.Pwd(path))
		}
		data, err := filex.ReadFile(path)
		if err != nil {
			return errorx.Wrap(err, "file reader read error")
		}
		r.Data = data
	}
	if err := marshalx.Apply(r.Name).Unmarshal(r.Data, v); err != nil {
		return errorx.Wrap(err, "file reader unmarshal error")
	}
	return nil
}

// Write 写入配置文件
func (r *FileReader) Write(v any) error {
	if err := marshalx.Apply(r.Name).Write(r.GetPath(), v); err != nil {
		return errorx.Wrap(err, "file reader write error")
	}
	return nil
}

// GetPath 获取配置文件路径
func (r *FileReader) GetPath() string {
	r.Anchor(defaultFileAnchor)
	return filepath.Join(r.Dir, r.Name)
}
