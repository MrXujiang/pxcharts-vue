package repo

import (
	model2 "mvtable/internal/app/mv_table_schema/model"

	"gorm.io/gorm"
)

func GetTableSchemaByViewFormID(db *gorm.DB, viewFormID string) (tableSchema *model2.MvTableSchema, err error) {
	err = db.Table("mv_view_form as vf").
		Select("ts.*").
		Joins("JOIN mv_view v ON v.id = vf.view_id").
		Joins("JOIN mv_table_schema ts ON ts.id = v.table_schema_id").
		Where("vf.id = ?", viewFormID).
		First(&tableSchema).
		Error
	return tableSchema, err
}
