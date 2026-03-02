package model

import (
	"mvtable/internal/storage/db"
	"time"
)

type MvProjectRecent struct {
	db.Model
	ProjectID      string    `gorm:"column:project_id" json:"projectId"`
	UserID         string    `gorm:"column:user_id" json:"userId"`
	LastAccessedAt time.Time `gorm:"column:last_accessed_at;default:CURRENT_TIMESTAMP" json:"lastAccessedAt"`
}

func (MvProjectRecent) TableName() string {
	return "mv_project_recent"
}
