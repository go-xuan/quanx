package logx

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogFormatter(timeFormat string) *LogFormatter {
	return &LogFormatter{timeFormat}
}

type LogFormatter struct {
	timeFormat string
}

// 日志格式化,用以实现logrus.Formatter接口
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	host, _ := os.Hostname()
	var b = bytes.Buffer{}
	b.WriteString(fmt.Sprintf("[%23s][%7s][%s]", time.Now().Format(f.timeFormat), entry.Level.String(), host))
	b.WriteString(entry.Message)
	for key, value := range entry.Data {
		b.WriteString(fmt.Sprintf(", %s:%+v", key, value))
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
