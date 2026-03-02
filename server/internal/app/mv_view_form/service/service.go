package service

import (
	"mvtable/internal/app/mv_view_form/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type MvViewFormService struct{}

func NewMvViewFormService() *MvViewFormService {
	return &MvViewFormService{}
}

func (s *MvViewFormService) UpdateMvViewForm(req *model.UpdateMvViewFormReq) error {
	var (
		viewForm     *model.MvViewForm
		updateFields []string
		err          error
	)
	viewForm, err = db.Get[model.MvViewForm](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv view form error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	if viewForm == nil {
		return errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	if req.Name != nil {
		viewForm.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.Description != nil {
		viewForm.Description = *req.Description
		updateFields = append(updateFields, "description")
	}
	if req.Cover != nil {
		viewForm.Cover = *req.Cover
		updateFields = append(updateFields, "cover")
	}
	if req.Layout != nil {
		viewForm.Layout = *req.Layout
		updateFields = append(updateFields, "layout")
	}
	if req.Stats != nil {
		viewForm.Stats = datatypes.JSON(*req.Stats)
		updateFields = append(updateFields, "stats")
	}
	if req.EnableSharing != nil {
		viewForm.EnableSharing = *req.EnableSharing
		updateFields = append(updateFields, "enable_sharing")
	}
	if req.EnableAnonymous != nil {
		viewForm.EnableAnonymous = *req.EnableAnonymous
		updateFields = append(updateFields, "enable_anonymous")
	}
	if req.Filter != nil {
		viewForm.Filter = *req.Filter
		updateFields = append(updateFields, "filter")
	}
	if req.FilterConfig != nil {
		viewForm.FilterConfig = *req.FilterConfig
		updateFields = append(updateFields, "filter_config")
	}
	if req.EnableNoLogin != nil {
		viewForm.EnableNoLogin = *req.EnableNoLogin
		updateFields = append(updateFields, "enable_no_login")
	}
	if req.EnableLimitSubmit != nil {
		viewForm.EnableLimitSubmit = *req.EnableLimitSubmit
		updateFields = append(updateFields, "enable_limit_submit")
	}
	if req.LimitSubmitType != nil {
		viewForm.LimitSubmitType = *req.LimitSubmitType
		updateFields = append(updateFields, "limit_submit_type")
	}
	if req.EnableLimitCollect != nil {
		viewForm.EnableLimitCollect = *req.EnableLimitCollect
		updateFields = append(updateFields, "enable_limit_collect")
	}
	if req.LimitCollectCount != nil {
		viewForm.LimitCollectCount = *req.LimitCollectCount
		updateFields = append(updateFields, "limit_collect_count")
	}
	if req.EnableEditAfterSubmit != nil {
		viewForm.EnableEditAfterSubmit = *req.EnableEditAfterSubmit
		updateFields = append(updateFields, "enable_edit_after_submit")
	}
	if req.Config != nil {
		viewForm.Config = *req.Config
		updateFields = append(updateFields, "config")
	}

	if err = db.Update(db.GetDB(), viewForm, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv view form error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	return nil
}
