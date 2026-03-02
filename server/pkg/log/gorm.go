package log

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	zapLogger     *zap.Logger
	logLevel      gormLogger.LogLevel
	slowThreshold time.Duration
}

func NewGormLogger(logLevel gormLogger.LogLevel, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		zapLogger:     logger, // 使用全局 logger
		logLevel:      logLevel,
		slowThreshold: slowThreshold,
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Info {
		l.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Warn {
		l.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Error {
		l.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && !errors.Is(err, gormLogger.ErrRecordNotFound):
		l.zapLogger.Sugar().Errorf("[%.3fms] [rows:%v] %s | err: %v", float64(elapsed.Milliseconds()), rows, sql, err)
	case l.slowThreshold != 0 && elapsed > l.slowThreshold:
		l.zapLogger.Sugar().Warnf("[SLOW QUERY %.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
	default:
		if l.logLevel == gormLogger.Info {
			l.zapLogger.Sugar().Infof("[%.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
		}
	}
}
