package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"time"
)

type MvProject struct {
	db.Model
	Name               string `gorm:"column:name" json:"name"`
	Description        string `gorm:"column:description;default:''" json:"description"`
	UserID             string `gorm:"column:user_id" json:"userId"`
	EnableAdvancedPerm bool   `gorm:"column:enable_advanced_perm;default:false" json:"enableAdvancedPerm"`
}

func (MvProject) TableName() string {
	return "mv_project"
}

type CreateMvProjectReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateMvProjectRes struct {
	ID string `json:"id"`
}

type UpdateMvProjectReq struct {
	ID                 string  `json:"id" binding:"required"`
	Name               *string `json:"name"`
	Description        *string `json:"description"`
	EnableAdvancedPerm *bool   `json:"enableAdvancedPerm"`
}

type DeleteMvProjectReq struct {
	ID string `json:"id" binding:"required"`
}

type SetFavoriteProjectReq struct {
	ProjectID  string `json:"projectId" binding:"required"`
	IsFavorite bool   `json:"isFavorite"`
}

type GetProjectReq struct {
	ProjectID string `form:"projectId" binding:"required"`
}

type GetProjectRes struct {
	MvProject
	IsFavorite bool `json:"isFavorite"`
}

type QueryProjectReq struct {
	db.Pagination
	// 查询项目类型: 1=最近访问 2=我创建的 3=我参与的 4=我收藏的
	Type constants.QueryProjectType `form:"type" binding:"required"`
}

type QueryProjectItem struct {
	ID          string    `gorm:"column:id" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	IsFavorite  bool      `gorm:"column:is_favorite" json:"isFavorite"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type QueryProjectRes struct {
	List  []*QueryProjectItem `json:"list"`
	Total int64               `json:"total"`
}

// RenameProjectNodeReq
// 在项目下重命名节点（文件夹 / 表格 / 表单 / 仪表盘）
// type 可选: table（表格）、form（表单视图）、dashboard（仪表盘）、folder（文件夹）
type RenameProjectNodeReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	TargetID  string `json:"targetId" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required"` // table / form / dashboard / folder
}

// CreateProjectNodeReq
// 在项目下创建节点（文件夹 / 表格 / 表单 / 仪表盘）
// type 可选: table（创建数据表及表格视图）、form（创建数据表及表单视图）、dashboard（创建仪表盘）、folder（创建文件夹）
type CreateProjectNodeReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	FolderID  string `json:"folderId"`                // 所在文件夹ID，根目录可为空
	Name      string `json:"name" binding:"required"` // 名称
	Type      string `json:"type" binding:"required"` // table / form / dashboard / folder
}

type CreateProjectNodeRes struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// 获取项目下数据表列表的请求和响应
type GetProjectTablesReq struct {
	ProjectID string `form:"projectId" binding:"required"`
}

type GetProjectTablesRes struct {
	List  []*TableSchemaItem `json:"list"`
	Total int64              `json:"total"`
}

type TableSchemaItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FolderID    string `json:"folderId"`
}
