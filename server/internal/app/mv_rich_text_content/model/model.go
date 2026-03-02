package model

import (
	"mvtable/internal/storage/db"
)

type MvRichTextContent struct {
	db.Model
	RecordID string `gorm:"column:record_id;not null"`
	FieldID  string `gorm:"column:field_id;not null"`
	Content  string `gorm:"column:content;type:text;not null;default:''"`
}

func (MvRichTextContent) TableName() string {
	return "mv_rich_text_content"
}

type GetRichTextContentReq struct {
	RecordID string `form:"recordId" binding:"required"`
	FieldID  string `form:"fieldId" binding:"required"`
}

type GetRichTextContentRes struct {
	Content string `json:"content"`
}
