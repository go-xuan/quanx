package logx

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/stringx"
)

const (
	TimeFormat = "2006-01-02 15:04:05.999"
)

type LogFormatter struct {
	TimeFormat string
}

// 日志格式化
func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timeFormat := stringx.IfZero(f.TimeFormat, TimeFormat)
	host, _ := os.Hostname()
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%23s][%7s][%s]", time.Now().Format(timeFormat), entry.Level.String(), host))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}

// gin框架生成日志Handler
func GinRequestLog(context *gin.Context) {
	start := time.Now()
	// 处理请求
	context.Next()
	// 日志格式
	log.Infof("[%3d][%8dms][%15s][%6s][%s]",
		context.Writer.Status(),
		time.Now().Sub(start).Milliseconds(),
		stringx.IfNot(context.ClientIP(), "::1", "localhost"),
		context.Request.Method,
		context.Request.RequestURI,
	)
}
