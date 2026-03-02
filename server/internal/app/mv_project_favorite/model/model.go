package model

import "mvtable/internal/storage/db"

type MvProjectFavorite struct {
	db.Model
	ProjectID string `gorm:"column:project_id" json:"projectId"`
	UserID    string `gorm:"column:user_id" json:"userId"`
}

func (MvProjectFavorite) TableName() string {
	return "mv_project_favorite"
}
