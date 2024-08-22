package respx

// Enum 错误枚举
type Enum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	SuccessCode      = &Enum{Code: 10000, Msg: "Success"}
	ErrorCode        = &Enum{Code: 10001, Msg: "error"}
	AuthFailedCode   = &Enum{Code: 10401, Msg: "authentication failed"}
	ParamErrorCode   = &Enum{Code: 10501, Msg: "request parameter error"}
	RequiredCode     = &Enum{Code: 10502, Msg: "request parameter required"}
	UploadFailedCode = &Enum{Code: 10601, Msg: "upload failed"}
	ImportFailedCode = &Enum{Code: 10601, Msg: "import failed"}
	ExportFailedCode = &Enum{Code: 10602, Msg: "export failed"}
)

func (e *Enum) Response(data any) *Response {
	return &Response{Code: e.Code, Msg: e.Msg, Data: data}
}
