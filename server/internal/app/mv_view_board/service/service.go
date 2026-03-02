package service

import (
	"mvtable/internal/app/mv_view_board/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvViewBoardService struct{}

func NewMvViewBoardService() *MvViewBoardService {
	return &MvViewBoardService{}
}

func (s *MvViewBoardService) UpdateMvViewBoard(req *model.UpdateMvViewBoardReq) error {
	var (
		viewBoard    *model.MvViewBoard
		updateFields []string
		err          error
	)
	viewBoard, err = db.Get[model.MvViewBoard](db.GetDB(), map[string]any{"view_id": req.ViewID})

	if err != nil {
		log.Error("get mv view board error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	if viewBoard == nil {
		return errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	if req.ShowFieldTitle != nil {
		viewBoard.ShowFieldTitle = *req.ShowFieldTitle
		updateFields = append(updateFields, "show_field_title")
	}
	if req.FilterConfig != nil {
		viewBoard.FilterConfig = *req.FilterConfig
		updateFields = append(updateFields, "filter_config")
	}
	if req.GroupConfig != nil {
		viewBoard.GroupConfig = *req.GroupConfig
		updateFields = append(updateFields, "group_config")
	}
	if req.SortConfig != nil {
		viewBoard.SortConfig = *req.SortConfig
		updateFields = append(updateFields, "sort_config")
	}

	if err = db.Update(db.GetDB(), viewBoard, map[string]any{"view_id": req.ViewID}, updateFields...); err != nil {
		log.Error("update mv view board error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	return nil
}
