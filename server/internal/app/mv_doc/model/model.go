package model

import "mvtable/internal/storage/db"

type MvDoc struct {
	db.Model
	ProjectID string `gorm:"column:project_id;default:''" json:"projectId"`
	FolderID  string `gorm:"column:folder_id;default:''" json:"folderId"`
	Name      string `gorm:"column:name;default:''" json:"name"`
	Content   string `gorm:"column:content;default:''" json:"content"`
	CreatedBy string `gorm:"column:created_by;default:''" json:"createdBy"`
	UpdatedBy string `gorm:"column:updated_by;default:''" json:"updatedBy"`
}

func (MvDoc) TableName() string {
	return "mv_doc"
}

// UserInfo 用户信息结构体
type UserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// CreateMvDocReq 创建文档请求
type CreateMvDocReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	FolderID  string `json:"folderId"`
	Name      string `json:"name" binding:"required"`
	Content   string `json:"content"`
}

// UpdateMvDocReq 更新文档请求
type UpdateMvDocReq struct {
	ID       string  `json:"id" binding:"required"`
	Name     *string `json:"name"`
	Content  *string `json:"content"`
	FolderID *string `json:"folderId"`
}

// DeleteMvDocReq 删除文档请求
type DeleteMvDocReq struct {
	ID string `json:"id" binding:"required"`
}

// GetMvDocReq 获取文档请求
type GetMvDocReq struct {
	ID string `form:"id" binding:"required"`
}

// GetMvDocRes 获取文档响应
type GetMvDocRes struct {
	MvDocItem
}

// ListMvDocReq 列表查询文档请求
type ListMvDocReq struct {
	ProjectID string `form:"projectId"`
	FolderID  string `form:"folderId"`
	Keyword   string `form:"keyword"`
	db.Pagination
}

// MvDocItem 文档列表项
type MvDocItem struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	FolderID  string    `json:"folderId"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Creator   *UserInfo `json:"creator"`
	Updater   *UserInfo `json:"updater"`
	Ct        string    `json:"ct"`
	Ut        string    `json:"ut"`
}

// ListMvDocRes 列表查询文档响应
type ListMvDocRes struct {
	List  []MvDocItem `json:"list"`
	Total int64       `json:"total"`
}
