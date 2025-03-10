package configx

import (
	"path/filepath"
	
	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

type From uint

const (
	FromLocal From = iota // 从本地文件读取配置
	FormNacos             // 从nacos读取配置
)

// Reader 配置读取器
type Reader interface {
	Location(...string) string // 配置文件定位
	ReadConfig(any) error      // 配置读取
}

// Configurator 配置器
type Configurator interface {
	Format() string          // 配置信息格式化
	Reader(from From) Reader // 配置文件读取
	Execute() error          // 配置器运行
}

type LocalFileReader struct {
	Dir  string `json:"dir" yaml:"dir"`   // 配置文件路径
	Name string `json:"name" yaml:"name"` // 配置文件名
}

func (r *LocalFileReader) Location(dir ...string) string {
	if len(dir) > 0 && dir[0] != "" {
		r.Dir = dir[0]
	}
	return filepath.Join(r.Dir, r.Name)
}

func (r *LocalFileReader) ReadConfig(v any) error {
	path, marshal := filepath.Join(r.Dir, r.Name), marshalx.Apply(r.Name)
	if err := marshal.Read(path, v); err != nil {
		return errorx.Wrap(err, "read config from file failed")
	}
	return nil
}

// ReadAndExecute 读取配置文件并运行
func ReadAndExecute(conf Configurator, from From, location ...string) error {
	if reader := conf.Reader(from); reader != nil {
		reader.Location(location...)
		if err := reader.ReadConfig(conf); err != nil {
			return errorx.Wrap(err, "read config failed")
		}
	}
	if err := conf.Execute(); err != nil {
		return errorx.Wrap(err, "configurator execute err")
	}
	return nil
}
