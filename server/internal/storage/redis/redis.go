package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"mvtable/internal/pkg/config"
	"mvtable/pkg/log"
)

var (
	Client *redis.Client
	once   sync.Once
	Nil    = redis.Nil
)

// Init 初始化Redis连接
func Init(cfg *config.RedisConfig) error {
	var initErr error
	once.Do(func() {
		Client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.Database,
			PoolSize: cfg.PoolSize,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := Client.Ping(ctx).Result()
		if err != nil {
			initErr = fmt.Errorf("redis连接失败: %w", err)
			return
		}

		log.Info("Redis连接成功",
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
			zap.Int("database", cfg.Database),
		)
	})
	return initErr
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// GetClient 获取Redis客户端
func GetClient() *redis.Client {
	return Client
}

// Set 设置键值
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return Client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
	return Client.Exists(ctx, keys...).Result()
}

// Expire 设置键过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return Client.Expire(ctx, key, expiration).Err()
}

// TTL 获取键剩余生存时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return Client.TTL(ctx, key).Result()
}
