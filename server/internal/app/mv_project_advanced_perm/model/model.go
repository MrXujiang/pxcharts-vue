package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvProjectAdvancedPerm struct {
	db.Model
	TableSchemaID   string                         `gorm:"column:table_schema_id" json:"tableSchemaId"`
	Role            constants.ProjectRole          `gorm:"column:role" json:"role"`
	DataAction      string                         `gorm:"column:data_action;default:'none'" json:"dataAction"`
	CanAdd          bool                           `gorm:"column:can_add;default:false" json:"canAdd"`
	CanDelete       bool                           `gorm:"column:can_delete;default:false" json:"canDelete"`
	OperateRange    string                         `gorm:"column:operate_range;default:'all'" json:"operateRange"`
	FieldAccess     string                         `gorm:"column:field_access;default:'all'" json:"fieldAccess"`
	CustomFieldPerm datatypes.JSONSlice[FieldPerm] `gorm:"column:custom_field_perm;default:'[]'" json:"customFieldPerm"`
	CanOperateView  bool                           `gorm:"column:can_operate_view;default:false" json:"canOperateView"`
	ViewAccess      string                         `gorm:"column:view_access;default:'all'" json:"viewAccess"`
	CanCheckViews   datatypes.JSONSlice[string]    `gorm:"column:can_check_views;default:'[]'" json:"canCheckViews"`
}

func (MvProjectAdvancedPerm) TableName() string {
	return "mv_project_advanced_perm"
}

type EnableMvProjectAdvancedPermReq struct {
	ProjectID string `json:"projectId"`
}

type DisableMvProjectAdvancedPermReq struct {
	ProjectID string `json:"projectId"`
}

type FieldPerm struct {
	FieldID string `json:"fieldId"`
	CanRead bool   `json:"canRead"`
	CanAdd  bool   `json:"canAdd"`
	CanEdit bool   `json:"canEdit"`
}

type UpdateMvProjectAdvancedPermReq struct {
	ID              string       `json:"id"`
	DataAction      string       `json:"dataAction" binding:"required"`
	CanAdd          *bool        `json:"canAdd"`
	CanDelete       *bool        `json:"canDelete"`
	OperateRange    *string      `json:"operateRange"`
	FieldAccess     *string      `json:"fieldAccess"`
	CustomFieldPerm *[]FieldPerm `json:"customFieldPerm"`
	CanOperateView  *bool        `json:"canOperateView"`
	ViewAccess      *string      `json:"viewAccess"`
	CanCheckViews   *[]string    `json:"canCheckViews"`
}

type GetMvProjectAdvancedPermReq struct {
	Role          string `form:"role" binding:"required"`
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
}

type GetMvProjectAdvancedPermRes struct {
	ID              string      `json:"id"`
	DataAction      string      `json:"dataAction"`
	CanAdd          bool        `json:"canAdd"`
	CanDelete       bool        `json:"canDelete"`
	OperateRange    string      `json:"operateRange"`
	FieldAccess     string      `json:"fieldAccess"`
	CustomFieldPerm []FieldPerm `json:"customFieldPerm"`
	ViewAccess      string      `json:"viewAccess"`
	CanCheckViews   []string    `json:"canCheckViews"`
}

func (a *GetMvProjectAdvancedPermRes) FillFrom(data *MvProjectAdvancedPerm) {
	a.ID = data.ID
	a.DataAction = data.DataAction
	a.CanAdd = data.CanAdd
	a.CanDelete = data.CanDelete
	a.OperateRange = data.OperateRange
	a.FieldAccess = data.FieldAccess
	a.CustomFieldPerm = data.CustomFieldPerm
	a.ViewAccess = data.ViewAccess
	a.CanCheckViews = data.CanCheckViews
}
