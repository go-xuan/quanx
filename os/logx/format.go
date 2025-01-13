package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/os/fmtx"
)

func DefaultFormatter() log.Formatter {
	hostname, _ := os.Hostname()
	return &defaultFormatter{
		timeFormat: TimeFormat,    // 默认2006-01-02 15:04:05.999
		hostname:   hostname,      // 默认当前hostname
		Output:     ConsoleOutput, // 默认控制台输出
		useColor:   true,          // 默认使用颜色
	}
}

type defaultFormatter struct {
	timeFormat string
	hostname   string
	Output     string
	useColor   bool
}

// Format 日志格式化,用以实现logrus.Formatter接口
func (f *defaultFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%-23s][%-7s][%s]", time.Now().Format(f.timeFormat), entry.Level.String(), f.hostname))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	if (f.Output == ConsoleOutput || f.Output == DefaultOutput) && f.useColor {
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

type jsonFormatter struct {
	timeFormat string
	hostname   string
}

// Format 日志格式化,用以实现logrus.Formatter接口
func (f *jsonFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b = bytes.Buffer{}
	var logRow = struct {
		Time    string      `json:"time"`
		Level   string      `json:"level"`
		Host    string      `json:"hostname"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Time:    entry.Time.Format(f.timeFormat),
		Level:   entry.Level.String(),
		Host:    f.hostname,
		Message: entry.Message,
		Data:    entry.Data,
	}

	if marshal, err := json.Marshal(logRow); err != nil {
		return nil, err
	} else {
		b.Write(marshal)
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
