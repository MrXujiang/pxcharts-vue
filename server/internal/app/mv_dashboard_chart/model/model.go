package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type MvDashboardChart struct {
	db.Model
	DashboardID   string         `gorm:"column:dashboard_id;default:''" json:"dashboardId"`
	TableSchemaID string         `gorm:"column:table_schema_id;default:''" json:"tableSchemaId"`
	Field1ID      string         `gorm:"column:field1_id;default:''" json:"field1Id"`
	Field2ID      string         `gorm:"column:field2_id;default:''" json:"field2Id"`
	Title         string         `gorm:"column:title;default:''" json:"title"`
	Type          string         `gorm:"column:type;default:''" json:"type"`
	Config        datatypes.JSON `gorm:"column:config;default:'{}'" json:"config"`
}

func (MvDashboardChart) TableName() string {
	return "mv_dashboard_chart"
}

// CreateMvDashboardChartReq 创建仪表盘图表请求
type CreateMvDashboardChartReq struct {
	DashboardID   string         `json:"dashboardId" binding:"required"`
	TableSchemaID string         `json:"tableSchemaId" binding:"required"`
	Field1ID      string         `json:"field1Id"`
	Field2ID      string         `json:"field2Id"`
	Title         string         `json:"title" binding:"required"`
	Type          string         `json:"type" binding:"required"`
	Config        map[string]any `json:"config"`
}

// UpdateMvDashboardChartReq 更新仪表盘图表请求
type UpdateMvDashboardChartReq struct {
	ID            string          `json:"id" binding:"required"`
	DashboardID   *string         `json:"dashboardId"`
	TableSchemaID *string         `json:"tableSchemaId"`
	Field1ID      *string         `json:"field1Id"`
	Field2ID      *string         `json:"field2Id"`
	Title         *string         `json:"title"`
	Type          *string         `json:"type"`
	Config        *datatypes.JSON `json:"config"`
}

// DeleteMvDashboardChartReq 删除仪表盘图表请求
type DeleteMvDashboardChartReq struct {
	ID string `json:"id" binding:"required"`
}
