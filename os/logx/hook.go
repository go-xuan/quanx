package logx

import (
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/stringx"
)

type Writers map[log.Level]io.Writer

type Hook struct {
	lock      *sync.Mutex
	writers   Writers
	levels    []log.Level
	formatter log.Formatter
}

func NewHook(writer any, formatter log.Formatter) *Hook {
	hook := &Hook{lock: new(sync.Mutex)}
	hook.SetFormatter(formatter)
	switch writer.(type) {
	case string:
		hook.SetWriter(&FileWriter{writer.(string)})
	case map[log.Level]string:
		var writers = make(Writers)
		for level, path := range writer.(map[log.Level]string) {
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

func (hook *Hook) Levels() []log.Level {
	return hook.levels
}

func (hook *Hook) Fire(entry *log.Entry) error {
	var caller = getCaller()
	_, fileName := stringx.Cut(caller.File, "/", -1)
	_, funcName := stringx.Cut(caller.Function, ".", -1)
	entry.WithField("position", fmt.Sprintf("%s:%04d:%s()", fileName, caller.Line, funcName))
	hook.lock.Lock()
	defer hook.lock.Unlock()
	return hook.Write(entry)
}

func (hook *Hook) SetFormatter(formatter log.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = &log.TextFormatter{
			DisableColors:          true,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
			DisableSorting:         false,
		}
	} else if textFormatter, ok := formatter.(*log.TextFormatter); ok {
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
	hook.writers = map[log.Level]io.Writer{
		log.TraceLevel: writer,
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}
	hook.levels = AllLevels()
}

// 输出到ioWriter
func (hook *Hook) Write(entry *log.Entry) (err error) {
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
	pcs := make([]uintptr, 32)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	callerOnce.Do(func() {
		callerPkg = getPackageName(runtime.FuncForPC(pcs[0]).Name())
	})
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if pkg := getPackageName(frame.Function); pkg != callerPkg {
			return &frame
		}
	}
	return nil
}

func getPackageName(function string) string {
	for {
		period, slash := stringx.Index(function, ".", -1), stringx.Index(function, "/", -1)
		if period > slash {
			function = function[:period]
		} else {
			break
		}
	}
	return function
}
