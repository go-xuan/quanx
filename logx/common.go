package logx

// 日志级别
const (
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelError = "error"
	LevelFatal = "fatal"
	LevelPanic = "panic"

	WriterConsole = "console" // 控制台打印
	WriterFile    = "file"    // 写入日志文件

	FormatterText = "text"                    // 文本格式化
	FormatterJson = "json"                    // json格式化
	TimeFormat    = "2006-01-02 15:04:05.000" // 时间格式化
)