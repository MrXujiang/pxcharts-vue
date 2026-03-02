package service

import (
	model2 "mvtable/internal/app/mv_project/model"
	"mvtable/internal/app/mv_project_state/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvProjectStateService struct{}

func NewMvProjectStateService() *MvProjectStateService {
	return &MvProjectStateService{}
}

func (s *MvProjectStateService) SetShareRange(userId string, req *model.SetShareRangeReq) error {
	var (
		err     error
		state   *model.MvProjectState
		project *model2.MvProject
	)

	if req.ShareRange == constants.ShareRangeTeam && req.TeamAction == nil {
		return errorx.New(errorx.ErrOperationFailed, "设置团队内公开时，团队内权限不能为空")
	}

	// 判断项目是否存在
	project, err = db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if project == nil {
		log.Error("project not found", zap.String("projectId", req.ProjectID))
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	// 获取分享配置
	state, err = db.Get[model.MvProjectState](db.GetDB(), map[string]any{"project_id": req.ProjectID})
	if err != nil {
		log.Error("get project state error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}

	if state == nil {
		log.Error("project state not found", zap.String("projectId", req.ProjectID))
		return errorx.New(errorx.ErrNotFound, "获取分享配置失败")
	}

	// 更新分享配置
	switch req.ShareRange {
	case constants.ShareRangePrivate:
		state.ShareRange = constants.ShareRangePrivate
		state.TeamAction = ""
		state.AccessPassword = ""
	case constants.ShareRangePublic:
		state.ShareRange = constants.ShareRangePublic
		state.TeamAction = ""
		if req.AccessPassword != nil && *req.AccessPassword != "" {
			state.AccessPassword = *req.AccessPassword
		} else {
			state.AccessPassword = ""
		}
	case constants.ShareRangeTeam:
		state.ShareRange = constants.ShareRangeTeam
		state.TeamAction = req.TeamAction.String()
		state.AccessPassword = ""
	default:
		return errorx.New(errorx.ErrOperationFailed, "设置的分享范围无效")
	}

	return nil
}
