package respx

// 错误枚举
type ErrorEnum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success     = ErrorEnum{Code: 2000, Msg: "success"}
	Error       = ErrorEnum{Code: 4444, Msg: "error"}
	ParamErr    = ErrorEnum{Code: 4001, Msg: "param error"}
	CookieErr   = ErrorEnum{Code: 4002, Msg: "cookie error"}
	AuthErr     = ErrorEnum{Code: 4002, Msg: "auth error"}
	PasswordErr = ErrorEnum{Code: 4003, Msg: "Password error"}
	RedisNilErr = ErrorEnum{Code: 4201, Msg: "redis nil"}
	ExportErr   = ErrorEnum{Code: 4002, Msg: "export error"}
)
