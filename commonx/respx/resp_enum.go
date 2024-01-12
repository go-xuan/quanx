package respx

// 错误枚举
type Enum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success     = Enum{Code: 10200, Msg: "success"}
	Error       = Enum{Code: 10444, Msg: "error"}
	ParamErr    = Enum{Code: 10401, Msg: "请求参数错误"}
	RequiredErr = Enum{Code: 10402, Msg: "请求参数必填"}
	AuthErr     = Enum{Code: 10403, Msg: "鉴权失败"}
	ExportErr   = Enum{Code: 10601, Msg: "导出失败"}
)
