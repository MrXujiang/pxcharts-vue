package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
)

type MvProjectState struct {
	db.Model
	ProjectID      string                `gorm:"column:project_id" json:"projectId"`
	ShareRange     constants.ShareRange  `gorm:"column:share_range;default:1" json:"shareRange"`          // 分享范围: 1=仅协作者可见, 2=互联网公开, 3=团队内公开
	TeamAction     constants.ProjectRole `gorm:"column:team_action;default:'read'" json:"teamAction"`     // 团队内权限(分享范围为团队内公开时该字段有意义): read=团队内可查看, edit=团队内可编辑
	AccessPassword string                `gorm:"column:access_password;default:''" json:"accessPassword"` // 访问密码(分享范围为互联网公开时该字段有意义)
}

func (MvProjectState) TableName() string {
	return "mv_project_state"
}

type SetShareRangeReq struct {
	ProjectID      string                 `json:"projectId" binding:"required"`
	ShareRange     constants.ShareRange   `json:"shareRange" binding:"required"`
	TeamAction     *constants.ProjectRole `json:"teamAction"`
	AccessPassword *string                `json:"accessPassword"`
}
