package model

import (
	"mvtable/internal/storage/db"
)

type MvFolder struct {
	db.Model
	ProjectID string `gorm:"column:project_id" json:"projectId"`
	Name      string `gorm:"column:name" json:"name"`
	ParentID  string `gorm:"column:parent_id;default:''" json:"parentId"`
}

func (MvFolder) TableName() string {
	return "mv_folder"
}

type CreateMvFolderReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ParentID  string `json:"parentId"`
}

type UpdateMvFolderReq struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DeleteMvFolderReq struct {
	ID string `json:"id" binding:"required"`
}

type SubNode struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	ParentID string     `json:"parentId"`
	Children []*SubNode `json:"children"`
}

type QuerySubNodeReq struct {
	ParentID  string `form:"parentId"`
	ProjectID string `form:"projectId"`
}

type QuerySubNodeRes struct {
	List []*SubNode `json:"list"`
}

// QueryAllNodeReq 查询项目下所有节点（树状）
type QueryAllNodeReq struct {
	ProjectID string `form:"projectId" binding:"required"`
}

// QueryFolderListReq 查询项目下所有文件夹（扁平列表）
type QueryFolderListReq struct {
	ProjectID string `form:"projectId" binding:"required"`
}

// FolderItem 扁平文件夹节点
type FolderItem struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	ParentID  string `json:"parentId"`
}

// QueryFolderListRes 查询项目下所有文件夹响应（扁平结构）
type QueryFolderListRes struct {
	List []*FolderItem `json:"list"`
}

type SearchReq struct {
	ProjectID string `form:"projectId" binding:"required"`
	Keyword   string `form:"keyword" binding:"required"`
}

type SearchItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	FolderID string `json:"folderId"`
}

type SearchRes struct {
	List []*SubNode `json:"list"`
}

// MoveNodeReq 移动节点请求
type MoveNodeReq struct {
	ProjectID  string `json:"projectId" binding:"required"`  // 项目ID
	TargetID   string `json:"targetId" binding:"required"`   // 要移动的节点ID
	Type       string `json:"type" binding:"required"`        // 节点类型: folder / table / form / dashboard
	TargetFolderID string `json:"targetFolderId"`             // 目标文件夹ID，为空字符串表示移动到根目录
}
