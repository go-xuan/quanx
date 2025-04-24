package configx

import (
	"path/filepath"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/utils/marshalx"
)

type LocalReader struct {
	Dir  string `json:"dir" yaml:"dir"`   // 配置文件路径
	Name string `json:"name" yaml:"name"` // 配置文件名
	Data []byte `json:"data" yaml:"data"` // 配置文件内容
}

func (r *LocalReader) Ready(dir ...string) {
	if len(dir) > 0 && dir[0] != "" {
		r.Dir = dir[0]
	}
}

func (r *LocalReader) Check(config any) error {
	if r.Data == nil {
		path := filepath.Join(r.Dir, r.Name)
		if !filex.Exists(path) {
			return errorx.Errorf("config file not exist: %s", filex.Pwd(path))
		}
		data, err := filex.ReadFile(path)
		if err != nil {
			return errorx.Wrap(err, "read file error")
		}
		if err = marshalx.Apply(r.Name).Unmarshal(data, config); err != nil {
			return errorx.Wrap(err, "unmarshal config from local error")
		}
		r.Data = data
	}
	return nil
}

func (r *LocalReader) Read(config any) error {
	if err := r.Check(config); err != nil {
		return errorx.Wrap(err, "check local config error")
	}
	return nil
}

func (r *LocalReader) Location() string {
	return "local@" + filepath.Join(r.Dir, r.Name)
}
