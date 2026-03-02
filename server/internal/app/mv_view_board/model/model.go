package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvViewBoard struct {
	db.Model
	ViewID         string         `gorm:"column:view_id" json:"viewId"`
	ShowFieldTitle bool           `gorm:"column:show_field_title;default:false" json:"showFieldTitle"`
	FilterConfig   datatypes.JSON `gorm:"column:filter_config;default:'[]'" json:"filterConfig"`
	GroupConfig    datatypes.JSON `gorm:"column:group_config;default:'{}'" json:"groupConfig"`
	SortConfig     datatypes.JSON `gorm:"column:sort_config;default:'[]'" json:"sortConfig"`
}

func (MvViewBoard) TableName() string {
	return "mv_view_board"
}

type UpdateMvViewBoardReq struct {
	ViewID string `json:"viewId" binding:"required"`
	// ShowFieldTitle 是否显示字段标题
	ShowFieldTitle *bool `json:"showFieldTitle"`
	// FilterConfig 筛选条件: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
	FilterConfig *datatypes.JSON `json:"filterConfig"`
	// GroupConfig 分组配置: {fieldId: ""}
	GroupConfig *datatypes.JSON `json:"groupConfig"`
	// SortConfig 排序配置: [{fieldId: "", order: "asc/desc"}]
	SortConfig *datatypes.JSON `json:"sortConfig"`
}
