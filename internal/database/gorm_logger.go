package database

import (
	"context"
	"errors"
	"log"
	"strings"
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
	table := extractTableName(sql)
	operation := extractOperation(sql)

	switch {
	case err != nil && l.logLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		log.Printf("DB GORM error | operation=%s | table=%s | duration=%.3fs | rowsAffected=%d | error=%v", operation, table, elapsed.Seconds(), rows, err)
	case elapsed > l.slowThreshold && l.logLevel >= logger.Warn:
		log.Printf("DB GORM slow query | operation=%s | table=%s | limit=%.3fs | duration=%.3fs | rowsAffected=%d", operation, table, l.slowThreshold.Seconds(), elapsed.Seconds(), rows)
	case l.logLevel >= logger.Info:
		log.Printf("DB GORM query | operation=%s | table=%s | duration=%.3fs | rowsAffected=%d", operation, table, elapsed.Seconds(), rows)
	}
}

func extractOperation(sql string) string {
	fields := strings.Fields(sql)
	if len(fields) == 0 {
		return "UNKNOWN"
	}

	return strings.ToUpper(fields[0])
}

func extractTableName(sql string) string {
	fields := strings.Fields(sql)
	for i, field := range fields {
		keyword := strings.ToUpper(field)
		if (keyword == "FROM" || keyword == "INTO" || keyword == "UPDATE") && i+1 < len(fields) {
			return cleanTableName(fields[i+1])
		}
	}

	return "unknown"
}

func cleanTableName(table string) string {
	table = strings.Trim(table, `"`)
	table = strings.Trim(table, "`")
	table = strings.TrimSuffix(table, ",")
	table = strings.TrimSuffix(table, ";")

	if strings.Contains(table, ".") {
		parts := strings.Split(table, ".")
		table = parts[len(parts)-1]
		table = strings.Trim(table, `"`)
	}

	return table
}
