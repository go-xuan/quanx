package nacosx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// NewReader 创建nacos配置读取器
func NewReader(dataId string, listen ...bool) *Reader {
	return &Reader{
		DataId: dataId,
		Type:   filex.GetSuffix(dataId),
		Listen: len(listen) > 0 && listen[0],
	}
}

// Reader nacos配置读取器
type Reader struct {
	DataId string // 配置文件id
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
		Type:    vo.ConfigType(r.GetType()),
	}
}

func (r *Reader) GetType() string {
	if r.Type == "" {
		r.Type = marshalx.JSON
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
		client := GetClient()
		// 配置文件锚点为group分组
		r.Anchor(client.GetGroup())

		// 读取配置
		param := r.ConfigParam()
		data, err := client.ReadConfig(v, param)
		if err != nil {
			return errorx.Wrap(err, "read nacos config failed")
		}
		r.Data = data

		if r.Listen {
			// 监听配置变化
			if err = client.ListenConfig(v, param); err != nil {
				return errorx.Wrap(err, "listen nacos config failed")
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
	client := GetClient()
	// 配置文件锚点为group分组
	r.Anchor(client.GetGroup())

	// 序列化配置
	data, err := marshalx.Apply(r.GetType()).Marshal(v)
	if err != nil {
		return errorx.Wrap(err, "marshal config failed")
	}
	r.Data = data

	// 发布配置
	param := r.ConfigParam()
	if err = client.PublishConfig(param); err != nil {
		return errorx.Wrap(err, "publish nacos config failed")
	}
	return nil
}
