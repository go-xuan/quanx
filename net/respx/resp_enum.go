package respx

// 错误枚举
type Enum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success         = Enum{Code: 10000, Msg: "success"}
	Error           = Enum{Code: 10001, Msg: "error"}
	AuthRequiredErr = Enum{Code: 10401, Msg: "token/cookie is required"}
	AuthExpiredErr  = Enum{Code: 10402, Msg: "token/cookie is expired"}
	AuthInvalidErr  = Enum{Code: 10403, Msg: "token/cookie is invalid"}
	ParamErr        = Enum{Code: 10411, Msg: "请求参数错误"}
	RequiredErr     = Enum{Code: 10412, Msg: "请求参数必填"}
	AuthErr         = Enum{Code: 10403, Msg: "鉴权失败"}
	ImportErr       = Enum{Code: 10601, Msg: "导入失败"}
	ExportErr       = Enum{Code: 10602, Msg: "导出失败"}
)
