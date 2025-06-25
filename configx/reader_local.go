package configx

import (
	"path/filepath"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
)

// FileReader 本地文件读取器
type FileReader struct {
	Dir  string `json:"dir" yaml:"dir"`   // 配置文件路径
	Name string `json:"name" yaml:"name"` // 配置文件名
	Data []byte `json:"data" yaml:"data"` // 配置文件内容
}

func (r *FileReader) Anchor(anchor string) {
	r.Dir = anchor
}

func (r *FileReader) Location() string {
	return "local@" + filepath.Join(r.Dir, r.Name)
}

func (r *FileReader) Read(config any) error {
	if r.Data == nil {
		if path := filepath.Join(r.Dir, r.Name); !filex.Exists(path) {
			return errorx.Errorf("local file not exist: %s", filex.Pwd(path))
		} else if data, err := filex.ReadFile(path); err != nil {
			return errorx.Wrap(err, "read file error")
		} else if err = marshalx.Apply(r.Name).Unmarshal(data, config); err != nil {
			return errorx.Wrap(err, "read config from local error")
		} else {
			r.Data = data
		}
	}
	return nil
}
