package nacosx

import (
	"fmt"
	"reflect"

	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// Scanner nacos配置扫描器
type Scanner struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	Listen bool          `yaml:"listen"` // 是否启用监听
}

func (s *Scanner) Info() string {
	return fmt.Sprintf("group=%s dataId=%s", s.Group, s.DataId)
}

// Scan nacos配置扫描
func (s *Scanner) Scan(v any) error {
	// 修改值必须是指针类型否则不可行
	if ref := reflect.ValueOf(v); ref.Type().Kind() != reflect.Ptr {
		return errorx.New("the scanned object must be of pointer type")
	}
	var param = vo.ConfigParam{
		DataId: s.DataId,
		Group:  s.Group,
		Type:   s.Type,
	}
	// 读取Nacos配置文本
	content, err := GetNacosConfigClient().GetConfig(param)
	if err != nil {
		log.Error("get nacos config content failed: ", s.Info(), err)
		return errorx.Wrap(err, "get nacos config content failed")
	}
	// 配置文本反序列化
	if err = marshalx.Apply(s.DataId).Unmarshal([]byte(content), v); err != nil {
		log.Error("scan nacos config failed: ", s.Info(), err)
		return errorx.Wrap(err, "scan nacos config failed")
	} else {
		log.Info("scan nacos config success: ", s.Info())
	}
	if s.Listen {
		// 设置Nacos配置监听
		GetConfigMonitor().Set(s.Group, s.DataId, content)
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
			log.Error("listen nacos config failed: ", s.Info(), err)
			return errorx.Wrap(err, "listen nacos config failed")
		} else {
			log.Info("listen nacos config success!")
		}
	}
	return nil
}
