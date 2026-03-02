package model

import "mvtable/internal/storage/db"

type MvDashboard struct {
	db.Model
	ProjectID string `gorm:"column:project_id;default:''" json:"projectId"`
	FolderID  string `gorm:"column:folder_id;default:''" json:"folderId"`
	Name      string `gorm:"column:name;default:''" json:"name"`
	BgType    string `gorm:"column:bg_type;default:''" json:"bgType"`
	BgColor   string `gorm:"column:bg_color;default:''" json:"bgColor"`
	BgImage   string `gorm:"column:bg_image;default:''" json:"bgImage"`
}

func (MvDashboard) TableName() string {
	return "mv_dashboard"
}

// CreateMvDashboardReq 创建仪表盘请求
type CreateMvDashboardReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	FolderID  string `json:"folderId"`
	Name      string `json:"name" binding:"required"`
	BgType    string `json:"bgType"`
	BgColor   string `json:"bgColor"`
	BgImage   string `json:"bgImage"`
}

// UpdateMvDashboardReq 更新仪表盘请求
type UpdateMvDashboardReq struct {
	ID       string  `json:"id" binding:"required"`
	Name     *string `json:"name"`
	BgType   *string `json:"bgType"`
	BgColor  *string `json:"bgColor"`
	BgImage  *string `json:"bgImage"`
	FolderID *string `json:"folderId"`
}

// DeleteMvDashboardReq 删除仪表盘请求
type DeleteMvDashboardReq struct {
	ID string `json:"id" binding:"required"`
}

// GetDashboardReq 获取仪表盘请求
type GetDashboardReq struct {
	ID string `form:"id" binding:"required"`
}

// DashboardChartItem 仪表盘下图表信息
type DashboardChartItem struct {
	ID            string `json:"id"`
	TableSchemaID string `json:"tableSchemaId"`
	Field1ID      string `json:"field1Id"`
	Field2ID      string `json:"field2Id"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	Config        any    `json:"config"`
}

// GetDashboardRes 获取仪表盘响应
type GetDashboardRes struct {
	MvDashboard
	Charts []DashboardChartItem `json:"charts"`
}

// CopyMvDashboardReq 复制仪表盘请求
type CopyMvDashboardReq struct {
	ID   string  `json:"id" binding:"required"` // 要复制的仪表盘ID
	Name *string `json:"name"`                  // 新仪表盘名称，为空则自动生成"原名称_副本"
}

// CopyMvDashboardRes 复制仪表盘响应
type CopyMvDashboardRes struct {
	ID string `json:"id"` // 新创建的仪表盘ID
}
