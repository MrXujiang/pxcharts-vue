package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvViewTable struct {
	db.Model
	ViewID       string         `gorm:"column:view_id" json:"viewId"`
	FilterConfig datatypes.JSON `gorm:"column:filter_config;default:'[]'" json:"filterConfig"`
	GroupConfig  datatypes.JSON `gorm:"column:group_config;default:'{}'" json:"groupConfig"`
	SortConfig   datatypes.JSON `gorm:"column:sort_config;default:'[]'" json:"sortConfig"`
	RowHeight    int            `gorm:"column:row_height;default:1" json:"rowHeight"`
	ColorConfig  datatypes.JSON `gorm:"column:color_config;default:'[]'" json:"colorConfig"`
}

func (MvViewTable) TableName() string {
	return "mv_view_table"
}

type UpdateMvViewTableReq struct {
	ViewID string `json:"viewId" binding:"required"`
	// FilterConfig 筛选条件: [{fieldId: "字段id", operator: "equal=等于/notEqual=不等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
	FilterConfig *datatypes.JSON `json:"filterConfig"`
	// GroupConfig 分组配置: {fieldId: "", order: "asc/desc"}
	GroupConfig *datatypes.JSON `json:"groupConfig"`
	// SortConfig 排序配置: [{fieldId: "", order: "asc/desc"}]
	SortConfig *datatypes.JSON `json:"sortConfig"`
	// RowHeight 行高: 1=常规, 2=中等, 3=高, 4=超高
	RowHeight *int `json:"rowHeight"`
	// ColorConfig 颜色配置: [{color: "#cccccc", scope: "cell/row/rowHeader/column", fieldId: "字段id", operator: "equal=等于/notEqual=不等于/greaterThan=大于/lessThan=小于/greaterThanOrEqual=大于等于/lessThanOrEqual=小于等于/contains=包含/notContains=不包含/null=为空/notNull=不为空/before=早于/after=晚于/beforeOrEqual=早于或等于/afterOrEqual=晚于或等于", value: "值"}]
	ColorConfig *datatypes.JSON `json:"colorConfig"`
}
