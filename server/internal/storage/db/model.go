package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        string         `gorm:"column:id;type:varchar(36);primaryKey;not null" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamptz;index:idx_user_deleted_at" json:"-"`
}

// BeforeCreate 在创建前设置ID
func (m *Model) BeforeCreate(*gorm.DB) error {
	if m.ID == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			return fmt.Errorf("uuid create failed, %w", err)
		}
		m.ID = id.String()
	}

	return nil
}

type Pagination struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

func Create[T any](db *gorm.DB, data *T) error {
	return db.Model(new(T)).Create(data).Error
}

func CreateBatch[T any](db *gorm.DB, data []*T) error {
	if len(data) == 0 {
		return nil
	}
	return db.Model(new(T)).Create(&data).Error
}

func Delete[T any](db *gorm.DB, conditions map[string]any, hardDelete ...bool) error {
	sql := applyConditions[T](db, conditions)

	// 判断是否硬删除
	if len(hardDelete) > 0 && hardDelete[0] {
		return sql.Unscoped().Delete(new(T)).Error
	}
	return sql.Delete(new(T)).Error
}

func List[T any](db *gorm.DB, currentPage, pageSize int, conditions map[string]any, orderBy []string) ([]*T, int64, error) {
	var list []*T
	var total int64

	sql := applyConditions[T](db, conditions)

	// 获取总数
	if err := sql.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 支持排序
	if len(orderBy) > 0 {
		for _, order := range orderBy {
			sql = sql.Order(order)
		}
	} else {
		sql = sql.Order("created_at DESC")
	}

	// 支持分页
	if pageSize > 0 && currentPage > 0 {
		sql = sql.Offset(pageSize * (currentPage - 1)).Limit(pageSize)
	}

	if err := sql.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func Get[T any](db *gorm.DB, conditions map[string]any) (*T, error) {
	var data T

	sql := applyConditions[T](db, conditions)

	err := sql.Take(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func Update[T any](db *gorm.DB, data *T, conditions map[string]any, fields ...string) error {
	sql := applyConditions[T](db, conditions)

	// 指定更新字段
	if len(fields) > 0 {
		sql = sql.Select(fields)
	}

	return sql.Updates(data).Error
}

func applyConditions[T any](db *gorm.DB, conditions map[string]any) *gorm.DB {
	sql := db.Model(new(T))
	for k, v := range conditions {
		switch val := v.(type) {
		case []any:
			if len(val) == 2 {
				sql = sql.Where(fmt.Sprintf("%s %s ?", k, val[0]), val[1])
			}
		default:
			sql = sql.Where(fmt.Sprintf("%s = ?", k), v)
		}
	}
	return sql
}
