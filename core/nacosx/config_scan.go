package nacosx

import (
	"fmt"
	"reflect"

	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

type Scanners []*Scanner

// Loading 批量加载nacos配置
func (list Scanners) Loading(conf any) error {
	for _, scanner := range list {
		if err := scanner.Scan(conf); err != nil {
			return errorx.Wrap(err, "scan nacos config failed")
		}
	}
	return nil
}

// Get 根据配置名获取配置
func (list Scanners) Get(dataId string) *Scanner {
	for _, scanner := range list {
		if scanner.DataId == dataId {
			return scanner
		}
	}
	return nil
}

// Scanner 配置扫描器
type Scanner struct {
	Type   vo.ConfigType `yaml:"type"`   // 配置类型
	Group  string        `yaml:"group"`  // 配置分组
	DataId string        `yaml:"dataId"` // 配置文件ID
	Listen bool          `yaml:"listen"` // 是否启用监听
}

// NewScanner 初始化扫描器
func NewScanner(group, dataId string, listen ...bool) *Scanner {
	return &Scanner{
		Group:  group,
		DataId: dataId,
		Type:   vo.ConfigType(filex.GetSuffix(dataId)),
		Listen: anyx.Default(false, listen...),
	}
}

func (s *Scanner) Info() string {
	return fmt.Sprintf("group=%s dataId=%s", s.Group, s.DataId)
}

// ToConfigParam 转化配置项
func (s *Scanner) ToConfigParam() vo.ConfigParam {
	return vo.ConfigParam{Group: s.Group, DataId: s.DataId, Type: s.Type}
}

// Scan nacos配置扫描
func (s *Scanner) Scan(conf any) error {
	valueRef := reflect.ValueOf(conf)
	// 修改值必须是指针类型否则不可行
	if valueRef.Type().Kind() != reflect.Ptr {
		return errorx.New("the input parameter is not a pointer type")
	}
	var param = s.ToConfigParam()
	// 读取Nacos配置
	var content string
	var err error
	if content, err = ReadConfigContent(s.Group, s.DataId); err != nil {
		log.Error("read nacos config content failed: ", s.Info(), err)
		return errorx.Wrap(err, "read nacos config content failed")
	}
	if err = marshalx.NewCase(s.DataId).Unmarshal([]byte(content), conf); err != nil {
		log.Error("loading nacos config failed: ", s.Info(), err)
		return errorx.Wrap(err, "load nacos config failed")
	}
	log.Info("loading nacos config successfully: ", s.Info())
	if s.Listen {
		// 设置Nacos配置监听
		GetConfigMonitor().Set(s.Group, s.DataId, content)
		// 配置监听响应方法
		param.OnChange = func(namespace, group, dataId, data string) {
			log.WithField("dataId", dataId).
				WithField("group", group).
				WithField("namespace", namespace).
				WithField("content", content).
				Error("The nacos config content has changed!!!")
			GetConfigMonitor().Set(group, dataId, data)
		}
		if err = This().ConfigClient.ListenConfig(param); err != nil {
			log.Error("listen nacos config failed: ", s.Info(), err)
			return errorx.Wrap(err, "listen nacos config failed")
		}
		log.Info("listen nacos config successfully !")
	}
	return nil
}
