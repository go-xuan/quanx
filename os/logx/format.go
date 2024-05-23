package logx

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/core/ginx"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/types/stringx"
)

func DefaultFormatter() logrus.Formatter {
	host, _ := os.Hostname()
	return &LogFormatter{timeFormat: TimeFormat, host: host, Output: ConsoleOutput, useColor: true}
}

type LogFormatter struct {
	timeFormat string
	host       string
	Output     string
	useColor   bool
}

// 日志格式化,用以实现logrus.Formatter接口
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%-23s][%-5s][%s] ", time.Now().Format(f.timeFormat), entry.Level.String(), f.host))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	if f.UserColor() {
		return Color(entry.Level).Bytes(b.String()), nil
	} else {
		return b.Bytes(), nil
	}
}

func (f *LogFormatter) UserColor() bool {
	return (f.Output == ConsoleOutput || f.Output == DefaultOutput) && f.useColor
}

func Color(level logrus.Level) fmtx.Color {
	switch level {
	case logrus.InfoLevel:
		return fmtx.Green
	case logrus.WarnLevel:
		return fmtx.Yellow
	case logrus.ErrorLevel:
		return fmtx.Red
	default:
		return 0
	}
}

// gin请求日志中间件
func GinRequestLog(ctx *gin.Context) {
	start := time.Now()
	// 处理请求
	ctx.Next()
	var ip string
	if ipv, ok := ctx.Get(ginx.IPKey); ok {
		ip = ipv.(string)
	} else {
		ip = stringx.IfNot(ctx.ClientIP(), "::1", "localhost")
	}
	// 日志格式
	logrus.Infof("[%3d][%8dms][%15s][%6s][%s]",
		ctx.Writer.Status(),
		time.Now().Sub(start).Milliseconds(),
		ip,
		ctx.Request.Method,
		ctx.Request.RequestURI,
	)
}
