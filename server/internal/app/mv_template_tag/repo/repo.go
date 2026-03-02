package repo

import (
	"mvtable/internal/app/mv_template_tag/model"
	"mvtable/internal/storage/db"
)

func QueryTagList(searchWord string) (list []*model.QueryMvTemplateTagItem, err error) {
	sql := db.GetDB().Model(&model.MvTemplateTag{}).
		Select("id, name, description").
		Where("deleted_at IS NULL")

	if searchWord != "" {
		sql = sql.Where("name LIKE ? OR description LIKE ?", "%"+searchWord+"%", "%"+searchWord+"%")
	}

	err = sql.Order("order_index ASC").Find(&list).Error
	return list, err
}
