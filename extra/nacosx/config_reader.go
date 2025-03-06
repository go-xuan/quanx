package nacosx

import (
	"fmt"
	"reflect"

	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// Reader nacos配置读取器
type Reader struct {
	Type   string `yaml:"type"`   // 配置类型
	Group  string `yaml:"group"`  // 配置分组
	DataId string `yaml:"dataId"` // 配置文件ID
	Listen bool   `yaml:"listen"` // 是否启用监听
}

func (r *Reader) Info() string {
	return fmt.Sprintf("group=%s dataId=%s", r.Group, r.DataId)
}

func (r *Reader) Location(group ...string) string {
	if len(group) > 0 && group[0] != "" {
		r.Group = group[0]
	}
	return fmt.Sprintf("%s@%s", r.Group, r.DataId)
}

// ReadConfig nacos配置读取
func (r *Reader) ReadConfig(v any) error {
	// 修改值必须是指针类型否则不可行
	if ref := reflect.ValueOf(v); ref.Type().Kind() != reflect.Ptr {
		return errorx.New("the scanned object must be of pointer type")
	}
	if r.Type == "" {
		r.Type = filex.GetSuffix(r.DataId)
	}
	var param = vo.ConfigParam{
		DataId: r.DataId,
		Group:  r.Group,
		Type:   vo.ConfigType(r.Type),
	}
	// 读取Nacos配置文本
	content, err := GetNacosConfigClient().GetConfig(param)
	if err != nil {
		log.Error("get nacos config content failed: ", r.Info(), err)
		return errorx.Wrap(err, "get nacos config content failed")
	}
	// 配置文本反序列化
	if err = marshalx.Apply(r.DataId).Unmarshal([]byte(content), v); err != nil {
		log.Error("scan nacos config failed: ", r.Info(), err)
		return errorx.Wrap(err, "scan nacos config failed")
	} else {
		log.Info("scan nacos config success: ", r.Info())
	}
	if r.Listen {
		// 设置Nacos配置监听
		GetConfigMonitor().Set(r.Group, r.DataId, content)
		// 配置监听响应方法
		param.OnChange = func(namespace, group, dataId, data string) {
			log.WithField("dataId", dataId).
				WithField("group", group).
				WithField("namespace", namespace).
				WithField("content", content).
				Error("the nacos config content has changed !!!")
			GetConfigMonitor().Set(group, dataId, data)
		}
		if err = GetNacosConfigClient().ListenConfig(param); err != nil {
			log.Error("listen nacos config failed: ", r.Info(), err)
			return errorx.Wrap(err, "listen nacos config failed")
		} else {
			log.Info("listen nacos config success!")
		}
	}
	return nil
}
