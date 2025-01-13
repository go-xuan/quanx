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

func newHook() *Hook {
	return &Hook{
		lock:    new(sync.Mutex),
		writers: make(map[log.Level]io.Writer),
	}
}

func NewHook(writer any, formatter log.Formatter) *Hook {
	hook := newHook()
	hook.SetFormatter(formatter)
	switch writer.(type) {
	case string:
		hook.InitWriter(&FileWriter{writer.(string)})
	case map[log.Level]string:
		for level, path := range writer.(map[log.Level]string) {
			hook.SetWriter(level, &FileWriter{path})
		}
	case Writers:
		hook.InitWriters(writer.(Writers))
	case io.Writer:
		hook.InitWriter(writer.(io.Writer))
	default:
		panic(fmt.Sprintf("unsupported writer type: %v", reflect.TypeOf(writer)))
	}
	return hook
}

type Writers map[log.Level]io.Writer

type Hook struct {
	lock      *sync.Mutex
	writers   Writers
	levels    []log.Level
	formatter log.Formatter
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

func (hook *Hook) InitWriter(writer io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.levels = AllLogrusLevels()
	for _, level := range hook.levels {
		hook.writers[level] = writer
	}
}

func (hook *Hook) InitWriters(writers Writers) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.writers = writers
	for level := range writers {
		hook.levels = append(hook.levels, level)
	}
}

func (hook *Hook) SetWriter(level log.Level, writer io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if _, ok := hook.writers[level]; !ok {
		hook.levels = append(hook.levels, level)
	}
	hook.writers[level] = writer
}

// 输出到ioWriter
func (hook *Hook) Write(entry *log.Entry) error {
	if hook.writers != nil {
		if writer, ok := hook.writers[entry.Level]; ok {
			if bytes, err := hook.formatter.Format(entry); err != nil {
				return err
			} else if _, err = writer.Write(bytes); err != nil {
				return err
			}
		}
	}
	return nil
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
