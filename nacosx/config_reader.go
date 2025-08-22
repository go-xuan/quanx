package nacosx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/go-xuan/utilx/stringx"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Reader nacos配置读取器
type Reader struct {
	DataId string // 配置文件ID
	Group  string // 配置所在分组
	Type   string // 配置文件类型
	Data   []byte // 配置文件内容
	Listen bool   // 是否启用监听
}

// ConfigParam 获取nacos配置参数
func (r *Reader) ConfigParam() vo.ConfigParam {
	return vo.ConfigParam{
		DataId:  r.DataId,
		Group:   r.Group,
		Content: string(r.Data),
		Type:    vo.ConfigType(stringx.IfZero(r.Type, filex.GetSuffix(r.DataId))),
	}
}

// Anchor 配置读取器锚点为配置分组
func (r *Reader) Anchor(anchor string) {
	if r.Group == "" {
		r.Group = anchor
	}
}

// Location 配置文件位置
func (r *Reader) Location() string {
	return fmt.Sprintf("nacos@%s@%s", r.Group, r.DataId)
}

// Read 从nacos中读取配置
func (r *Reader) Read(v any) error {
	if r.Data == nil {
		if !Initialized() {
			return errorx.New("nacos not initialized")
		}
		param := r.ConfigParam()
		data, err := this().ReadConfig(v, param)
		if err != nil {
			return errorx.Wrap(err, "read nacos config error")
		}
		r.Data = data
		if r.Listen {
			if err = this().ListenConfig(v, param); err != nil {
				return errorx.Wrap(err, "listen nacos config error")
			}
		}
	}
	return nil
}

// Write 将配置写入nacos
func (r *Reader) Write(config any) error {
	if !Initialized() {
		return errorx.New("nacos not initialized")
	}
	if r.Type == "" {
		r.Type = filex.GetSuffix(r.DataId)
	}
	type_ := stringx.IfZero(r.Type, filex.GetSuffix(r.DataId))
	data, err := marshalx.Apply(type_).Marshal(config)
	if err != nil {
		return errorx.Wrap(err, "marshal config error")
	}
	r.Data = data
	if err = this().PublishConfig(r.ConfigParam()); err != nil {
		return errorx.Wrap(err, "publish config to nacos error")
	}
	return nil
}
