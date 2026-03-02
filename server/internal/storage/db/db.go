package db

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"mvtable/internal/pkg/config"
	"mvtable/pkg/log"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init 初始化数据库连接
func Init(cfg *config.DatabaseConfig) error {
	var initErr error
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database,
		)

		gormLogger := log.NewGormLogger(logger.Info, 200*time.Millisecond)
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})
		if err != nil {
			initErr = fmt.Errorf("连接数据库失败: %w", err)
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			initErr = fmt.Errorf("获取数据库实例失败: %w", err)
			return
		}

		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

		if err = sqlDB.Ping(); err != nil {
			initErr = fmt.Errorf("数据库连接测试失败: %w", err)
			return
		}

		log.Info("数据库连接成功",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.String("database", cfg.Database),
		)
	})
	return initErr
}

func GetDB() *gorm.DB {
	return db
}

func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func Transaction(fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}
