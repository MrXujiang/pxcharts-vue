package repo

import (
	"mvtable/internal/app/invite_code/model"
	"mvtable/internal/storage/db"
)

func List(page, size int, isUsed *bool) (list []*model.InviteCodeListItem, total int64, err error) {
	sql := db.GetDB().Model(&model.InviteCodeListItem{}).
		Table("invite_code AS ic").
		Select("ic.id, ic.value, ic.is_used, ic.used_by, ic.created_at, u.nickname AS creator").
		Joins(`LEFT JOIN "user" AS u ON ic.user_id = u.id`).
		Where("ic.deleted_at IS NULL")

	if isUsed != nil {
		sql = sql.Where("ic.is_used = ?", *isUsed)
	}

	err = sql.Count(&total).Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}
