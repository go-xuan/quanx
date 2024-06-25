package errorx

/*
错误码： 10-0000-00-0001
错误码格式：[功能码#2][服务码#3][级别码#2][自定义码#1]
功能码：标识所属功能，枚举类型，
级别码：对应错误级别
解释说明：错误码由4部分构成：
*/

const (
	LevelCodeNormal = "00" // 正常
	LevelCodeIgnore = "01" // 忽略
	LevelCodeError  = "02" // 报错
	LevelCodePanic  = "03" // 恐慌
)

const (
	FuncCodeNormal    = "00" // 正常
	FuncCodeSystem    = "01" // 系统
	FuncCodeRequest   = "02" // 服务请求
	FuncCodeConfig    = "03" // 配置加载
	FuncCodeInit      = "04" // 服务初始化
	FuncCodeSerialize = "05" // 序列化
	FuncCodeHttp      = "06" // http请求
	FuncCodeIo        = "07" // IO操作
	FuncCodeCache     = "08" // 缓存
	FuncCodeDB        = "09" // 数据库
	FuncCodeQuery     = "11" // 查询
	FuncCodeCreate    = "12" // 新建
	FuncCodeUpdate    = "13" // 更新
	FuncCodeDelete    = "14" // 删除
	FuncCodeUpload    = "15" // 上传
	FuncCodeExport    = "16" // 导出
)
