package respx

// Enum 错误枚举
type Enum struct {
	Code int    `json:"code"` // 响应状态码
	Msg  string `json:"msg"`  // 响应消息
}

var (
	SuccessEnum      = &Enum{Code: 10000, Msg: "Success"}
	ErrorEnum        = &Enum{Code: 10001, Msg: "error"}
	AuthFailedEnum   = &Enum{Code: 10401, Msg: "auth failed"}
	ParamErrorEnum   = &Enum{Code: 10501, Msg: "request parameter error"}
	RequiredEnum     = &Enum{Code: 10502, Msg: "request parameter required"}
	UploadFailedEnum = &Enum{Code: 10601, Msg: "upload failed"}
	ImportFailedEnum = &Enum{Code: 10601, Msg: "import failed"}
	ExportFailedEnum = &Enum{Code: 10602, Msg: "export failed"}
)

func (e *Enum) Response(data any) *Response {
	return &Response{Code: e.Code, Msg: e.Msg, Data: data}
}
