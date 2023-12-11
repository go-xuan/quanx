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
type Items []Item
type Item struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	Listen bool          `yaml:"listen"` // 是否启用监听
}

// 加载nacos配置
func LoadNacosConfig(group string, config *LoadConfig, target interface{}) {
	if instance.NamingClient == nil {
		log.Error("未初始化nacos配置中心客户端!")
		return
	}
	if config == nil || config.Basic == "" {
		log.Error("当前应用未配置nacos!")
		return
	}
	var listenMap = make(map[string]bool)
	if config.Listen != "" {
		listenIds := strings.Split(config.Listen, ",")
		for _, dataId := range listenIds {
			listenMap[strings.TrimSpace(dataId)] = true
		}
	}
	initIds := strings.Split(config.Basic, ",")
	var final Item
	var pres Items
	// 初始化默认配置项
	for _, dataId := range initIds {
		dataId = strings.TrimSpace(dataId)
		var opt = Item{Group: group, DataId: dataId, Listen: listenMap[dataId]}
		if strings.Contains(dataId, "server") {
			final = opt
		} else {
			pres = append(pres, opt)
		}
	}
	// 预加载配置项
	var err = pres.LoadConfig(target)
	if err != nil {
		log.Error("加载nacos预加载配置项失败！", err)
	}
	// 最终加载配置
	defer func() {
		err = final.LoadConfig(target)
		if err != nil {
			log.Error("加载nacos最终加载配置失败!", err)
		}
	}()
}

// 配置信息格式化
func (item Item) Format() string {
	return fmt.Sprintf("group=%s dataId=%s type=%s", item.Group, item.DataId, item.ConfigType())
}

func (item Item) IsServer() bool {
	var dataId = item.DataId
	if strings.Contains(dataId, ".") {
		i := strings.LastIndex(dataId, ".")
		dataId = dataId[:i]
	}
	return dataId == "server"
}

// 转化配置项
func (item Item) TransToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{DataId: item.DataId, Group: item.Group, Type: item.ConfigType()}
}

// 获取配置文件类型
func (item Item) ConfigType() (confType vo.ConfigType) {
	if item.Type == "" {
		for i := len(item.DataId) - 1; i >= 0; i-- {
			if item.DataId[i] == '.' {
				confType = vo.ConfigType(item.DataId[i+1:])
				return
			}
		}
	} else {
		confType = item.Type
	}
	return
}

// 加载nacos配置
func (item Item) LoadConfig(config interface{}) (err error) {
	valueRef := reflect.ValueOf(config)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		log.Error("加载nacos配置异常，入参conf不是指针类型")
		return errors.New("入参conf不是指针类型")
	}
	// 读取Nacos配置
	var content string
	var param = item.TransToConfigParam()
	content, err = instance.ConfigClient.GetConfig(param)
	if err != nil {
		log.Error("获取nacos配置内容失败 ", err)
		return
	}
	switch param.Type {
	case vo.YAML:
		err = yaml.Unmarshal([]byte(content), config)
	case vo.JSON:
		err = json.Unmarshal([]byte(content), &config)
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
	msg := item.Format()
	if err != nil {
		log.Error("加载nacos配置-失败! ", msg, " error : ", err)
		return
	}
	log.Info("加载nacos配置-成功! ", msg)
	// 设置Nacos配置监听
	if item.Listen {
		// 初始化nacos配置监控
		GetNacosConfigMonitor().AddConfigData(item.Group, item.DataId, content)
		err = ListenConfigChange(param)
		if err != nil {
			log.Error("监听nacos配置-失败! ", msg, " error : ", err)
			return
		}
		log.Info("监听nacos配置-成功! ", msg)
	}
	return
}

func (item Item) GetNacosConfig() (conf map[string]interface{}, err error) {
	// 读取Nacos配置
	var content string
	var param = item.TransToConfigParam()
	content, err = instance.ConfigClient.GetConfig(param)
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
func (items Items) LoadConfig(config interface{}) (err error) {
	for _, item := range items {
		err = item.LoadConfig(config)
		if err != nil {
			return
		}
	}
	return
}

// 根据配置名获取配置
func (items Items) Get(dataId string) (target Item) {
	for _, item := range items {
		if item.DataId == dataId {
			target = item
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
	return instance.ConfigClient.ListenConfig(param)
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
