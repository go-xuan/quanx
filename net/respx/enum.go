package respx

// 错误枚举
type Enum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	Success     = Enum{Code: 10000, Msg: "success"}
	Error       = Enum{Code: 10001, Msg: "error"}
	AuthErr     = Enum{Code: 10401, Msg: "authentication failed"}
	ParamErr    = Enum{Code: 10501, Msg: "request parameter error"}
	RequiredErr = Enum{Code: 10502, Msg: "request parameter required"}
	UploadErr   = Enum{Code: 10601, Msg: "upload failed"}
	ImportErr   = Enum{Code: 10601, Msg: "import failed"}
	ExportErr   = Enum{Code: 10602, Msg: "export failed"}
)
