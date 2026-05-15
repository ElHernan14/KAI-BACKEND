package database

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	slowThreshold time.Duration
	logLevel      logger.LogLevel
}

func NewGormLogger() GormLogger {
	return GormLogger{
		slowThreshold: 200 * time.Millisecond,
		logLevel:      logger.Warn,
	}
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		log.Printf(msg, data...)
	}
}

func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		log.Printf(msg, data...)
	}
}

func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		log.Printf(msg, data...)
	}
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		log.Printf("gorm_error duration=%.3fs rows=%d error=%v sql=%s", elapsed.Seconds(), rows, err, sql)
	case elapsed > l.slowThreshold && l.logLevel >= logger.Warn:
		log.Printf("gorm_slow_sql threshold=%.3fs duration=%.3fs rows=%d sql=%s", l.slowThreshold.Seconds(), elapsed.Seconds(), rows, sql)
	case l.logLevel >= logger.Info:
		log.Printf("gorm_query duration=%.3fs rows=%d sql=%s", elapsed.Seconds(), rows, sql)
	}
}
