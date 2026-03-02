package model

import (
	"mvtable/internal/storage/db"

	"gorm.io/datatypes"
)

type WebsiteConfig struct {
	db.Model
	Schema datatypes.JSON `gorm:"column:schema" json:"schema"`
}

func (*WebsiteConfig) TableName() string {
	return "website_config"
}

type WebsiteConfigSchema struct {
	Logo            string   `json:"logo"`
	Slogan          string   `json:"slogan"`
	RegisterMethod  []string `json:"registerMethod"`
	CommunityRQCode string   `json:"communityRQCode"`
	InviteRQCode    string   `json:"inviteRQCode"`
	Copyright       string   `json:"copyright"`
	UploadMethod    string   `json:"uploadMethod"`
}
