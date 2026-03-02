package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
)

type MvProjectPerm struct {
	db.Model
	TargetID  string                `gorm:"column:target_id" json:"targetID"`
	Target    string                `gorm:"column:target;default:'user'" json:"target"`
	ProjectID string                `gorm:"column:project_id" json:"projectId"`
	Role      constants.ProjectRole `gorm:"column:role;default:'reader'" json:"role"`
}

func (MvProjectPerm) TableName() string {
	return "mv_project_perm"
}

type Member struct {
	ID     string
	Target string
}

type BatchCreateMemberReq struct {
	MemberList []Member              `json:"memberList" binding:"required"`
	ProjectID  string                `json:"projectId" binding:"required"`
	Role       constants.ProjectRole `json:"role" binding:"required"`
}

type UpdateMemberReq struct {
	ProjectID string                `json:"projectId" binding:"required"`
	TargetID  string                `json:"targetID" binding:"required"`
	Target    string                `json:"target" binding:"required"`
	Role      constants.ProjectRole `json:"role" binding:"required"`
}

type DeleteMemberReq struct {
	ProjectID string `json:"projectId" binding:"required"`
	TargetID  string `json:"targetID" binding:"required"`
	Target    string `json:"target" binding:"required"`
}

type SearchMemberReq struct {
	ProjectID string `form:"projectId" binding:"required"`
	Keywords  string `form:"keywords" binding:"required"`
}

type MemberListItem struct {
	TargetID string `gorm:"column:target_id" json:"targetID"` // 用户id/团队id
	Target   string `gorm:"column:target" json:"target"`      // user/team
	Name     string `gorm:"column:name" json:"name"`          // 团队名/用户昵称
	Avatar   string `gorm:"column:avatar" json:"avatar"`      // 团队头像/用户头像
}

type SearchMemberRes struct {
	List []*MemberListItem `json:"list"`
}

// ProjectMemberListReq 查询项目成员请求
type ProjectMemberListReq struct {
	ProjectID string `form:"projectId" binding:"required"`
}

// ProjectMemberItem 项目成员信息（包含角色）
type ProjectMemberItem struct {
	TargetID string                `gorm:"column:target_id" json:"targetID"`
	Target   string                `gorm:"column:target" json:"target"`
	Name     string                `gorm:"column:name" json:"name"`
	Avatar   string                `gorm:"column:avatar" json:"avatar"`
	Role     constants.ProjectRole `gorm:"column:role" json:"role"`
}

// ProjectMemberListRes 查询项目成员响应
type ProjectMemberListRes struct {
	List []*ProjectMemberItem `json:"list"`
}
