package respx

// 错误枚举
type ErrEnum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

// 常用异常
var (
	// 应用错误
	Success     = ErrEnum{Code: 2000, Msg: "服务响应成功"}
	Error       = ErrEnum{Code: 4444, Msg: "服务响应异常"}
	ParamErr    = ErrEnum{Code: 4001, Msg: "请求参数错误"}
	AuthErr     = ErrEnum{Code: 4002, Msg: "鉴权失败"}
	PasswordErr = ErrEnum{Code: 4003, Msg: "密码错误"}
	UniqueErr   = ErrEnum{Code: 4004, Msg: "数据重复"}
	DBSelectErr = ErrEnum{Code: 4101, Msg: "数据库查询错误"}
	DBInsertErr = ErrEnum{Code: 4102, Msg: "数据库写入失败"}
	DBUpdateErr = ErrEnum{Code: 4103, Msg: "数据库更新失败"}
	DBDeleteErr = ErrEnum{Code: 4104, Msg: "数据库删除失败"}
	RedisGetNil = ErrEnum{Code: 4201, Msg: "redis: nil"}
)
