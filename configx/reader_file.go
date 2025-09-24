package configx

import (
	"path/filepath"

	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
)

// NewFileReader 默认本地文件读取器
func NewFileReader(name string) *FileReader {
	return &FileReader{
		Name: name,
	}
}

// FileReader 本地文件读取器
type FileReader struct {
	Dir  string `json:"dir"`  // 配置文件路径
	Name string `json:"name"` // 配置文件名
	Data []byte `json:"data"` // 配置文件内容
}

func (r *FileReader) Anchor(dir string) {
	if r.Dir == "" {
		r.Dir = dir
	}
}

func (r *FileReader) Read(v any) error {
	if r.Data == nil {
		r.Anchor(constx.DefaultConfigDir)
		path := filepath.Join(r.Dir, r.Name)
		if !filex.Exists(path) {
			return errorx.Errorf("file not exist: %s", filex.Pwd(path))
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

func (r *FileReader) Location() string {
	return "local@" + filepath.Join(r.Dir, r.Name)
}

func (r *FileReader) Write(v any) error {
	if r.Data == nil {
		data, err := marshalx.Apply(r.Name).Marshal(v)
		if err != nil {
			return errorx.Wrap(err, "file reader marshal error")
		}
		r.Data = data
	}
	r.Anchor(constx.DefaultConfigDir)
	path := filepath.Join(r.Dir, r.Name)
	if err := filex.WriteFile(path, r.Data); err != nil {
		return errorx.Wrap(err, "file reader write error")
	}
	return nil
}
