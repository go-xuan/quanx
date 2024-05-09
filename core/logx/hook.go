package logx

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/stringx"
)

type Writers map[logrus.Level]io.Writer

type Hook struct {
	lock      *sync.Mutex
	writers   Writers
	levels    []logrus.Level
	formatter logrus.Formatter
}

func NewHook(writer any, formatter logrus.Formatter) *Hook {
	hook := &Hook{lock: new(sync.Mutex)}
	hook.SetFormatter(formatter)
	switch writer.(type) {
	case string:
		hook.SetWriter(&FileWriter{writer.(string)})
	case map[logrus.Level]string:
		var writers = make(Writers)
		for level, path := range writer.(map[logrus.Level]string) {
			writers[level] = &FileWriter{path}
		}
		hook.SetWriters(writers)
	case io.Writer:
		hook.SetWriter(writer.(io.Writer))
	case Writers:
		hook.SetWriters(writer.(Writers))
	default:
		panic(fmt.Sprintf("unsupported writer type: %v", reflect.TypeOf(writer)))
	}
	return hook
}

func (hook *Hook) Levels() []logrus.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	var caller = getCaller()
	_, fileName := stringx.Cut(caller.File, "/", -1)
	_, funcName := stringx.Cut(caller.Function, ".", -1)
	entry.WithField("position", fmt.Sprintf("%s:%04d:%s()", fileName, caller.Line, funcName))
	hook.lock.Lock()
	defer hook.lock.Unlock()
	return hook.Write(entry)
}

func (hook *Hook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = &logrus.TextFormatter{
			DisableColors:          true,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			DisableSorting:         false,
		}
	} else if textFormatter, ok := formatter.(*logrus.TextFormatter); ok {
		textFormatter.DisableColors = true
	}
	hook.formatter = formatter
}

func (hook *Hook) SetWriters(writers Writers) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.writers = writers
	for level := range writers {
		hook.levels = append(hook.levels, level)
	}
}

func (hook *Hook) SetWriter(writer io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.writers = map[logrus.Level]io.Writer{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}
	hook.levels = AllLevels()
}

// 输出到ioWriter
func (hook *Hook) Write(entry *logrus.Entry) (err error) {
	if hook.writers != nil {
		if writer, ok := hook.writers[entry.Level]; ok {
			var bytes []byte
			if bytes, err = hook.formatter.Format(entry); err != nil {
				return
			}
			_, err = writer.Write(bytes)
		}
	}
	return
}

var (
	callerPkg  string
	callerOnce sync.Once
)

func getCaller() *runtime.Frame {
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	callerOnce.Do(func() {
		callerPkg = getPackageName(runtime.FuncForPC(pcs[0]).Name())
	})
	for f, again := frames.Next(); again; f, again = frames.Next() {
		if pkg := getPackageName(f.Function); pkg != callerPkg {
			return &f
		}
	}
	return nil
}

func getPackageName(f string) string {
	for {
		lastPeriod, lastSlash := stringx.Index(f, ".", -1), stringx.Index(f, "/", -1)
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}
