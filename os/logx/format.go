package logx

import (
	"bytes"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/fmtx"
)

func DefaultFormatter() log.Formatter {
	host, _ := os.Hostname()
	return &LogFormatter{
		timeFormat: TimeFormat,    // 默认2006-01-02 15:04:05.999
		host:       host,          // 默认当前机器host
		Output:     ConsoleOutput, // 默认控制台输出
		useColor:   true,          // 默认使用颜色
	}
}

type LogFormatter struct {
	timeFormat string
	host       string
	Output     string
	useColor   bool
}

func (f *LogFormatter) UseColor() bool {
	return (f.Output == ConsoleOutput || f.Output == DefaultOutput) && f.useColor
}

// Format 日志格式化,用以实现logrus.Formatter接口
func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%-23s][%-7s][%s]", time.Now().Format(f.timeFormat), entry.Level.String(), f.host))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	if f.UseColor() {
		return Color(entry.Level).Bytes(b.String()), nil
	} else {
		return b.Bytes(), nil
	}
}

func Color(level log.Level) fmtx.Color {
	switch level {
	case log.InfoLevel:
		return fmtx.Green
	case log.WarnLevel:
		return fmtx.Yellow
	case log.ErrorLevel:
		return fmtx.Red
	default:
		return 0
	}
}
