package logx

import (
	"context"

	"github.com/go-xuan/utilx/contextx"
	"github.com/go-xuan/utilx/idx"
	log "github.com/sirupsen/logrus"
)

func NewEntry() *log.Entry {
	return log.WithField("app", _config.Name)
}

// NewLogger 创建日志实例
func NewLogger() *Logger {
	return &Logger{
		app:    _config.Name,
		trace:  idx.UUID(),
		entity: NewEntry(),
	}
}

// NewCtxLogger 创建日志实例
func NewCtxLogger(ctx context.Context) *Logger {
	return NewLogger().WithContext(ctx)
}

type Logger struct {
	app    string
	trace  string
	entity *log.Entry
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if app := contextx.GetValue(ctx, "app"); app.Valid() {
		l.app = app.String()
		l.entity = l.entity.WithField("app", l.app)
	}
	if trace := contextx.GetValue(ctx, "trace"); trace.Valid() {
		l.trace = trace.String()
		l.entity = l.entity.WithField("trace", l.trace)
	}
	l.entity = l.entity.WithContext(ctx)
	return l
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.entity = l.entity.WithField(key, value)
	return l
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.entity = l.entity.WithFields(fields)
	return l
}

func (l *Logger) WithError(err error) *Logger {
	l.entity = l.entity.WithField("error", err.Error())
	return l
}

func (l *Logger) WithEntity(entity Entry) *Logger {
	l.entity = log.WithFields(entity.LogFields())
	return l
}

func (l *Logger) Trace(args ...interface{}) {
	l.entity.Trace(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.entity.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.entity.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.entity.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entity.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entity.Fatal(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.entity.Panic(args...)
}

type Entry interface {
	LogFields() log.Fields // 日志字段
}

func WithEntity(entity Entry) *log.Entry {
	return log.WithFields(entity.LogFields())
}
