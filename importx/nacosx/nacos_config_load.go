package nacosx

import (
	"errors"
	"fmt"
	"github.com/go-xuan/quanx/importx/marshalx"
	"reflect"

	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/utilx/filex"
)

// 配置项
type Configs []*Config

// 批量加载nacos配置
func (list Configs) LoadConfig(config interface{}) (err error) {
	for _, conf := range list {
		err = conf.LoadConfig(config)
		if err != nil {
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

func NewConfig(group, dataId string) *Config {
	return &Config{Group: group, DataId: dataId, Type: vo.ConfigType(filex.Suffix(dataId))}
}

// 配置信息格式化
func (c *Config) ToString(title string) string {
	return fmt.Sprintf("%s => group=%s dataId=%s", title, c.Group, c.DataId)
}

// 转化配置项
func (c *Config) ToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{Group: c.Group, DataId: c.DataId, Type: c.Type}
}

// 加载nacos配置
func (c *Config) LoadConfig(config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		log.Error("load nacos config failed!")
		return errors.New("the input parameter is not of pointer type")
	}
	var param = c.ToConfigParam()
	// 读取Nacos配置
	var content string
	content, err = GetConfigContent(c.Group, c.DataId)
	if err != nil {
		log.Error("get config from nacos failed : ", err)
		return
	}
	if err = marshalx.UnmarshalToPointer(config, []byte(content), filex.Suffix(c.DataId)); err != nil {
		log.Error(c.ToString("load nacos config failed!"))
		log.Error(" error : ", err)
		return
	}
	log.Info(c.ToString("load nacos config successful!"))
	// 设置Nacos配置监听
	if c.Listen {
		// 新增nacos配置监听
		GetNacosConfigMonitor().AddConfigData(c.Group, c.DataId, content)
		param.OnChange = ConfigChangedMonitor
		if err = This().ConfigClient.ListenConfig(param); err != nil {
			log.Error(c.ToString("listen nacos config failed!"))
			log.Error(" error : ", err)
			return
		}
		log.Info(c.ToString("listen nacos config successful!"))
	}
	return
}

func ConfigChangedMonitor(namespace, group, dataId, data string) {
	log.Errorf("Listen nacos config changed!!!\n dataId=%s group=%s namespace=%s\n改动后内容如下:\n%s", dataId, group, namespace, data)
	GetNacosConfigMonitor().UpdateConfigData(group, dataId, data)
}
