package respx

// 错误枚举
type ErrEnum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success     = ErrEnum{Code: 2000, Msg: "success"}
	Error       = ErrEnum{Code: 4444, Msg: "error"}
	ParamErr    = ErrEnum{Code: 4001, Msg: "param error"}
	CookieErr   = ErrEnum{Code: 4002, Msg: "cookie error"}
	AuthErr     = ErrEnum{Code: 4002, Msg: "auth error"}
	PasswordErr = ErrEnum{Code: 4003, Msg: "Password error"}
	RedisGetNil = ErrEnum{Code: 4201, Msg: "redis nil"}
	ExportErr   = ErrEnum{Code: 4002, Msg: "export error"}
)
