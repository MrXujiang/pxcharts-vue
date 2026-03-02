package model

import (
	"mvtable/internal/storage/db"
)

type MvView struct {
	db.Model
	TableSchemaID string `gorm:"column:table_schema_id"`
	Type          string `gorm:"column:type;default:'table'"`
	Name          string `gorm:"column:name;default:''"`
	Description   string `gorm:"column:description;default:''"`
	OrderIndex    string `gorm:"column:order_index;default:'mmmmmmmmmm'"`
}

func (MvView) TableName() string {
	return "mv_view"
}

type QueryMvViewReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
	ProjectID     string `form:"projectId" binding:"required"`
}

type MvViewItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	IsDefault   bool   `json:"isDefault"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type QueryMvViewRes struct {
	List []MvViewItem `json:"list"`
}

type CreateMvViewReq struct {
	TableSchemaID string `json:"tableSchemaId" binding:"required"`
	Type          string `json:"type" binding:"required"` // table=表格视图 form=表单视图 board=看板视图
}

type UpdateMvViewReq struct {
	ID          string  `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type DeleteMvViewReq struct {
	ID string `json:"id" binding:"required"`
}

type GetMvViewReq struct {
	ID string `form:"id" binding:"required"`
}

type SwitchActiveViewReq struct {
	TableSchemaID string `json:"tableSchemaId" binding:"required"`
	ViewID        string `json:"viewId" binding:"required"`
}

type GetMvViewRes struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	TableConfig *MvViewTableConfig `json:"tableConfig,omitempty"`
	FormConfig  *MvViewFormConfig  `json:"formConfig,omitempty"`
	BoardConfig *MvViewBoardConfig `json:"boardConfig,omitempty"`
}

type MvViewTableConfig struct {
	FilterConfig any `json:"filterConfig"`
	GroupConfig  any `json:"groupConfig"`
	SortConfig   any `json:"sortConfig"`
	RowHeight    int `json:"rowHeight"`
	ColorConfig  any `json:"colorConfig"`
}

type MvViewFormConfig struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	Cover                 string `json:"cover"`
	Layout                string `json:"layout"`
	Stats                 string `json:"stats"`
	EnableSharing         bool   `json:"enableSharing"`
	EnableAnonymous       bool   `json:"enableAnonymous"`
	Filter                string `json:"filter"`
	FilterConfig          any    `json:"filterConfig"`
	EnableNoLogin         bool   `json:"enableNoLogin"`
	EnableLimitSubmit     bool   `json:"enableLimitSubmit"`
	LimitSubmitType       string `json:"limitSubmitType"`
	EnableLimitCollect    bool   `json:"enableLimitCollect"`
	LimitCollectCount     int64  `json:"limitCollectCount"`
	EnableCycleRemind     bool   `json:"enableCycleRemind"`
	CycleRemindConfig     any    `json:"cycleRemindConfig"`
	EnableEditAfterSubmit bool   `json:"enableEditAfterSubmit"`
	Config                any    `json:"config"`
}

type MvViewBoardConfig struct {
	ShowFieldTitle bool `json:"showFieldTitle"`
	FilterConfig   any  `json:"filterConfig"`
	GroupConfig    any  `json:"groupConfig"`
	SortConfig     any  `json:"sortConfig"`
}
