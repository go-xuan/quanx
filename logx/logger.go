package logx

import (
	"context"

	"github.com/go-xuan/utilx/contextx"
	log "github.com/sirupsen/logrus"
)

// NewLogger 创建日志实例
func NewLogger() *Logger {
	return &Logger{
		entity: log.WithField("app", GetConfig().Name),
	}
}

// NewCtxLogger 创建日志实例
func NewCtxLogger(ctx context.Context) *Logger {
	return NewLogger().WithContext(ctx)
}

// Logger 日志实例
type Logger struct {
	entity *log.Entry
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if app := contextx.GetValue(ctx, "app"); app.Valid() {
		l.entity = l.entity.WithField("app", app)
	}
	if trace := contextx.GetValue(ctx, "trace"); trace.Valid() {
		l.entity = l.entity.WithField("trace", trace)
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
	l.entity = l.entity.WithError(err)
	return l
}

func (l *Logger) WithErrorMsg(err error) *Logger {
	l.entity = l.entity.WithField(log.ErrorKey, err.Error())
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
