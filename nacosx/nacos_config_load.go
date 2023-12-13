package nacosx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xuan/quanx/utilx/structx"
	"reflect"
	"regexp"
	"strings"

	"github.com/magiconair/properties"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

// 配置项
type Configs []*Config
type Config struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	Listen bool          `yaml:"listen"` // 是否启用监听
}

func NewConfig(group string, dataId string, listen bool) *Config {
	return &Config{Group: group, DataId: dataId, Listen: listen}
}

// 配置信息格式化
func (c *Config) ToString() string {
	return fmt.Sprintf("group=%s dataId=%s", c.Group, c.DataId)
}

func (c *Config) IsServer() bool {
	var dataId = c.DataId
	if strings.Contains(dataId, ".") {
		i := strings.LastIndex(dataId, ".")
		dataId = dataId[:i]
	}
	return dataId == "server"
}

func (c *Config) Exist() bool {
	if !Initialized() {
		return false
	}
	content, err := handler.ConfigClient.GetConfig(c.ToConfigParam())
	if err == nil && content != "" {
		return true
	}
	return false
}

// 转化配置项
func (c *Config) ToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{DataId: c.DataId, Group: c.Group, Type: c.ConfigType()}
}

// 获取配置文件类型
func (c *Config) ConfigType() (confType vo.ConfigType) {
	if c.Type == "" {
		for i := len(c.DataId) - 1; i >= 0; i-- {
			if c.DataId[i] == '.' {
				confType = vo.ConfigType(c.DataId[i+1:])
				return
			}
		}
	} else {
		confType = c.Type
	}
	return
}

// 加载nacos配置
func (c *Config) LoadConfig(config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		log.Error("加载nacos配置异常，入参conf不是指针类型")
		return errors.New("入参conf不是指针类型")
	}
	// 读取Nacos配置
	var content string
	var param = c.ToConfigParam()
	content, err = handler.ConfigClient.GetConfig(param)
	if err != nil {
		log.Error("获取nacos配置内容失败 ", err)
		return
	}
	var msg = c.ToString()
	if err = structx.ParseBytesToPointer(config, []byte(content), param.DataId); err != nil {
		log.Error("nacos配置加载失败! ", msg, " error : ", err)
		return
	}
	log.Info("nacos配置加载成功! ", msg)
	// 设置Nacos配置监听
	if c.Listen {
		// 初始化nacos配置监控
		GetNacosConfigMonitor().AddConfigData(c.Group, c.DataId, content)
		if err = ListenConfigChange(param); err != nil {
			log.Error("监听nacos配置-失败! ", msg, " error : ", err)
			return
		}
		log.Info("监听nacos配置-成功! ", msg)
	}
	return
}

// 获取配置项键值对
func (c *Config) GetKeyValueMap() (kvm map[string]interface{}, err error) {
	// 读取Nacos配置
	var content string
	content, err = handler.ConfigClient.GetConfig(c.ToConfigParam())
	if err != nil {
		log.Error("获取nacos配置内容失败 ", err)
		return
	}
	kvm = make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &kvm)
	if err != nil {
		return
	}
	return
}

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

// 监听nacos配置改动
func ListenConfigChange(param vo.ConfigParam) error {
	param.OnChange = func(namespace, group, dataId, data string) {
		log.Errorf("监听到nacos配置已改动!!! \n dataId=%s group=%s namespace=%s\n改动后内容如下:\n%s", dataId, group, namespace, data)
		GetNacosConfigMonitor().UpdateConfigData(group, dataId, data)
	}
	return handler.ConfigClient.ListenConfig(param)
}

// 通过配置标签刷新配置
func refreshValueByTag(p *properties.Properties, v reflect.Value) {
	m := p.Map()
	for i := 0; i < v.Elem().NumField(); i++ {
		//先判断有没有nacos的配置
		tag := v.Elem().Type().Field(i).Tag.Get("nacos")
		r, _ := regexp.Compile("\\${.*?}")
		gs := r.FindAllString(tag, -1)
		for _, str := range gs {
			if len(str) <= 3 {
				tag = strings.ReplaceAll(tag, str, "")
			} else {
				envStr := str[2 : len(str)-1]
				tag = strings.ReplaceAll(tag, str, strings.Split(v.FieldByName(envStr).String(), ".")[0])
			}
		}
		if tag == "" && reflect.Struct != v.Elem().Field(i).Kind() {
			continue
		}
		switch v.Elem().Field(i).Kind() {
		case reflect.String:
			temp, ok := p.Get(tag)
			if ok {
				v.Elem().Field(i).SetString(temp)
			}
		case reflect.Bool:
			_, ok := m[tag]
			if ok {
				temp := p.GetBool(tag, false)
				v.Elem().Field(i).SetBool(temp)
			}
		case reflect.Int:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Int8:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Int16:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Int32:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Int64:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Uint8:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Uint16:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Uint32:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Uint64:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Elem().Field(i).SetInt(int64(temp))
			}
		case reflect.Float32:
			_, ok := m[tag]
			if ok {
				temp := p.GetFloat64(tag, 0)
				v.Elem().Field(i).SetFloat(temp)
			}
		case reflect.Float64:
			_, ok := m[tag]
			if ok {
				temp := p.GetFloat64(tag, 0)
				v.Elem().Field(i).SetFloat(temp)
			}
		case reflect.Struct:
			refreshStructValueByTag(p, v.Elem().Field(i))
		default:
			fmt.Printf("未匹配到type %s", v.Elem().Field(i).Kind())
		}
	}
}

// 通过配置标签刷新结构体配置
func refreshStructValueByTag(p *properties.Properties, v reflect.Value) {
	m := p.Map()
	for i := 0; i < v.NumField(); i++ {
		//先判断有没有nacos的配置
		tag := v.Type().Field(i).Tag.Get("nacos")
		r, _ := regexp.Compile("\\${.*?}")
		gs := r.FindAllString(tag, -1)
		for _, str := range gs {
			if len(str) <= 3 {
				tag = strings.ReplaceAll(tag, str, "")
			} else {
				envStr := str[2 : len(str)-1]
				tag = strings.ReplaceAll(tag, str, strings.Split(v.FieldByName(envStr).String(), ".")[0])
			}
		}

		if tag == "" && reflect.Struct != v.Field(i).Kind() {
			continue
		}

		switch v.Field(i).Kind() {
		case reflect.String:
			temp, ok := p.Get(tag)
			if ok {
				v.Field(i).SetString(temp)
			}
		case reflect.Bool:
			_, ok := m[tag]
			if ok {
				temp := p.GetBool(tag, false)
				v.Field(i).SetBool(temp)
			}
		case reflect.Int:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Int8:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Int16:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Int32:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Int64:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Uint8:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Uint16:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Uint32:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Uint64:
			_, ok := m[tag]
			if ok {
				temp := p.GetInt(tag, 0)
				v.Field(i).SetInt(int64(temp))
			}
		case reflect.Float32:
			_, ok := m[tag]
			if ok {
				temp := p.GetFloat64(tag, 0)
				v.Field(i).SetFloat(temp)
			}
		case reflect.Float64:
			_, ok := m[tag]
			if ok {
				temp := p.GetFloat64(tag, 0)
				v.Field(i).SetFloat(temp)
			}
		case reflect.Struct:
			refreshStructValueByTag(p, v.Field(i))
		default:
			fmt.Printf("未匹配到type %s", v.Field(i).Kind())
		}
	}
}
