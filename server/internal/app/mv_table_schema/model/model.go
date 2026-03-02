package model

import (
	"mvtable/internal/storage/db"
	"time"

	"gorm.io/datatypes"
)

type MvTableSchema struct {
	db.Model
	FolderID    string         `gorm:"column:folder_id;default:''" json:"folderId"`
	ProjectID   string         `gorm:"column:project_id;default:''" json:"projectId"`
	Name        string         `gorm:"column:name" json:"name"`
	Version     int            `gorm:"column:version;default:1" json:"version"`
	CreatedBy   string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedBy   *string        `gorm:"column:updated_by" json:"updatedBy"`
	Description string         `gorm:"column:description;default:''" json:"description"`
	Config      datatypes.JSON `gorm:"column:config;default:'{}'" json:"config"`
	Stats       datatypes.JSON `gorm:"column:stats;default:'{}'" json:"stats"`
	RowName     string         `gorm:"column:row_name;default:''" json:"rowName"`
}

func (MvTableSchema) TableName() string {
	return "mv_table_schema"
}

type CreateMvTableSchemaReq struct {
	ProjectID string `json:"projectId"`
	FolderID  string `json:"folderId"`
	Name      string `json:"name" binding:"required"`
}

type UpdateMvTableSchemaReq struct {
	ID          string  `json:"id" binding:"required"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	RowName     *string `json:"rowName"`
}

type DeleteMvTableSchemaReq struct {
	ID string `json:"id" binding:"required"`
}

// CopyMvTableSchemaReq 复制数据表请求
// range: 1=仅复制数据表结构，2=复制数据表结构和所有记录
type CopyMvTableSchemaReq struct {
	ProjectID string `json:"projectId" binding:"required"` // 项目ID
	SourceID  string `json:"sourceId" binding:"required"`  // 源数据表ID
	Name      string `json:"name" binding:"required"`      // 新数据表名称（项目中不允许重复）
	Range     int    `json:"range" binding:"required"`     // 1=仅结构 2=结构+记录
}

type CopyMvTableSchemaRes struct {
	ID string `json:"id"` // 新创建的数据表ID
}

type GetTableDataReq struct {
	ViewID string `form:"viewId" binding:"required"`
}

type RecordItem map[string]any

type WatermarkConfig struct {
	Content WatermarkContent `json:"content"`
	Style   WatermarkStyle   `json:"style"`
}

type WatermarkContent struct {
	Type     string  `json:"type"` // "text" or "image"
	Text     *string `json:"text,omitempty"`
	ImageUrl *string `json:"imageUrl,omitempty"`
}

type WatermarkStyle struct {
	Opacity  float64 `json:"opacity"`
	Rotation float64 `json:"rotation"`
	FontSize int     `json:"fontSize"`
	Color    string  `json:"color"`
	Spacing  int     `json:"spacing"`
}

type TableConfig struct {
	EnableWatermark bool             `json:"enableWatermark"`
	WatermarkConfig *WatermarkConfig `json:"watermarkConfig"`
	EnableRevision  bool             `json:"enableRevision"`
	MaxRecords      int              `json:"maxRecords"`
	AllowExport     bool             `json:"allowExport"`
	AllowImport     bool             `json:"allowImport"`
}

type GetTableDataRes struct {
	TableSchemaID string               `json:"tableSchemaId"`
	Name          string               `json:"name"`
	Version       int                  `json:"version"`
	Ct            time.Time            `json:"ct"`
	Ut            time.Time            `json:"ut"`
	CreatedBy     string               `json:"createdBy"`
	UpdatedBy     string               `json:"updatedBy"`
	Config        TableConfig          `json:"config"`
	Fields        []GeneralFieldConfig `json:"fields"`
	Records       []RecordItem         `json:"records"`
}

// 筛选条件配置
type FilterCondition struct {
	FieldID  string      `json:"fieldId"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// 排序配置
type SortCondition struct {
	FieldID string `json:"fieldId"`
	Order   string `json:"order"` // "asc" 或 "desc"
}

// 分组配置
type GroupCondition struct {
	FieldID string `json:"fieldId"`
	Order   string `json:"order"` // "asc" 或 "desc"
}
