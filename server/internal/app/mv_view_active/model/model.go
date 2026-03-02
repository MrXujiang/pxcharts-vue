package model

import "mvtable/internal/storage/db"

type MvViewActive struct {
	db.Model
	TableSchemaID string `gorm:"column:table_schema_id"`
	ViewID        string `gorm:"column:view_id"`
	UserID        string `gorm:"column:user_id"`
}

func (MvViewActive) TableName() string {
	return "mv_view_active"
}
