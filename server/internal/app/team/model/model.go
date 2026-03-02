package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"time"
)

type Team struct {
	db.Model
	UserID      string `gorm:"column:user_id" json:"userId"`
	Name        string `gorm:"column:name;default:''" json:"name"`
	Description string `gorm:"column:description;default:''" json:"description"`
	Logo        string `gorm:"column:logo;default:''" json:"logo"`
}

func (Team) TableName() string {
	return "team"
}

type TeamListItem struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Logo              string                 `json:"logo"`
	MemberCount       int                    `json:"memberCount"`
	DisplayAvatarList []string               `json:"displayAvatarList"`
	Creator           string                 `json:"creator"`
	Identity          constants.TeamIdentity `json:"identity"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

func (a *TeamListItem) FillFromTeam(team *Team) {
	a.ID = team.ID
	a.Name = team.Name
	a.Description = team.Description
	a.Logo = team.Logo
	a.CreatedAt = team.CreatedAt
	a.UpdatedAt = team.UpdatedAt
}

type ListTeamReq struct {
	List []TeamListItem `json:"list"`
}

type UpdateTeamReq struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type CreateTeamReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type DeleteTeamReq struct {
	ID string `json:"id" binding:"required"`
}

type TeamMember struct {
	UserID    string                 `gorm:"column:user_id" json:"userId"`
	Avatar    string                 `gorm:"column:avatar" json:"avatar"`
	Nickname  string                 `gorm:"column:nickname" json:"nickname"`
	Email     string                 `gorm:"column:email" json:"email"`
	Identity  constants.TeamIdentity `gorm:"column:identity" json:"identity"`
	UpdatedAt time.Time              `gorm:"column:updated_at" json:"updatedAt"`
	CreatedAt time.Time              `gorm:"column:created_at" json:"createdAt"`
}

type GetMemberListReq struct {
	db.Pagination
	TeamID     string `form:"teamId" binding:"required"`
	SearchWord string `form:"searchWord"`
}

type GetMemberListRes struct {
	List  []*TeamMember `json:"list"`
	Total int64         `json:"total"`
}

type AddMemberReq struct {
	TeamID   string                 `json:"teamId" binding:"required"`
	UserID   string                 `json:"userId" binding:"required"`
	Identity constants.TeamIdentity `json:"identity" binding:"required"`
}

type UpdateMemberReq struct {
	TeamID   string                 `json:"teamId" binding:"required"`
	UserID   string                 `json:"userId" binding:"required"`
	Identity constants.TeamIdentity `json:"identity" binding:"required"`
}

type DeleteMemberReq struct {
	TeamID string `json:"teamId" binding:"required"`
	UserID string `json:"userId" binding:"required"`
}

type AdminGetTeamListReq struct {
	db.Pagination
	SearchWord string `form:"searchWord"`
}

type AdminGetTeamListItem struct {
	ID              string    `gorm:"column:id" json:"id"`
	CreatorNickname string    `gorm:"column:creator_nickname" json:"creatorNickname"`
	CreatorEmail    string    `gorm:"column:creator_email" json:"creatorEmail"`
	Name            string    `gorm:"column:name" json:"name"`
	Description     string    `gorm:"column:description" json:"description"`
	Logo            string    `gorm:"column:logo" json:"logo"`
	MemberCount     int64     `gorm:"column:member_count" json:"memberCount"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

type AdminGetTeamListRes struct {
	List  []*AdminGetTeamListItem `json:"list"`
	Total int64                   `json:"total"`
}

type BatchGetUsersReq struct {
	UserIDs []string `json:"userIds" binding:"required,min=1"`
}

type BatchGetUsersRes struct {
	List []*UserBasicInfo `json:"list"`
}

type SearchUserReq struct {
	SearchWord string `form:"searchWord" binding:"required"`
}

type SearchUserRes struct {
	List []*UserBasicInfo `json:"list"`
}

type UserBasicInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
