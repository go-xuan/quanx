package configx

import (
	"path/filepath"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

type LocalReader struct {
	Dir  string `json:"dir" yaml:"dir"`   // 配置文件路径
	Name string `json:"name" yaml:"name"` // 配置文件名
}

func (r *LocalReader) Location(dir ...string) string {
	if len(dir) > 0 && dir[0] != "" {
		r.Dir = dir[0]
	}
	return filepath.Join(r.Dir, r.Name)
}

func (r *LocalReader) ReadConfig(config any) error {
	path, marshal := filepath.Join(r.Dir, r.Name), marshalx.Apply(r.Name)
	if err := marshal.Read(path, config); err != nil {
		return errorx.Wrap(err, "read config from local failed")
	}
	return nil
}
