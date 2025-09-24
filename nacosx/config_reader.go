package nacosx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewReader 创建nacos配置读取器
func NewReader(dataId string) *Reader {
	return &Reader{DataId: dataId}
}

// Reader nacos配置读取器
type Reader struct {
	DataId string `json:"dataId"` // 配置文件ID
	Group  string `json:"group"`  // 配置所在分组
	Type   string `json:"type"`   // 配置文件类型
	Data   []byte `json:"data"`   // 配置文件内容
	Listen bool   `json:"listen"` // 是否启用监听
}

func (r *Reader) ConfigParam() vo.ConfigParam {
	return vo.ConfigParam{
		DataId:  r.DataId,
		Group:   r.Group,
		Content: string(r.Data),
		Type:    vo.ConfigType(r.GetType()),
	}
}

func (r *Reader) GetType() string {
	if r.Type == "" {
		r.Type = filex.GetSuffix(r.DataId)
	}
	return r.Type
}

func (r *Reader) Anchor(group string) {
	if r.Group == "" {
		r.Group = group
	}
}

// Read 从nacos中读取配置
func (r *Reader) Read(v any) error {
	if r.Data == nil {
		if !Initialized() {
			return errorx.New("nacos not initialized")
		}

		// 配置文件锚点为Group分组
		r.Anchor(this().GetGroup())

		param := r.ConfigParam()
		if data, err := this().ReadConfig(v, param); err != nil {
			return errorx.Wrap(err, "read nacos config error")
		} else {
			r.Data = data
		}

		if r.Listen {
			// 监听配置变化
			if err := this().ListenConfig(v, param); err != nil {
				return errorx.Wrap(err, "listen nacos config error")
			}
		}
	}
	return nil
}

// Location 配置文件位置
func (r *Reader) Location() string {
	return fmt.Sprintf("nacos@%s@%s", r.Group, r.DataId)
}

// Write 将配置写入nacos
func (r *Reader) Write(v any) error {
	if !Initialized() {
		return errorx.New("nacos not initialized")
	}

	// 序列化配置
	if data, err := marshalx.Apply(r.GetType()).Marshal(v); err != nil {
		return errorx.Wrap(err, "marshal config error")
	} else {
		r.Data = data
	}

	// 配置文件锚点为Group分组
	r.Anchor(this().GetGroup())

	// 发布配置
	param := r.ConfigParam()
	if err := this().PublishConfig(param); err != nil {
		return errorx.Wrap(err, "publish config to nacos error")
	}
	return nil
}
