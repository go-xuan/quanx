package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/osx"
)

func DefaultJsonFormatter() log.Formatter {
	return &jsonFormatter{TimeFormat, osx.Hostname()}
}

type jsonFormatter struct {
	timeFormat string
	hostname   string
}

// Format 日志格式化,用以实现logrus.Formatter接口
func (f *jsonFormatter) Format(entry *log.Entry) ([]byte, error) {
	var buffer = bytes.Buffer{}
	if marshal, err := json.Marshal(LogRecord{
		Time:     entry.Time.Format(f.timeFormat),
		Level:    LevelString(entry.Level, 5),
		Hostname: f.hostname,
		Message:  entry.Message,
		Data:     entry.Data,
	}); err != nil {
		return nil, err
	} else {
		buffer.Write(marshal)
	}
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}

func DefaultFormatter() log.Formatter {
	return &textFormatter{TimeFormat, osx.Hostname(), true}
}

type textFormatter struct {
	timeFormat string
	hostname   string
	color      bool
}

// Format 日志格式化，用以实现logrus.Formatter接口
func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%-23s][%s][%-5s]", time.Now().Format(f.timeFormat), f.hostname, LevelString(entry.Level, 5)))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	if f.color {
		return Color(entry.Level, b.Bytes()), nil
	}
	return b.Bytes(), nil
}

func Color(level log.Level, s []byte) []byte {
	switch level {
	case log.InfoLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, string(s)))
	case log.WarnLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, string(s)))
	case log.ErrorLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 31, string(s)))
	default:
		return s
	}
}

func LevelString(level log.Level, length ...int) string {
	var str string
	switch level {
	case log.TraceLevel:
		str = "trace"
	case log.DebugLevel:
		str = "debug"
	case log.InfoLevel:
		str = "info"
	case log.WarnLevel:
		str = "warn"
	case log.ErrorLevel:
		str = "error"
	case log.FatalLevel:
		str = "fatal"
	case log.PanicLevel:
		str = "panic"
	default:
		str = "unknown"
	}
	if len(length) > 0 && length[0] > 0 && length[0] < len(str) {
		return str[:length[0]]
	}
	return str
}

type LogRecord struct {
	Time     string      `json:"create_time" bson:"create_time"`
	Level    string      `json:"level" bson:"level"`
	Hostname string      `json:"hostname" bson:"hostname"`
	Message  string      `json:"message" bson:"message"`
	Data     interface{} `json:"data" bson:"data"`
}
