package model

import "mvtable/internal/storage/db"

type InviteCode struct {
	db.Model
	UserID string `gorm:"column:user_id"`
	Value  string `gorm:"column:value"`
	UsedBy string `gorm:"column:used_by"`
	IsUsed bool   `gorm:"column:is_used;default:false"`
}

func (*InviteCode) TableName() string {
	return "invite_code"
}

type BatchCreateInviteCodeReq struct {
	Count int `json:"count" binding:"required"`
}

type GetInviteCodeListReq struct {
	IsUsed *bool `form:"isUsed"`
	db.Pagination
}

type InviteCodeListItem struct {
	ID        string `json:"id" gorm:"column:id"`
	Value     string `json:"value" gorm:"column:value"`
	IsUsed    bool   `json:"isUsed" gorm:"column:is_used"`
	UsedBy    string `json:"usedBy" gorm:"column:used_by"`
	Creator   string `json:"creator" gorm:"column:creator"`
	CreatedAt string `json:"createdAt" gorm:"column:created_at"`
}

type GetInviteCodeListRes struct {
	Total int64                 `json:"total"`
	List  []*InviteCodeListItem `json:"list"`
}
