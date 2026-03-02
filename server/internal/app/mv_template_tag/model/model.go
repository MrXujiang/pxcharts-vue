package model

import (
	"mvtable/internal/storage/db"
)

type MvTemplateTag struct {
	db.Model
	Name        string `gorm:"column:name;default:''" json:"name"`
	Description string `gorm:"column:description;default:''" json:"description"`
	OrderIndex  string `gorm:"column:order_index;default:'mmmmmmmmmm'" json:"orderIndex"`
}

func (MvTemplateTag) TableName() string {
	return "mv_template_tag"
}

type CreateMvTemplateTagReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateMvTemplateTagReq struct {
	ID          string  `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeleteMvTemplateTagReq struct {
	ID string `json:"id" binding:"required"`
}

type UpdateTagSortReq struct {
	PrevTagID    string `json:"prevTagId"`
	NextTagID    string `json:"nextTagId"`
	CurrentTagID string `json:"currentTagId" binding:"required"`
}

type QueryMvTemplateTagReq struct {
	SearchWord string `form:"searchWord"`
}

type QueryMvTemplateTagItem struct {
	ID          string `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}

type QueryMvTemplateTagRes struct {
	List []*QueryMvTemplateTagItem `json:"list"`
}
