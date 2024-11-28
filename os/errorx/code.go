package errorx

type Enum struct {
	Code string
	Msg  string
}

// 标签码
var (
	ErrorTagSuccess = Enum{Code: "0", Msg: "成功"} // 正常
	ErrorTagFailed  = Enum{Code: "1", Msg: "失败"}
	ErrorTagConfig  = Enum{Code: "2", Msg: "服务配置错误"}
	ErrorTagInit    = Enum{Code: "3", Msg: "服务初始化异常"}
	ErrorTagRequest = Enum{Code: "5", Msg: "请求异常"}
	ErrorTagHttp    = Enum{Code: "6", Msg: "外部请求异常"}
	ErrorTagLock    = Enum{Code: "7", Msg: "锁异常"}
	ErrorTagIO      = Enum{Code: "8", Msg: "IO异常"}
	ErrorTagCache   = Enum{Code: "9", Msg: "缓存异常"}
	ErrorTagDB      = Enum{Code: "0", Msg: "数据库异常"}
)

// 错误级别码
var (
	LevelNormal = Enum{Code: "0", Msg: "正常"}
	LevelIgnore = Enum{Code: "1", Msg: "忽略"}
	LevelError  = Enum{Code: "2", Msg: "报错"}
	LevelPanic  = Enum{Code: "3", Msg: "恐慌"}
)
