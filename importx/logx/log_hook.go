package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// We are logging to file, strip colors to make the output more readable.
var defaultFormatter = &logrus.TextFormatter{
	DisableColors:          true,
	DisableTimestamp:       true,
	DisableLevelTruncation: true,
	DisableSorting:         false}

// is map for mapping a log level to a file's path.
// Multiple levels may share a file, but multiple files may not be used for one level.
type PathMap map[logrus.Level]string

// is map for mapping a log level to an io.Writer.
// Multiple levels may share a writer, but multiple writers may not be used for one level.
type WriterMap map[logrus.Level]io.Writer

// LfsHook is a hook to handle writing to local log files.
type LfsHook struct {
	paths     PathMap
	writers   WriterMap
	levels    []logrus.Level
	lock      *sync.Mutex
	formatter logrus.Formatter

	defaultPath      string
	defaultWriter    io.Writer
	hasDefaultPath   bool
	hasDefaultWriter bool
}

// returns new LFS hook.
// Output can be a string, io.Writer, WriterMap or PathMap.
// If using io.Writer or WriterMap, user is responsible for closing the used io.Writer.
func NewHook(output interface{}, formatter logrus.Formatter) *LfsHook {
	hook := &LfsHook{
		lock: new(sync.Mutex),
	}

	hook.SetFormatter(formatter)

	switch output.(type) {
	case string:
		hook.SetDefaultPath(output.(string))
		break
	case io.Writer:
		hook.SetDefaultWriter(output.(io.Writer))
		break
	case PathMap:
		hook.paths = output.(PathMap)
		for level := range output.(PathMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	case WriterMap:
		hook.writers = output.(WriterMap)
		for level := range output.(WriterMap) {
			hook.levels = append(hook.levels, level)
		}
		break
	default:
		panic(fmt.Sprintf("unsupported level map type: %v", reflect.TypeOf(output)))
	}

	return hook
}

// sets the format that will be used by hook.
// If using text formatter, this method will disable color output to make the log file more readable.
func (hook *LfsHook) SetFormatter(formatter logrus.Formatter) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if formatter == nil {
		formatter = defaultFormatter
	} else {
		switch formatter.(type) {
		case *logrus.TextFormatter:
			textFormatter := formatter.(*logrus.TextFormatter)
			textFormatter.DisableColors = true
		}
	}

	hook.formatter = formatter
}

// sets default path for levels that don't have any defined output path.
func (hook *LfsHook) SetDefaultPath(defaultPath string) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultPath = defaultPath
	hook.hasDefaultPath = true
}

// sets default writer for levels that don't have any defined writer.
func (hook *LfsHook) SetDefaultWriter(defaultWriter io.Writer) {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	hook.defaultWriter = defaultWriter
	hook.hasDefaultWriter = true
}

// writes the log file to defined path or using the defined writer.
// Title who run this function needs to write permissions to the file or directory if the file does not yet exist.
func (hook *LfsHook) Fire(entry *logrus.Entry) error {

	caller := getCaller()
	splitsFile := strings.Split(caller.File, "/")
	splitsFunc := strings.Split(caller.Function, ".")
	callStr := fmt.Sprintf("%s:%04d:%s()", splitsFile[len(splitsFile)-1], caller.Line, splitsFunc[len(splitsFunc)-1]) //caller.File + fmt.Sprintf("%d", caller.Line) + caller.Function
	entry.Data["source"] = fmt.Sprintf("%s", callStr)
	hook.lock.Lock()
	defer hook.lock.Unlock()
	if hook.writers != nil || hook.hasDefaultWriter {
		return hook.ioWrite(entry)
	} else if hook.paths != nil || hook.hasDefaultPath {
		return hook.fileWrite(entry)
	}
	return nil
}

// Write a log line to an io.Writer.
func (hook *LfsHook) ioWrite(entry *logrus.Entry) (err error) {
	var (
		writer io.Writer
		msg    []byte
		ok     bool
	)

	if writer, ok = hook.writers[entry.Level]; !ok {
		if hook.hasDefaultWriter {
			writer = hook.defaultWriter
		} else {
			return
		}
	}

	if msg, err = hook.formatter.Format(entry); err != nil {
		log.Println("failed to generate string for entry:", err)
		return
	}
	if _, err = writer.Write(msg); err != nil {
		return
	}
	return
}

// Write a log line directly to a file.
func (hook *LfsHook) fileWrite(entry *logrus.Entry) (err error) {
	var (
		fd   *os.File
		path string
		msg  []byte
		ok   bool
	)

	if path, ok = hook.paths[entry.Level]; !ok {
		if hook.hasDefaultPath {
			path = hook.defaultPath
		} else {
			return
		}
	}

	dir := filepath.Dir(path)
	_ = os.MkdirAll(dir, os.ModePerm)

	if fd, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		log.Println("failed to open logfile:", path, err)
		return
	}
	defer fd.Close()
	if msg, err = hook.formatter.Format(entry); err != nil {
		log.Println("failed to generate string for entry:", err)
		return
	}
	if _, err = fd.Write(msg); err != nil {
		return
	}
	return
}

// Levels returns configured log levels.
func (hook *LfsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var (
	// qualified package name, cached at first use
	logrusPackage string
	// Used for caller information initialisation
	callerInitOnce sync.Once
)

// reduces a fully qualified function name to the package name
// There really ought to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// Restrict the look back frames to avoid runaway lookups
	pcs := make([]uintptr, 25)
	depth := runtime.Callers(4, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		logrusPackage = getPackageName(runtime.FuncForPC(pcs[0]).Name())
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage {
			return &f
		}
	}
	// if we got here, we failed to find the caller's context
	return nil
}
