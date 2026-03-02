package model

import (
	"mvtable/internal/storage/db"

	"github.com/lib/pq"
)

type MvTemplateProject struct {
	db.Model
	UserID      string         `gorm:"column:user_id" json:"userId"`
	Name        string         `gorm:"column:name;default:''" json:"name"`
	Description string         `gorm:"column:description;default:''" json:"description"`
	Cover       string         `gorm:"column:cover;default:''" json:"cover"`
	Tags        pq.StringArray `gorm:"type:text[];column:tags;default:'{}'::text[]" json:"tags"`
	UseCount    int            `gorm:"column:use_count;default:0" json:"useCount"`
}

func (MvTemplateProject) TableName() string {
	return "mv_template_project"
}

type QueryMvTemplateProjectReq struct {
	db.Pagination
	Tag string `form:"tag"`
}

type TempQueryItem struct {
	ID          string         `gorm:"column:id"`
	Creator     string         `gorm:"column:creator"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Cover       string         `gorm:"column:cover"`
	Tags        pq.StringArray `gorm:"type:text[];column:tags"`
	UseCount    int            `gorm:"column:use_count"`
}

type QueryMvTemplateProjectItem struct {
	ID          string   `json:"id"`
	Creator     string   `json:"creator"` // 创建者的昵称
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
	UseCount    int      `json:"useCount"`
}

type QueryMvTemplateProjectRes struct {
	List  []*QueryMvTemplateProjectItem `json:"list"`
	Total int64                         `json:"total"`
}

type SaveAsTemplateReq struct {
	ProjectID   string   `json:"projectId" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Cover       string   `json:"cover"`
	Tags        []string `json:"tags"`
}

type SaveAsTemplateRes struct {
	ID string `json:"id"`
}

type UpdateMvTemplateProjectReq struct {
	ID          string    `json:"id" binding:"required"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Cover       *string   `json:"cover"`
	Tags        *[]string `json:"tags"`
}

type DeleteMvTemplateProjectReq struct {
	ID string `json:"id" binding:"required"`
}
