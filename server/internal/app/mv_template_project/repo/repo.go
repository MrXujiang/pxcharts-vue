package repo

import (
	"mvtable/internal/app/mv_template_project/model"
	"mvtable/internal/storage/db"
)

// tempQueryItem 临时结构体，用于扫描数据库（使用 pq.StringArray）

func QueryTemplateProjectList(page, size int, tag string) (list []*model.QueryMvTemplateProjectItem, total int64, err error) {
	sql := db.GetDB().Model(&model.MvTemplateProject{}).
		Table("mv_template_project AS tp").
		Select("tp.id, tp.name, tp.description, tp.cover, tp.tags, tp.use_count, u.nickname AS creator").
		Joins(`LEFT JOIN "user" AS u ON tp.user_id = u.id`).
		Where("tp.deleted_at IS NULL")

	// 如果提供了标签，使用PostgreSQL的数组包含操作符 @>
	if tag != "" {
		sql = sql.Where("tp.tags @> ARRAY[?]::text[]", tag)
	}

	err = sql.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && size > 0 {
		sql = sql.Offset((page - 1) * size).Limit(size)
	}

	// 使用临时结构体扫描数据库
	var tempList []model.TempQueryItem
	err = sql.Order("tp.created_at DESC").Find(&tempList).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为目标结构体
	list = make([]*model.QueryMvTemplateProjectItem, len(tempList))
	for i, item := range tempList {
		list[i] = &model.QueryMvTemplateProjectItem{
			ID:          item.ID,
			Creator:     item.Creator,
			Name:        item.Name,
			Description: item.Description,
			Cover:       item.Cover,
			Tags:        []string(item.Tags), // 转换为 []string
			UseCount:    item.UseCount,
		}
	}

	return list, total, nil
}
