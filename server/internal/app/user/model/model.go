package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/utils"
	"mvtable/internal/storage/db"
	"time"
)

type User struct {
	db.Model
	Email           string                 `gorm:"column:email" json:"email"`
	Password        string                 `gorm:"column:password" json:"password"`
	Nickname        string                 `gorm:"column:nickname;default:''" json:"nickname"`
	Avatar          string                 `gorm:"column:avatar;default:''" json:"avatar"`
	Status          int8                   `gorm:"column:status;default:1" json:"status"`
	LastLoginAt     *time.Time             `gorm:"column:last_login_at" json:"lastLoginAt"`
	Tags            string                 `gorm:"column:tags;default:''" json:"tags"`
	Identity        constants.UserIdentity `gorm:"column:identity;default:1" json:"identity"`
	MaxProjectCount int                    `gorm:"column:max_project_count;default:0" json:"maxProjectCount"`
	MaxTeamCount    int                    `gorm:"column:max_team_count;default:0" json:"maxTeamCount"`
	Role            string                 `gorm:"column:role;default:'user'" json:"role"`
}

func (*User) TableName() string {
	return "user"
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		Nickname:        u.Nickname,
		Avatar:          u.Avatar,
		Status:          u.Status,
		Tags:            utils.SplitAndFilter(u.Tags, ","),
		Identity:        u.Identity,
		Role:            u.Role,
		MaxProjectCount: u.MaxProjectCount,
		MaxTeamCount:    u.MaxTeamCount,
		LastLoginAt:     u.LastLoginAt,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

type UserResponse struct {
	ID              string                 `json:"id"`
	Email           string                 `json:"email"`
	Nickname        string                 `json:"nickname"`
	Avatar          string                 `json:"avatar"`
	Status          int8                   `json:"status"`
	Tags            []string               `json:"tags"`
	Identity        constants.UserIdentity `json:"identity"`
	Role            string                 `json:"role"`
	MaxProjectCount int                    `json:"maxProjectCount"`
	MaxTeamCount    int                    `json:"maxTeamCount"`
	LastLoginAt     *time.Time             `json:"lastLoginAt"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

type RegisterReq struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"`
}

type SendEmailVerifyCodeReq struct {
	Email string `json:"email" binding:"required"`
}

type RegisterRes struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
	UserInfo     *UserResponse `json:"userInfo"`
}

type InviteCodeRegisterReq struct {
	Email      string `json:"email" binding:"required"`
	InviteCode string `json:"inviteCode" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UpdateInfoReq struct {
	Nickname *string   `json:"nickname"`
	Avatar   *string   `json:"avatar"`
	Tags     *[]string `json:"tags"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRes struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
	UserInfo     *UserResponse `json:"userInfo"`
}

type ResetPasswordReq struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	VerifyCode string `json:"verifyCode" binding:"required"`
}

// 管理后台

type AdminLoginReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminGetUserListReq struct {
	db.Pagination
	Email    string `form:"email"`
	Nickname string `form:"nickname"`
	Status   int8   `form:"status"`
}

type AdminGetUserListRes struct {
	List  []*UserResponse `json:"list"`
	Total int64           `json:"total"`
}

type AdminUpdateInfoReq struct {
	ID       string                  `json:"id" binding:"required"`
	Nickname *string                 `json:"nickname"`
	Avatar   *string                 `json:"avatar"`
	Tags     *[]string               `json:"tags"`
	Status   *int8                   `json:"status"`
	Role     *string                 `json:"role"`
	Identity *constants.UserIdentity `json:"identity"`
}

type AdminDeleteUserReq struct {
	ID string `json:"id" binding:"required"`
}
