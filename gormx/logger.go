package gormx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

const (
	Debug = logger.Info + 1
)

// DefaultLogger 默认日志
func DefaultLogger() *Logger {
	return &Logger{
		LogLevel:      logger.Silent,
		SlowThreshold: 200 * time.Millisecond,
	}
}

// LogLevel 日志级别映射，默认Silent
func LogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return Debug
	case "info":
		return logger.Info
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	default:
		return logger.Silent
	}
}

// Logger 日志
type Logger struct {
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *Logger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.WithContext(ctx).Infof(format, args...)
	}
}

func (l *Logger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.WithContext(ctx).Warnf(format, args...)
	}
}

func (l *Logger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.WithContext(ctx).Errorf(format, args...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, affected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && !errors.Is(err, logger.ErrRecordNotFound):
		sql, rows := fc()
		log.WithContext(ctx).WithFields(log.Fields{
			"location": utils.FileWithLineNum(),
			"rows":     rows,
			"elapsed":  elapsed.String(),
			"sql":      sql,
		}).Error(err.Error())
	case elapsed > l.SlowThreshold && l.SlowThreshold <= 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		log.WithContext(ctx).WithFields(log.Fields{
			"location": utils.FileWithLineNum(),
			"rows":     rows,
			"elapsed":  elapsed.String(),
			"sql":      sql,
		}).Warnf("slow sql more than %v", l.SlowThreshold)
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		log.WithContext(ctx).WithFields(log.Fields{
			"location": utils.FileWithLineNum(),
			"rows":     rows,
			"elapsed":  elapsed.String(),
			"sql":      sql,
		}).Info()
	case l.LogLevel == Debug:
		sql, rows := fc()
		fmt.Printf("[GORM-DEBUG] %s [%s] [%d] %s \n", utils.FileWithLineNum(), elapsed.String(), rows, sql)
	}
}
