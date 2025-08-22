package configx

import (
	"path/filepath"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
)

// FileReader 本地文件读取器
type FileReader struct {
	Dir  string // 配置文件路径
	Name string // 配置文件名
	Data []byte // 配置文件内容
}

// Anchor 本地文件读取器锚点为配置文件夹路径
func (r *FileReader) Anchor(anchor string) {
	if r.Dir == "" {
		r.Dir = anchor
	}
}

// Location 配置文件位置
func (r *FileReader) Location() string {
	return "local@" + filepath.Join(r.Dir, r.Name)
}

// Read 从本地文件中读取配置
func (r *FileReader) Read(v any) error {
	if r.Data == nil {
		if path := filepath.Join(r.Dir, r.Name); !filex.Exists(path) {
			return errorx.Errorf("local file not exist: %s", filex.Pwd(path))
		} else if data, err := filex.ReadFile(path); err != nil {
			return errorx.Wrap(err, "read file error")
		} else if err = marshalx.Apply(r.Name).Unmarshal(data, v); err != nil {
			return errorx.Wrap(err, "read config from local error")
		} else {
			r.Data = data
		}
	}
	return nil
}
