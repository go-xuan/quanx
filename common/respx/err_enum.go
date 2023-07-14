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
	Error       = ErrEnum{Code: 5000, Msg: "服务响应异常"}
	Fail        = ErrEnum{Code: 4000, Msg: "服务请求失败"}
	ParamErr    = ErrEnum{Code: 2001, Msg: "请求参数错误"}
	UniqueErr   = ErrEnum{Code: 2002, Msg: "唯一校验失败"}
	AuthErr     = ErrEnum{Code: 1001, Msg: "鉴权失败"}
	PasswordErr = ErrEnum{Code: 1002, Msg: "密码错误"}
	DBSelectErr = ErrEnum{Code: 6001, Msg: "数据库查询错误"}
	DBInsertErr = ErrEnum{Code: 6002, Msg: "数据库写入失败"}
	DBUpdateErr = ErrEnum{Code: 6003, Msg: "数据库更新失败"}
	DBDeleteErr = ErrEnum{Code: 6004, Msg: "数据库删除失败"}
	RedisGetNil = ErrEnum{Code: 6101, Msg: "redis: nil"}
)
