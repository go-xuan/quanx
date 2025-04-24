package randx

import (
	"strconv"
)

const (
	typeInt      = "int"      // 数字
	typeFloat    = "float"    // 浮点数
	typeSequence = "sequence" // 序列
	typeTime     = "time"     // 时间
	typeDate     = "date"     // 日期
	typeUUID     = "uuid"     // uuid
	typePhone    = "phone"    // 手机号
	typeName     = "name"     // 姓名
	typeIdCard   = "id_card"  // 身份证
	typePlateNo  = "plate_no" // 车牌号
	typeEmail    = "email"    // 邮箱
	typeIP       = "ip"       // ip地址
	typeProvince = "province" // 省
	typeCity     = "city"     // 市
	typePassword = "password" // 密码
	typeEnum     = "enum"     // 枚举
)

// Options 随机生成
type Options struct {
	Param   *Param // 约束条件参数
	Type    string // 数据类型
	Default string // 默认值
	Offset  int    // 偏移量
}

// NewString 生成随机数
func (opt *Options) NewString() string {
	var value string
	if opt.Type == typeSequence {
		value = strconv.Itoa(opt.Param.Sequence(opt.Offset))
	} else {
		value = opt.Param.String(opt.Type)
	}
	return opt.Param.Modify(value)
}
