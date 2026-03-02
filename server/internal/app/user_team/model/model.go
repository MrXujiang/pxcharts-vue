package model

import (
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
)

type UserTeam struct {
	db.Model
	UserID   string                 `gorm:"column:user_id" json:"userId"`
	TeamID   string                 `gorm:"column:team_id" json:"teamId"`
	Identity constants.TeamIdentity `gorm:"column:identity;default:'member'" json:"identity"`
}

func (UserTeam) TableName() string {
	return "user_team"
}
