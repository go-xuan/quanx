package logx

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/file/filex"
	"github.com/go-xuan/quanx/types/stringx"
)

type PathMap map[logrus.Level]string

type WriterMap map[logrus.Level]io.Writer

type Hook struct {
	lock             *sync.Mutex
	paths            PathMap
	defaultPath      string
	hasDefaultPath   bool
	writers          WriterMap
	defaultWriter    io.Writer
	hasDefaultWriter bool
	levels           []logrus.Level
	formatter        logrus.Formatter
}

func NewHook(output any, formatter logrus.Formatter) *Hook {
	hook := &Hook{
		lock: new(sync.Mutex),
	}
	hook.SetFormatter(formatter)

	switch output.(type) {
	case string:
		hook.SetDefaultPath(output.(string))
	case io.Writer:
		hook.SetDefaultWriter(output.(io.Writer))
	case PathMap:
		hook.paths = output.(PathMap)
		for level := range output.(PathMap) {
			hook.levels = append(hook.levels, level)
		}
	case WriterMap:
		hook.writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.levels = append(hook.levels, level)
		}
	default:
		panic(fmt.Sprintf("unsupported output type: %v", reflect.TypeOf(output)))
	}
	return hook
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
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}
	hook.formatter = formatter
}

func (hook *Hook) SetDefaultPath(defaultPath string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultPath = defaultPath
	hook.hasDefaultPath = true
}

func (hook *Hook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
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
	if hook.writers != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	} else if hook.paths != nil || hook.hasDefaultPath {
		return hook.fileWrite(entry)
	}
	return nil
}

// 写入ioWriter
func (hook *Hook) ioWrite(entry *logrus.Entry) (err error) {
	if writer, ok := hook.writers[entry.Level]; !ok {
		if hook.hasDefaultWriter {
			writer = hook.defaultWriter
		} else {
			return
		}
	} else {
		var bytes []byte
		if bytes, err = hook.formatter.Format(entry); err != nil {
			return
		}
		if _, err = writer.Write(bytes); err != nil {
			return
		}
	}
	return
}

// 写入文件
func (hook *Hook) fileWrite(entry *logrus.Entry) (err error) {
	if path, ok := hook.paths[entry.Level]; !ok {
		if hook.hasDefaultPath {
			path = hook.defaultPath
		} else {
			return
		}
	} else {
		filex.CreateDirNotExist(path)
		var file *os.File
		if file, err = os.OpenFile(path, filex.AppendOnly, 0666); err != nil {
			return
		}
		defer file.Close()
		var bytes []byte
		if bytes, err = hook.formatter.Format(entry); err != nil {
			return
		}
		if _, err = file.Write(bytes); err != nil {
			return
		}
	}
	return
}

var (
	pcPkg      string
	callerOnce sync.Once
)

func getCaller() *runtime.Frame {
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	callerOnce.Do(func() {
		pcPkg = getPackageName(runtime.FuncForPC(pcs[0]).Name())
	})
	for f, again := frames.Next(); again; f, again = frames.Next() {
		if pkg := getPackageName(f.Function); pkg != pcPkg {
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
