package logx

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Record 日志记录
type Record struct {
	Time     string      `json:"time" bson:"time"`
	Level    string      `json:"level" bson:"level"`
	Hostname string      `json:"hostname" bson:"hostname"`
	Msg      string      `json:"msg" bson:"msg"`
	Data     interface{} `json:"data" bson:"data"`
}

// Formatter 日志格式化器
type Formatter struct {
	Formatter  string // 日志格式化器，json或text
	TimeLayout string // 时间格式化
	Hostname   string // 主机名
	Color      bool   // 是否使用颜色
}

// Format 日志格式化
func (f *Formatter) Format(entry *log.Entry) ([]byte, error) {
	var data []byte
	var err error
	switch f.Formatter {
	case FormatterJson:
		data, err = f.FormatJson(entry)
	case FormatterText:
		data, err = f.FormatText(entry)
	}
	if err != nil {
		return nil, err
	}
	if f.Color {
		data = Color(entry.Level, data)
	}
	return data, nil
}

// FormatJson 日志json格式化
func (f *Formatter) FormatJson(entry *log.Entry) ([]byte, error) {
	time, level := entry.Time.Format(f.TimeLayout), LevelString(entry.Level, 5)
	data, err := json.Marshal(Record{
		Time:     time,
		Level:    level,
		Hostname: f.Hostname,
		Msg:      entry.Message,
		Data:     entry.Data,
	})
	if err != nil {
		return nil, err
	}
	return append(data, "\n"...), nil
}

// FormatText 日志文本格式化
func (f *Formatter) FormatText(entry *log.Entry) ([]byte, error) {
	buffer := bytes.Buffer{}
	time, level := entry.Time.Format(f.TimeLayout), LevelString(entry.Level, 5)
	buffer.WriteString(fmt.Sprintf("[%-23s][%s][%-5s]", time, f.Hostname, level))
	buffer.WriteString(entry.Message)
	for key, value := range entry.Data {
		buffer.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}

// Color 日志颜色，根据日志级别对日志内容染色
func Color(level log.Level, data []byte) []byte {
	var color string
	switch level {
	case log.InfoLevel:
		color = "\x1b[32m" // green
	case log.WarnLevel:
		color = "\x1b[33m" // yellow
	case log.ErrorLevel:
		color = "\x1b[31m" // red
	default:
		return data
	}
	buffer := bytes.Buffer{}
	buffer.WriteString(color)
	buffer.Write(data)
	buffer.WriteString("\x1b[0m")
	return buffer.Bytes()
}

// LevelString 日志级别字符串
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
