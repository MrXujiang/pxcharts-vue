package service

import (
	"mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvViewTableService struct{}

func NewMvViewTableService() *MvViewTableService {
	return &MvViewTableService{}
}

func (s *MvViewTableService) UpdateMvViewTable(req *model.UpdateMvViewTableReq) error {
	var (
		viewTable    *model.MvViewTable
		updateFields []string
		err          error
	)
	viewTable, err = db.Get[model.MvViewTable](db.GetDB(), map[string]any{"view_id": req.ViewID})

	if err != nil {
		log.Error("get mv view table error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	if viewTable == nil {
		return errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	if req.FilterConfig != nil {
		viewTable.FilterConfig = *req.FilterConfig
		updateFields = append(updateFields, "filter_config")
	}
	if req.GroupConfig != nil {
		viewTable.GroupConfig = *req.GroupConfig
		updateFields = append(updateFields, "group_config")
	}
	if req.SortConfig != nil {
		viewTable.SortConfig = *req.SortConfig
		updateFields = append(updateFields, "sort_config")
	}
	if req.RowHeight != nil {
		viewTable.RowHeight = *req.RowHeight
		updateFields = append(updateFields, "row_height")
	}
	if req.ColorConfig != nil {
		viewTable.ColorConfig = *req.ColorConfig
		updateFields = append(updateFields, "color_config")
	}

	if err = db.Update(db.GetDB(), viewTable, map[string]any{"view_id": req.ViewID}, updateFields...); err != nil {
		log.Error("update mv view table error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	return nil
}
