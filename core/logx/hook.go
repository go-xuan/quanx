package logx

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/types/stringx"
)

type OutputMap map[logrus.Level]*Output

type Hook struct {
	lock      *sync.Mutex
	output    OutputMap
	levels    []logrus.Level
	formatter logrus.Formatter
}

func NewHook(output any, formatter logrus.Formatter) *Hook {
	hook := &Hook{lock: new(sync.Mutex)}
	hook.SetFormatter(formatter)
	switch output.(type) {
	case *Output:
		hook.SetOutput(output.(*Output))
	case OutputMap:
		hook.SetOutputMap(output.(OutputMap))
	default:
		panic(fmt.Sprintf("unsupported output type: %v", reflect.TypeOf(output)))
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
	if hook.output != nil {
		return hook.OutputToWriter(entry)
	}
	return nil
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

func (hook *Hook) SetOutputMap(outputMap OutputMap) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.output = outputMap
	for level := range outputMap {
		hook.levels = append(hook.levels, level)
	}
}

func (hook *Hook) SetOutput(output *Output) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.output = map[logrus.Level]*Output{
		logrus.TraceLevel: output,
		logrus.DebugLevel: output,
		logrus.InfoLevel:  output,
		logrus.WarnLevel:  output,
		logrus.ErrorLevel: output,
		logrus.FatalLevel: output,
		logrus.PanicLevel: output,
	}
	hook.levels = AllLevels()
}

// 输出到ioWriter
func (hook *Hook) OutputToWriter(entry *logrus.Entry) (err error) {
	if output, ok := hook.output[entry.Level]; ok {
		var bytes []byte
		if bytes, err = hook.formatter.Format(entry); err != nil {
			return
		}
		_, err = output.Writer.Write(bytes)
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
