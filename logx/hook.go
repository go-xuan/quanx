package logx

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/go-xuan/utilx/stringx"
	log "github.com/sirupsen/logrus"
)

// NewHook 创建日志钩子
func NewHook(formatter log.Formatter) *Hook {
	return &Hook{
		lock:      new(sync.Mutex),
		writers:   make(map[log.Level]io.Writer),
		levels:    make([]log.Level, 0),
		formatter: formatter,
	}
}

// HookConfig 日志钩子配置
type HookConfig struct {
	Writer string   `json:"writer" yaml:"writer" default:"console"`
	Levels []string `json:"levels" yaml:"levels"`
}

// NewHook 创建日志钩子
func (c *HookConfig) NewHook(name string, formatter log.Formatter) *Hook {
	hook := NewHook(formatter)
	for _, lv := range c.Levels {
		level := LogrusLevel(lv)
		hook.levels = append(hook.levels, level)
		if writer := NewWriter(c.Writer, name, lv); writer != nil {
			hook.writers[level] = writer
		}
	}
	return hook
}

// Hook 日志钩子
type Hook struct {
	lock      *sync.Mutex
	writers   map[log.Level]io.Writer
	levels    []log.Level
	formatter log.Formatter
}

// Levels 获取日志钩子级别
func (h *Hook) Levels() []log.Level {
	return h.levels
}

// Fire 日志钩子触发
func (h *Hook) Fire(entry *log.Entry) error {
	if caller := getCaller(); caller != nil {
		_, fileName := stringx.Cut(caller.File, "/", -1)
		_, funcName := stringx.Cut(caller.Function, ".", -1)
		entry.WithField("position", fmt.Sprintf("%s:%04d:%s()", fileName, caller.Line, funcName))
	}
	if writer, ok := h.writers[entry.Level]; ok {
		if bytes, err := h.formatter.Format(entry); err != nil {
			return err
		} else if _, err = writer.Write(bytes); err != nil {
			return err
		}
	}
	return nil
}

// SetFormatter 设置日志格式化器
func (h *Hook) SetFormatter(formatter log.Formatter) {
	if formatter == nil {
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	h.formatter = formatter
}

// AddWriter 添加日志hook级别以及Writer
func (h *Hook) AddWriter(level log.Level, writer io.Writer) {
	if writer == nil {
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	if _, ok := h.writers[level]; !ok {
		h.levels = append(h.levels, level)
	}
	h.writers[level] = writer
}

var (
	callerPkg  string
	callerOnce sync.Once
)

// 获取调用栈信息
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

// 获取包名
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
