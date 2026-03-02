package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvViewForm struct {
	db.Model
	ViewID                string         `gorm:"column:view_id" json:"viewId"`
	Name                  string         `gorm:"column:name" json:"name"`
	Description           string         `gorm:"column:description" json:"description"`
	Cover                 string         `gorm:"column:cover" json:"cover"`
	Layout                string         `gorm:"column:layout" json:"layout"`
	Stats                 datatypes.JSON `gorm:"column:stats" json:"stats"`
	EnableSharing         bool           `gorm:"column:enable_sharing" json:"enableSharing"`
	EnableAnonymous       bool           `gorm:"column:enable_anonymous" json:"enableAnonymous"`
	Filter                string         `gorm:"column:filter" json:"filter"`
	FilterConfig          datatypes.JSON `gorm:"column:filter_config" json:"filterConfig"`
	EnableNoLogin         bool           `gorm:"column:enable_no_login" json:"enableNoLogin"`
	EnableLimitSubmit     bool           `gorm:"column:enable_limit_submit" json:"enableLimitSubmit"`
	LimitSubmitType       string         `gorm:"column:limit_submit_type" json:"limitSubmitType"`
	EnableLimitCollect    bool           `gorm:"column:enable_limit_collect" json:"enableLimitCollect"`
	LimitCollectCount     int64          `gorm:"column:limit_collect_count" json:"limitCollectCount"`
	EnableCycleRemind     bool           `gorm:"column:enable_cycle_remind" json:"enableCycleRemind"`
	CycleRemindConfig     datatypes.JSON `gorm:"column:cycle_remind_config" json:"cycleRemindConfig"`
	EnableEditAfterSubmit bool           `gorm:"column:enable_edit_after_submit" json:"enableEditAfterSubmit"`
	Config                datatypes.JSON `gorm:"column:config" json:"config"`
}

func (MvViewForm) TableName() string {
	return "mv_view_form"
}

type UpdateMvViewFormReq struct {
	ID                    string  `json:"id" binding:"required"`
	Name                  *string `json:"name"`
	Description           *string `json:"description"`
	Cover                 *string `json:"cover"`
	Layout                *string `json:"layout"`
	Stats                 *string `json:"stats"`
	EnableSharing         *bool   `json:"enableSharing"`
	EnableAnonymous       *bool   `json:"enableAnonymous"`
	Filter                *string `json:"filter"`
	FilterConfig          *[]byte `json:"filterConfig"`
	EnableNoLogin         *bool   `json:"enableNoLogin"`
	EnableLimitSubmit     *bool   `json:"enableLimitSubmit"`
	LimitSubmitType       *string `json:"limitSubmitType"`
	EnableLimitCollect    *bool   `json:"enableLimitCollect"`
	LimitCollectCount     *int64  `json:"limitCollectCount"`
	EnableEditAfterSubmit *bool   `json:"enableEditAfterSubmit"`
	Config                *[]byte `json:"config"`
}
