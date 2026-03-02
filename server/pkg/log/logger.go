package log

import (
	"mvtable/internal/pkg/config"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
)

// logger 全局日志实例
var (
	logger *zap.Logger
	once   sync.Once
)

func Init(cfg *config.LogConfig) error {
	var err error
	once.Do(func() {
		// 设置日志级别
		var zapLevel zapcore.Level
		switch cfg.Level {
		case "debug":
			zapLevel = zapcore.DebugLevel
		case "info":
			zapLevel = zapcore.InfoLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		case "panic":
			zapLevel = zapcore.PanicLevel
		case "fatal":
			zapLevel = zapcore.FatalLevel
		default:
			zapLevel = zapcore.InfoLevel
		}

		// 设置编码器配置
		var encoderConfig zapcore.EncoderConfig
		if cfg.Format == "json" {
			encoderConfig = zap.NewProductionEncoderConfig()
		} else {
			encoderConfig = zap.NewDevelopmentEncoderConfig()
		}

		// 自定义时间格式
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		// 创建编码器
		var encoder zapcore.Encoder
		if cfg.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		// 设置输出
		var writeSyncer zapcore.WriteSyncer
		if cfg.Output == "file" {
			// 确保日志目录存在
			if err = os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
				return
			}

			// 使用lumberjack进行日志轮转
			fileWriter := &lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,
				MaxAge:     cfg.MaxAge,
				MaxBackups: cfg.MaxBackups,
				Compress:   cfg.Compress,
			}
			writeSyncer = zapcore.AddSync(fileWriter)
		} else {
			writeSyncer = zapcore.AddSync(os.Stdout)
		}

		// 创建核心
		core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

		// 创建logger
		logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})
	return err
}

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 记录错误级别日志
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Panic 记录恐慌级别日志
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal 记录致命级别日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Sync 同步日志缓冲区
func Sync() error {
	return logger.Sync()
}

// With 创建带有字段的logger
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sugar 获取sugared logger
func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}
