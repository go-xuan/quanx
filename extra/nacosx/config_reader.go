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

func (r *Reader) ConfigParam() vo.ConfigParam {
	return vo.ConfigParam{
		DataId: r.DataId,
		Group:  r.Group,
		Type:   vo.ConfigType(r.Type),
	}
}

// ReadConfig nacos配置读取
func (r *Reader) ReadConfig(config any) error {
	// 配置值必须是指针类型，否则不允许读取
	if ref := reflect.ValueOf(config); ref.Type().Kind() != reflect.Ptr {
		return errorx.New("the scanned object must be of pointer type")
	}
	if r.Type == "" {
		r.Type = filex.GetSuffix(r.DataId)
	}
	// 读取Nacos配置文本
	content, err := GetNacosConfigClient().GetConfig(r.ConfigParam())
	if err != nil {
		log.Error("read nacos config content failed: ", r.Info(), err)
		return errorx.Wrap(err, "read config from nacos failed")
	}
	// 配置文本反序列化
	if err = marshalx.Apply(r.DataId).Unmarshal([]byte(content), config); err != nil {
		log.Error("unmarshal nacos config failed: ", r.Info(), err)
		return errorx.Wrap(err, "unmarshal config from nacos failed")
	} else {
		log.Info("unmarshal nacos config success: ", r.Info())
	}
	if err = r.ListenConfig(config); err != nil {
		return errorx.Wrap(err, "listen config failed")
	}
	return nil
}

// ListenConfig 监听nacos配置
func (r *Reader) ListenConfig(config any) error {
	if r.Listen {
		var param = r.ConfigParam()
		// 配置监听响应方法
		param.OnChange = func(namespace, group, dataId, data string) {
			log.WithField("dataId", dataId).
				WithField("group", group).
				WithField("namespace", namespace).
				WithField("data", data).
				Info("the nacos config content has changed !!!")
			if err := marshalx.Apply(dataId).Unmarshal([]byte(data), config); err != nil {
				log.Errorf("update config error, group: %s; dataId: %s; data: %s", group, dataId, data)
			}
		}
		if err := GetNacosConfigClient().ListenConfig(param); err != nil {
			log.Error("listen nacos config failed: ", r.Info(), err)
			return errorx.Wrap(err, "listen nacos config failed")
		}
		log.Info("listen nacos config success!")
	}
	return nil
}
