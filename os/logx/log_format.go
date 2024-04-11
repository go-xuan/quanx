package logx

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/utils/anyx"
)

const (
	TimeFormat = "2006-01-02 15:04:05.999999"
)

type LogFormatter struct {
	TimeFormat string
}

// 日志格式化
func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timeFormat := anyx.IfZero(f.TimeFormat, TimeFormat)
	msg := entry.Message
	for key, value := range entry.Data {
		msg += fmt.Sprintf(" , %s : %+v", key, value)
	}
	host, _ := os.Hostname()
	sf := fmt.Sprintf("[%26s][%7s][%s]%15s\n",
		time.Now().Format(timeFormat),
		entry.Level.String(),
		host,
		msg)
	return []byte(sf), nil
}

// gin框架生成日志Handler
func LoggerToFile() gin.HandlerFunc {
	return func(context *gin.Context) {
		// 处理请求
		context.Next()
		ip := context.ClientIP()
		if ip == "::1" {
			ip = "localhost"
		}
		// 日志格式
		log.Infof("[%3d][%10v][%15s][%4s][%s]",
			context.Writer.Status(),
			time.Now().Sub(time.Now()),
			ip,
			context.Request.Method,
			context.Request.RequestURI,
		)
	}
}
