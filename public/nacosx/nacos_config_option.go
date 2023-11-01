package nacosx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/magiconair/properties"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// 配置项
type Options []Option
type Option struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	IsSelf bool          `yaml:"isSelf"` // 私有配置
	Listen bool          `yaml:"listen"` // 是否启用监听
}

// 加载nacos配置
func LoadNacosConfig(options Options, config interface{}) {
	if len(options) == 0 {
		return
	}
	if CTL.NamingClient == nil {
		log.Error("未初始化nacos配置中心客户端!")
		return
	}
	if options != nil && len(options) > 0 {
		var server Option  // 应用服务配置
		var others Options // 其他配置
		for _, option := range options {
			if option.IsServer() {
				server = option
			} else if !option.IsSelf {
				others = append(others, option)
			}
		}
		// 加载其他配置
		err := others.LoadNacosConfig(config)
		if err != nil {
			log.Error("加载nacos其他配置失败！", err)
		}
		// 加载主要配置
		defer func() {
			err = server.LoadNacosConfig(config)
			if err != nil {
				log.Error("加载nacos应用配置失败!", err)
			}
		}()
	}
}

// 配置信息格式化
func (opt Option) Format() string {
	return fmt.Sprintf("group=%s dataId=%s type=%s", opt.Group, opt.DataId, opt.ConfigType())
}

func (opt Option) IsServer() bool {
	var dataId = opt.DataId
	if strings.Contains(dataId, ".") {
		i := strings.LastIndex(dataId, ".")
		dataId = dataId[:i]
	}
	return dataId == "server"
}

// 转化配置项
func (opt Option) TransToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{DataId: opt.DataId, Group: opt.Group, Type: opt.ConfigType()}
}

// 获取配置文件类型
func (opt Option) ConfigType() (confType vo.ConfigType) {
	if opt.Type == "" {
		for i := len(opt.DataId) - 1; i >= 0; i-- {
			if opt.DataId[i] == '.' {
				confType = vo.ConfigType(opt.DataId[i+1:])
				return
			}
		}
	} else {
		confType = opt.Type
	}
	return
}

// 加载nacos配置
func (opt Option) LoadNacosConfig(config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		log.Error("加载nacos配置异常，入参conf不是指针类型")
		return errors.New("入参conf不是指针类型")
	}
	// 读取Nacos配置
	var content string
	var param = opt.TransToConfigParam()
	content, err = CTL.ConfigClient.GetConfig(param)
	if err != nil {
		log.Error("获取nacos配置内容失败 ", err)
		return
	}
	switch param.Type {
	case vo.YAML:
		err = yaml.Unmarshal([]byte(content), config)
		break
	case vo.JSON:
		err = json.Unmarshal([]byte(content), &config)
		break
	case vo.PROPERTIES:
		var p *properties.Properties
		p, err = properties.LoadString(content)
		if err != nil {
			break
		}
		refreshValueByTag(p, valueRef)
	default:
		err = errors.New("配置项option.type不符合规范")
	}
	msg := opt.Format()
	if err != nil {
		log.Error("加载nacos配置-失败! ", msg, " error : ", err)
		return
	}
	log.Info("加载nacos配置-成功! ", msg)
	// 设置Nacos配置监听
	if opt.Listen {
		// 初始化nacos配置监控
		GetNacosConfigMonitor().AddConfigData(opt.Group, opt.DataId, content)
		err = ListenConfigChange(param)
		if err != nil {
			log.Error("监听nacos配置-失败! ", msg, " error : ", err)
			return
		}
		log.Info("监听nacos配置-成功! ", msg)
	}
	return
}

func (opt Option) GetNacosConfig() (conf map[string]interface{}, err error) {
	// 读取Nacos配置
	var content string
	var param = opt.TransToConfigParam()
	content, err = CTL.ConfigClient.GetConfig(param)
	if err != nil {
		log.Error("获取nacos配置内容失败 ", err)
		return
	}
	conf = make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &conf)
	if err != nil {
		return
	}
	return
}

// 批量加载nacos配置
func (list Options) LoadNacosConfig(config interface{}) (err error) {
	for _, item := range list {
		err = item.LoadNacosConfig(config)
		if err != nil {
			return
		}
	}
	return
}

// 根据配置名获取配置
func (list Options) Get(dataId string) (target Option) {
	for _, item := range list {
		if item.DataId == dataId {
			target = item
			return
		}
	}
	return
}

// 将配置项进行单项分离
func (list Options) Separate() (Options, Options) {
	var needs Options
	var notNeeds Options
	for _, item := range list {
		if item.IsSelf {
			needs = append(needs, item)
		} else {
			notNeeds = append(notNeeds, item)
		}
	}
	return needs, notNeeds
}

// 监听nacos配置改动
func ListenConfigChange(param vo.ConfigParam) error {
	param.OnChange = func(namespace, group, dataId, data string) {
		log.Errorf("nacos配置已改动 \nnamespace :%s\nGroup     :%s\nData Id   :%s\n%s", namespace, group, dataId, data)
		GetNacosConfigMonitor().UpdateConfigData(group, dataId, data)
	}
	return CTL.ConfigClient.ListenConfig(param)
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
