package respx

// 错误枚举
type ErrorEnum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success     = ErrorEnum{Code: 10200, Msg: "success"}
	Error       = ErrorEnum{Code: 10444, Msg: "error"}
	ParamErr    = ErrorEnum{Code: 10401, Msg: "请求参数错误"}
	AuthErr     = ErrorEnum{Code: 10403, Msg: "鉴权失败"}
	PasswordErr = ErrorEnum{Code: 10402, Msg: "密码错误"}
	RedisNilErr = ErrorEnum{Code: 10501, Msg: "缓存为空"}
	ExportErr   = ErrorEnum{Code: 10601, Msg: "导出失败"}
)
