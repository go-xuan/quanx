package logx

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-xuan/utilx/osx"
	"github.com/go-xuan/utilx/stringx"
	log "github.com/sirupsen/logrus"
)

func NewFormatter(formatter, timeFormat string) log.Formatter {
	timeFormat = stringx.IfZero(timeFormat, TimeFormat)
	switch formatter {
	case FormatterJson:
		return &jsonFormatter{
			timeFormat: timeFormat,
			hostname:   osx.Hostname(),
		}
	case FormatterText:
		return &textFormatter{
			timeFormat: timeFormat,
			hostname:   osx.Hostname(),
		}
	default:
		return DefaultFormatter()
	}
}

func DefaultFormatter() log.Formatter {
	return &textFormatter{
		timeFormat: TimeFormat,
		hostname:   osx.Hostname(),
	}
}

type jsonFormatter struct {
	timeFormat string
	hostname   string
}

// Format 日志格式化
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

type textFormatter struct {
	timeFormat string
	hostname   string
	color      bool
}

func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var buffer = bytes.Buffer{}
	timeStr, level := entry.Time.Format(f.timeFormat), LevelString(entry.Level, 5)
	buffer.WriteString(fmt.Sprintf("[%-23s][%s][%-5s]", timeStr, f.hostname, level))
	buffer.WriteString(entry.Message)
	for key, value := range entry.Data {
		buffer.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	buffer.WriteString("\n")
	if f.color {
		return Color(entry.Level, buffer.Bytes()), nil
	}
	return buffer.Bytes(), nil
}

func Color(level log.Level, data []byte) []byte {
	switch level {
	case log.InfoLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 32, string(data)))
	case log.WarnLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 33, string(data)))
	case log.ErrorLevel:
		return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m", 31, string(data)))
	default:
		return data
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
