package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"

	schemaModel "mvtable/internal/app/mv_table_schema/model"

	"gorm.io/datatypes"
)

type MvField struct {
	db.Model
	TableSchemaID string              `gorm:"column:table_schema_id"`
	Title         string              `gorm:"column:title"`
	Type          constants.FieldType `gorm:"column:type"`
	Config        datatypes.JSON      `gorm:"column:config;default:'{}'"`
	OrderIndex    string              `gorm:"column:order_index;default:'mmmmmmmmmm'"`
	IsPrimaryKey  bool                `gorm:"column:is_primary_key;default:false"`
}

func (MvField) TableName() string {
	return "mv_field"
}

type GetMvFieldListReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
}

type GetMvFieldListRes struct {
	List []schemaModel.GeneralFieldConfig `json:"fields"`
}

type CreateMvFieldReq struct {
	TableSchemaID string              `json:"tableSchemaId" binding:"required"`
	Title         string              `json:"title" binding:"required"`
	Type          constants.FieldType `json:"type" binding:"required"`
	Config        any                 `json:"config"`
}

type UpdateMvFieldReq struct {
	ID     string              `json:"id" binding:"required"`
	Title  string              `json:"title" binding:"required"`
	Type   constants.FieldType `json:"type" binding:"required"`
	Config any                 `json:"config" binding:"required"`
}

type DeleteMvFieldReq struct {
	ID string `json:"id" binding:"required"`
}

type SortMvFieldReq struct {
	TableSchemaID string   `json:"tableSchemaId" binding:"required"`
	FieldIDs      []string `json:"fieldIds" binding:"required,min=1"`
}
