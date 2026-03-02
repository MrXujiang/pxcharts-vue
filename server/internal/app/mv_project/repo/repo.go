package repo

import (
	"mvtable/internal/app/mv_project/model"
	model2 "mvtable/internal/app/mv_project_recent/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetRecentProjects(db *gorm.DB, userId string) (projects []*model.QueryProjectItem, total int64, err error) {
	projects = make([]*model.QueryProjectItem, 0)

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err = db.Table("mv_project_recent AS r").
		Select("p.id, p.name, p.description, p.created_at, p.updated_at").
		Joins("JOIN mv_project AS p ON p.id = r.project_id").
		Where("r.user_id = ? AND r.last_accessed_at >= ?", userId, sevenDaysAgo).
		Order("r.last_accessed_at DESC").
		Count(&total).
		Scan(&projects).Error
	return
}

func GetCreatedProjects(db *gorm.DB, userId string, page, size int) (projects []*model.QueryProjectItem, total int64, err error) {
	projects = make([]*model.QueryProjectItem, 0, size)

	sql := db.Model(&model.MvProject{}).
		Select("id, name, description, created_at, updated_at").
		Where("user_id = ? AND deleted_at IS NULL", userId).
		Order("created_at DESC").
		Count(&total)

	if page > 0 && size > 0 {
		sql = sql.Offset((page - 1) * size).Limit(size)
	}

	err = sql.Scan(&projects).Error
	return projects, total, err
}

func GetSharedProjects(db *gorm.DB, userId string, page, size int) (projects []*model.QueryProjectItem, total int64, err error) {
	projects = make([]*model.QueryProjectItem, 0, size)

	sql := db.Table("mv_project_perm AS perm").
		Select("p.id, p.name, p.description, p.created_at, p.updated_at").
		Joins("JOIN mv_project AS p ON p.id = perm.project_id").
		Where("perm.target = ? AND perm.target_id = ? AND p.user_id <> ? AND p.deleted_at IS NULL AND perm.deleted_at IS NULL", "user", userId, userId).
		Order("p.created_at DESC").
		Count(&total)

	if page > 0 && size > 0 {
		sql = sql.Offset((page - 1) * size).Limit(size)
	}

	err = sql.Scan(&projects).Error
	return
}

func GetFavoriteProjects(db *gorm.DB, userId string, page, size int) (projects []*model.QueryProjectItem, total int64, err error) {
	projects = make([]*model.QueryProjectItem, 0, size)

	sql := db.Table("mv_project_favorite AS f").
		Select("p.id, p.name, p.description, p.created_at, p.updated_at").
		Joins("JOIN mv_project AS p ON p.id = f.project_id").
		Where("f.user_id = ? AND f.deleted_at IS NULL", userId).
		Order("f.created_at DESC").
		Count(&total)

	if page > 0 && size > 0 {
		sql = sql.Offset((page - 1) * size).Limit(size)
	}

	err = sql.Scan(&projects).Error
	return
}

func UpsertRecentProject(db *gorm.DB, userId, projectId string) error {
	recent := model2.MvProjectRecent{
		UserID:         userId,
		ProjectID:      projectId,
		LastAccessedAt: time.Now(),
	}
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "project_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"last_accessed_at", "updated_at"}),
	}).Create(&recent).Error
}
