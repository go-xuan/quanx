package nacosx

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/file/filex"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// 配置项
type Configs []*Config

// 批量加载nacos配置
func (list Configs) Loading(v any) (err error) {
	for _, conf := range list {
		if err = conf.Loading(v); err != nil {
			return
		}
	}
	return
}

// 根据配置名获取配置
func (list Configs) Get(dataId string) (target *Config) {
	for _, conf := range list {
		if conf.DataId == dataId {
			target = conf
			return
		}
	}
	return
}

type Config struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	Listen bool          `yaml:"listen"` // 是否启用监听
}

// 初始化
func NewConfig(group, dataId string, listen ...bool) *Config {
	return &Config{
		Group:  group,
		DataId: dataId,
		Type:   vo.ConfigType(filex.Suffix(dataId)),
		Listen: anyx.Default(false, listen...),
	}
}

// 配置信息格式化
func (c *Config) Info() string {
	return fmt.Sprintf("group=%s dataId=%s", c.Group, c.DataId)
}

// 转化配置项
func (c *Config) ToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{Group: c.Group, DataId: c.DataId, Type: c.Type}
}

// 加载nacos配置
func (c *Config) Loading(v any) (err error) {
	valueRef := reflect.ValueOf(v)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		return errors.New("the input parameter is not a pointer type")
	}
	var param = c.ToConfigParam()
	// 读取Nacos配置
	var info = c.Info()
	var content string
	if content, err = ReadConfigContent(c.Group, c.DataId); err != nil {
		log.Error("Read Nacos Config Content Failed: ", info, err)
		return
	}
	if err = marshalx.NewCase(c.DataId).Unmarshal([]byte(content), v); err != nil {
		log.Error("Loading Nacos Config Failed: ", info, err)
		return
	}
	log.Info("Loading Nacos Config Successful: ", info)
	if c.Listen {
		// 设置Nacos配置监听
		GetNacosConfigMonitor().Set(c.Group, c.DataId, content)
		// 配置监听响应方法
		param.OnChange = func(namespace, group, dataId, data string) {
			log.Errorf("The config on nacos has changed!!!\n dataId=%s group=%s namespace=%s\nThe latest config content is :\n%s", dataId, group, namespace, data)
			GetNacosConfigMonitor().Set(group, dataId, data)
		}
		if err = This().ConfigClient.ListenConfig(param); err != nil {
			log.Error("Listen Nacos Config Failed: ", info, err)
			return
		}
		log.Info("Listen Nacos Config Successful: ", info)
	}
	return
}
