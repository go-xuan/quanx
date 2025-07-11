package nacosx

import (
	"fmt"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/stringx"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Reader nacos配置读取器
type Reader struct {
	Type   string // 配置文件类型
	Group  string // 配置所在分组
	DataId string // 配置文件ID
	Data   []byte // 配置文件内容
	Listen bool   // 是否启用监听
}

func (r *Reader) ConfigParam() vo.ConfigParam {
	return vo.ConfigParam{
		DataId: r.DataId,
		Group:  r.Group,
		Type:   vo.ConfigType(stringx.IfZero(r.Type, filex.GetSuffix(r.DataId))),
	}
}

func (r *Reader) Anchor(anchor string) {
	r.Group = anchor
}

func (r *Reader) Location() string {
	return fmt.Sprintf("nacos@%s@%s", r.Group, r.DataId)
}

func (r *Reader) Read(config any) error {
	if r.Data == nil {
		if !Initialized() {
			return errorx.New("nacos not initialized")
		}
		var param = r.ConfigParam()
		if data, err := this().ReadConfig(config, param); err != nil {
			return errorx.Wrap(err, "read config from nacos error")
		} else {
			r.Data = data
		}
		if r.Listen {
			if err := this().ListenConfig(config, param); err != nil {
				return errorx.Wrap(err, "listen nacos config error")
			}
		}
	}
	return nil
}
