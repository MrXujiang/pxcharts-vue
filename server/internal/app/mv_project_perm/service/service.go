package service

import (
	model2 "mvtable/internal/app/mv_project/model"
	"mvtable/internal/app/mv_project_perm/model"
	"mvtable/internal/app/mv_project_perm/repo"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"slices"

	"go.uber.org/zap"
)

type MvProjectPermissionService struct{}

func NewMvProjectPermissionService() *MvProjectPermissionService {
	return &MvProjectPermissionService{}
}

func (s *MvProjectPermissionService) BatchCreateMember(userId string, req *model.BatchCreateMemberReq) error {
	if len(req.MemberList) == 0 {
		return errorx.BadRequest("协作者不能为空")
	}

	if !slices.Contains([]constants.ProjectRole{constants.ProjectActionReader, constants.ProjectActionEditor, constants.ProjectActionAdmin}, req.Role) {
		return errorx.BadRequest("协作者角色不合法")
	}

	// 判断项目是否存在
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if project == nil {
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	permissions := make([]*model.MvProjectPerm, len(req.MemberList))

	for i, v := range req.MemberList {
		permissions[i] = &model.MvProjectPerm{
			ProjectID: req.ProjectID,
			TargetID:  v.ID,
			Target:    v.Target,
			Role:      req.Role,
		}
	}

	if err = db.CreateBatch(db.GetDB(), permissions); err != nil {
		log.Error("create mv project perm error", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}
	return nil
}

func (s *MvProjectPermissionService) UpdateMember(userId string, req *model.UpdateMemberReq) error {
	if !slices.Contains([]constants.ProjectRole{constants.ProjectActionReader, constants.ProjectActionEditor, constants.ProjectActionAdmin}, req.Role) {
		return errorx.BadRequest("操作类型不合法")
	}

	// 判断项目是否存在
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if project == nil {
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	// 判断协作者是否存在
	perm, err := db.Get[model.MvProjectPerm](db.GetDB(), map[string]any{"project_id": req.ProjectID, "target_id": req.TargetID, "target": req.Target})
	if err != nil {
		log.Error("get mv project perm error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if perm == nil {
		return errorx.New(errorx.ErrNotFound, "协作者不存在")
	}

	perm.Role = req.Role
	if err = db.Update(db.GetDB(), perm, map[string]any{"id": perm.ID}, "role"); err != nil {
		log.Error("update mv project perm error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

func (s *MvProjectPermissionService) DeleteMember(userId string, req *model.DeleteMemberReq) error {
	// 判断项目是否存在
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if project == nil {
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	// 判断协作者是否存在
	perm, err := db.Get[model.MvProjectPerm](db.GetDB(), map[string]any{"project_id": req.ProjectID, "target_id": req.TargetID, "target": req.Target})
	if err != nil {
		log.Error("get mv project perm error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if perm == nil {
		return errorx.New(errorx.ErrNotFound, "协作者不存在")
	}

	if err = db.Delete[model.MvProjectPerm](db.GetDB(), map[string]any{"id": perm.ID}); err != nil {
		log.Error("delete mv project perm error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}
	return nil
}

func (s *MvProjectPermissionService) SearchMember(userId string, req *model.SearchMemberReq) (*model.SearchMemberRes, error) {
	list, err := repo.QueryAvailableMembers(db.GetDB(), userId, req.ProjectID, req.Keywords)
	if err != nil {
		log.Error("search member error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	return &model.SearchMemberRes{
		List: list,
	}, nil
}

func (s *MvProjectPermissionService) ListProjectMembers(userId string, req *model.ProjectMemberListReq) (*model.ProjectMemberListRes, error) {
	// 校验项目存在
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if project == nil {
		return nil, errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	list, err := repo.QueryProjectMembers(db.GetDB(), req.ProjectID)
	if err != nil {
		log.Error("query project members error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	return &model.ProjectMemberListRes{
		List: list,
	}, nil
}
