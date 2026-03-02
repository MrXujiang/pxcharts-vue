package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"mvtable/internal/app/collaboration/model"
	"mvtable/internal/storage/redis"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

const (
	// LockExpiration 锁的过期时间（秒）
	LockExpiration = 300 // 5分钟
	// LockRefreshInterval 锁刷新间隔（秒）
	LockRefreshInterval = 60 // 1分钟
)

// LockManager 锁管理器
type LockManager struct{}

// NewLockManager 创建锁管理器
func NewLockManager() *LockManager {
	return &LockManager{}
}

// getLockKey 获取锁的 Redis key
func (lm *LockManager) getLockKey(tableID, recordID, fieldID string) string {
	return fmt.Sprintf("collab:lock:%s:%s:%s", tableID, recordID, fieldID)
}

// AcquireLock 申请锁
func (lm *LockManager) AcquireLock(ctx context.Context, tableID, recordID, fieldID, userID string) (bool, error) {
	key := lm.getLockKey(tableID, recordID, fieldID)

	// 尝试获取当前锁的持有者
	currentUserID, err := redis.Get(ctx, key)
	if err != nil && err != redis.Nil {
		log.Error("get lock error", zap.Error(err))
		return false, err
	}

	// 如果锁已存在且不是当前用户持有，返回失败
	if err == nil && currentUserID != userID {
		return false, nil
	}

	// 设置锁（如果不存在则设置，如果存在且是当前用户则更新过期时间）
	err = redis.Set(ctx, key, userID, LockExpiration*time.Second)
	if err != nil {
		log.Error("set lock error", zap.Error(err))
		return false, err
	}

	return true, nil
}

// ReleaseLock 释放锁
func (lm *LockManager) ReleaseLock(ctx context.Context, tableID, recordID, fieldID, userID string) error {
	key := lm.getLockKey(tableID, recordID, fieldID)

	// 检查锁是否由当前用户持有
	currentUserID, err := redis.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			// 锁不存在，直接返回成功
			return nil
		}
		log.Error("get lock error", zap.Error(err))
		return err
	}

	// 如果锁不是当前用户持有，返回错误
	if currentUserID != userID {
		return fmt.Errorf("lock is not held by user %s", userID)
	}

	// 删除锁
	err = redis.Del(ctx, key)
	if err != nil {
		log.Error("delete lock error", zap.Error(err))
		return err
	}

	return nil
}

// RefreshLock 刷新锁的过期时间
func (lm *LockManager) RefreshLock(ctx context.Context, tableID, recordID, fieldID, userID string) error {
	key := lm.getLockKey(tableID, recordID, fieldID)

	// 检查锁是否由当前用户持有
	currentUserID, err := redis.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("lock does not exist")
		}
		log.Error("get lock error", zap.Error(err))
		return err
	}

	// 如果锁不是当前用户持有，返回错误
	if currentUserID != userID {
		return fmt.Errorf("lock is not held by user %s", userID)
	}

	// 更新锁的过期时间
	err = redis.Set(ctx, key, userID, LockExpiration*time.Second)
	if err != nil {
		log.Error("refresh lock error", zap.Error(err))
		return err
	}

	return nil
}

// GetLockOwner 获取锁的持有者
func (lm *LockManager) GetLockOwner(ctx context.Context, tableID, recordID, fieldID string) (string, error) {
	key := lm.getLockKey(tableID, recordID, fieldID)
	userID, err := redis.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return "", nil // 锁不存在
		}
		return "", err
	}
	return userID, nil
}

// ReleaseAllLocksByField 释放指定字段的所有锁
func (lm *LockManager) ReleaseAllLocksByField(ctx context.Context, tableID, fieldID string) ([]model.LockInfo, error) {
	pattern := fmt.Sprintf("collab:lock:%s:*:%s", tableID, fieldID)

	var releasedLocks []model.LockInfo
	var cursor uint64 = 0
	var count int64 = 100 // 每次扫描的数量

	for {
		// 使用 SCAN 命令遍历匹配的键
		keys, nextCursor, err := redis.GetClient().Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			log.Error("scan locks error", zap.Error(err))
			return nil, err
		}

		// 解析锁信息并删除
		for _, key := range keys {
			// 从 key 中提取信息：collab:lock:{tableID}:{recordID}:{fieldID}
			parts := strings.Split(key, ":")
			if len(parts) == 5 && parts[0] == "collab" && parts[1] == "lock" {
				recordID := parts[3]

				// 获取锁的持有者（用于通知）
				userID, err := redis.Get(ctx, key)
				if err == nil && userID != "" {
					releasedLocks = append(releasedLocks, model.LockInfo{
						RecordID: recordID,
						FieldID:  fieldID,
						UserID:   userID,
					})
				}

				// 删除锁
				if err := redis.Del(ctx, key); err != nil {
					log.Warn("delete lock error", zap.String("key", key), zap.Error(err))
				}
			}
		}

		// 如果 cursor 为 0，说明扫描完成
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	log.Info("released locks by field",
		zap.String("tableID", tableID),
		zap.String("fieldID", fieldID),
		zap.Int("count", len(releasedLocks)))

	return releasedLocks, nil
}

// ReleaseAllLocksByUser 释放指定用户持有的所有锁，返回按 tableID 分组的释放锁信息
func (lm *LockManager) ReleaseAllLocksByUser(ctx context.Context, userID string) (map[string][]model.LockInfo, error) {
	pattern := "collab:lock:*:*:*"
	cursor := uint64(0)
	var count int64 = 100

	releasedByTable := make(map[string][]model.LockInfo)

	for {
		keys, nextCursor, err := redis.GetClient().Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			log.Error("scan locks error", zap.Error(err))
			return nil, err
		}

		for _, key := range keys {
			parts := strings.Split(key, ":")
			// expect format collab:lock:{tableID}:{recordID}:{fieldID}
			if len(parts) != 5 {
				continue
			}
			tableID := parts[2]
			recordID := parts[3]
			fieldID := parts[4]

			owner, err := redis.Get(ctx, key)
			if err != nil {
				if err == redis.Nil {
					continue
				}
				log.Warn("get lock owner error", zap.String("key", key), zap.Error(err))
				continue
			}
			if owner != userID {
				continue
			}

			// 删除锁
			if err := redis.Del(ctx, key); err != nil {
				log.Warn("delete lock error", zap.String("key", key), zap.Error(err))
			}

			releasedByTable[tableID] = append(releasedByTable[tableID], model.LockInfo{
				RecordID: recordID,
				FieldID:  fieldID,
				UserID:   owner,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return releasedByTable, nil
}
