package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvFormSubmit struct {
	db.Model
	UserID     string `gorm:"column:user_id" json:"userId"`
	ViewFormID string `gorm:"column:view_form_id" json:"viewFormId"`
	RecordID   string `gorm:"column:record_id" json:"recordId"`
}

func (m *MvFormSubmit) TableName() string {
	return "mv_form_submit"
}

type SubmitFormReq struct {
	FormID   string         `json:"formId" binding:"required"`
	FormData datatypes.JSON `json:"formData"`
}
