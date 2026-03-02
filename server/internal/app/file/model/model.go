package model

import "mvtable/internal/storage/db"

type File struct {
	db.Model
	UserID   string `gorm:"column:user_id" json:"userId"`
	Filename string `gorm:"column:filename;default:''" json:"filename"`
	Filesize int64  `gorm:"column:filesize;default:0" json:"filesize"`
	Filetype string `gorm:"column:filetype;default:''" json:"filetype"`
	OSS      string `gorm:"column:oss;default:''" json:"oss"`
	Path     string `gorm:"column:path;default:''" json:"path"`
	Remark   string `gorm:"column:remark;default:''" json:"remark"`
}

func (File) TableName() string {
	return "file"
}

type UploadRes struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type AdminGetListReq struct {
	db.Pagination
	Filename string `form:"filename"`
	OSS      string `form:"oss"`
	Remark   string `form:"remark"`
}

type AdminGetListItem struct {
	File
	Uploader string `json:"uploader"`
	URL      string `json:"url"`
}

type AdminGetListRes struct {
	List  []*AdminGetListItem `json:"list"`
	Total int64               `json:"total"`
}
