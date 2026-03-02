package model

import (
	schemaModel "mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvRecord struct {
	db.Model
	TableSchemaID string         `gorm:"column:table_schema_id"`
	CreatedBy     string         `gorm:"column:created_by"`
	UpdatedBy     string         `gorm:"column:updated_by;default:null"`
	RowData       datatypes.JSON `gorm:"column:row_data;default:'{}'"`
	OrderIndex    string         `gorm:"column:order_index;default:'mmmmmmmmmm'"`
}

func (MvRecord) TableName() string {
	return "mv_record"
}

type InsertMvRecordReq struct {
	PrevID        string          `json:"prevId"`
	TableSchemaID string          `json:"tableSchemaId" binding:"required"`
	RowData       *datatypes.JSON `json:"rowData"`
}

type SaveRowSortReq struct {
	PrevRecordID    string `json:"prevRecordId"`
	NextRecordID    string `json:"nextRecordId"`
	CurrentRecordID string `json:"currentId"`
}

type UpdateCellReq struct {
	RecordID string `json:"recordId" binding:"required"`
	FieldID  string `json:"fieldId" binding:"required"`
	Value    any    `json:"value"`
}

type UpdateRowReq map[string]any

type DeleteRecordsReq struct {
	RecordIDs []string `json:"recordIds" binding:"required,min=1"`
}

type ImportRecordsReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
}

type GetRecordReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
	RecordID      string `form:"recordId" binding:"required"`
}

type GetRecordsReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
}

type GetRecordsRes struct {
	Fields  []schemaModel.GeneralFieldConfig `json:"fields"`
	Records []any                            `json:"records"`
}

type ImportRecordsRes struct {
	TotalCount   int           `json:"totalCount"`   // 总行数
	SuccessCount int           `json:"successCount"` // 成功导入数量
	FailCount    int           `json:"failCount"`    // 失败数量
	Errors       []ImportError `json:"errors"`       // 错误详情
}

type ImportError struct {
	Row    int    `json:"row"`    // 行号（从1开始，包含标题行）
	Column string `json:"column"` // 列名
	Error  string `json:"error"`  // 错误信息
}

type GetRecordOptionsReq struct {
	TableSchemaID string `form:"tableSchemaId" binding:"required"`
}

type GetRecordOptionsRes struct {
	Options []RecordOption `json:"options"`
}

type RecordOption struct {
	ID    string `json:"id"`    // rowId
	Value any    `json:"value"` // 第一个字段对应的数据
}
