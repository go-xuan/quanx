package nacosx

import (
	"fmt"
	"reflect"
	
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// Reader nacos配置读取器
type Reader struct {
	Type   string `yaml:"type"`   // 配置类型
	Group  string `yaml:"group"`  // 配置分组
	DataId string `yaml:"dataId"` // 配置文件ID
	Data   []byte `yaml:"data"`   // 配置文件内容
	Listen bool   `yaml:"listen"` // 是否启用监听
}

func (r *Reader) Info() string {
	return fmt.Sprintf("group=%s dataId=%s", r.Group, r.DataId)
}

func (r *Reader) ConfigParam() vo.ConfigParam {
	return vo.ConfigParam{
		DataId: r.DataId,
		Group:  r.Group,
		Type:   vo.ConfigType(stringx.IfZero(r.Type, filex.GetSuffix(r.DataId))),
	}
}

func (r *Reader) Ready(group ...string) {
	if len(group) > 0 && group[0] != "" {
		r.Group = group[0]
	}
}

func (r *Reader) Check(config any) error {
	if r.Data == nil {
		if !Initialized() {
			return errorx.New("nacos client not initialized")
		}
		// 配置值必须是指针类型，否则不允许读取
		if ref := reflect.ValueOf(config); ref.Type().Kind() != reflect.Ptr {
			return errorx.New("the scanned object must be of pointer type")
		}
		content, err := GetNacosConfigClient().GetConfig(r.ConfigParam())
		if err != nil {
			return errorx.Wrap(err, "read config from nacos error")
		}
		data := []byte(content)
		if err = marshalx.Apply(r.DataId).Unmarshal(data, config); err != nil {
			return errorx.Wrap(err, "unmarshal config from nacos error")
		}
		r.Data = data
	}
	return nil
}

func (r *Reader) Read(config any) error {
	if err := r.Check(config); err != nil {
		return errorx.Wrap(err, "check nacos config error")
	}
	if err := r.ListenConfig(config); err != nil {
		return errorx.Wrap(err, "listen nacos config error")
	}
	return nil
}

func (r *Reader) Location() string {
	return fmt.Sprintf("nacos@%s@%s", r.Group, r.DataId)
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
