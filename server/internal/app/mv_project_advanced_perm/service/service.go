package service

import (
	model2 "mvtable/internal/app/mv_project/model"
	"mvtable/internal/app/mv_project_advanced_perm/model"
	model3 "mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvProjectAdvancedPermService struct{}

func NewMvProjectAdvancedPermService() *MvProjectAdvancedPermService {
	return &MvProjectAdvancedPermService{}
}

func (s *MvProjectAdvancedPermService) EnableMvProjectAdvancedPermReq(req *model.EnableMvProjectAdvancedPermReq) error {
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if project == nil {
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	if project.EnableAdvancedPerm {
		return errorx.New(errorx.ErrOperationFailed, "项目已开启高级权限")
	}

	// 查询项目下所有的表格
	tableSchemas, _, err := db.List[model3.MvTableSchema](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("get table schemas error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, v := range tableSchemas {
			ReaderPerm := getDefaultReaderPerm(v.ID)
			EditorPerm := getDefaultEditorPerm(v.ID)
			AdminPerm := getDefaultAdminPerm(v.ID)
			OwnerPerm := getDefaultOwnerPerm(v.ID)

			if err := db.Create(tx, ReaderPerm); err != nil {
				log.Error("create reader perm error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
			if err := db.Create(tx, EditorPerm); err != nil {
				log.Error("create editor perm error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
			if err := db.Create(tx, AdminPerm); err != nil {
				log.Error("create admin perm error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
			if err := db.Create(tx, OwnerPerm); err != nil {
				log.Error("create owner perm error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}

		// 更新项目高级权限字段
		project.EnableAdvancedPerm = true
		if err := db.Update(db.GetDB(), project, map[string]any{"id": req.ProjectID}, "enable_advanced_perm"); err != nil {
			log.Error("update project advanced perm error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}

		return nil
	})

	if err != nil {
		log.Error("enable mv project advanced perm error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	return nil
}

func (s *MvProjectAdvancedPermService) DisableMvProjectAdvancedPermReq(req *model.DisableMvProjectAdvancedPermReq) error {
	project, err := db.Get[model2.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	if project == nil {
		return errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	if !project.EnableAdvancedPerm {
		return errorx.New(errorx.ErrOperationFailed, "项目未开启高级权限")
	}

	// 直接设置项目的高级权限为关闭
	project.EnableAdvancedPerm = false
	if err := db.Update(db.GetDB(), project, map[string]any{"id": req.ProjectID}, "enable_advanced_perm"); err != nil {
		log.Error("update project advanced perm error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	return nil
}

func (s *MvProjectAdvancedPermService) GetMvProjectAdvancedPerm(req *model.GetMvProjectAdvancedPermReq) (*model.GetMvProjectAdvancedPermRes, error) {
	perm, err := db.Get[model.MvProjectAdvancedPerm](db.GetDB(), map[string]any{"role": req.Role, "table_schema_id": req.TableSchemaID})
	if err != nil {
		log.Error("get project advanced perm error", zap.Error(err))
		return nil, errorx.InternalServerError("操作失败")
	}
	if perm == nil {
		return nil, errorx.New(errorx.ErrNotFound, "高级权限不存在")
	}

	res := &model.GetMvProjectAdvancedPermRes{}
	res.FillFrom(perm)
	return res, nil
}

func (s *MvProjectAdvancedPermService) UpdateMvProjectAdvancedPermReq(req *model.UpdateMvProjectAdvancedPermReq) error {
	perm, err := db.Get[model.MvProjectAdvancedPerm](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get project advanced perm error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	if perm == nil {
		return errorx.New(errorx.ErrNotFound, "高级权限不存在")
	}

	perm.DataAction = req.DataAction

	if req.CanAdd != nil {
		perm.CanAdd = *req.CanAdd
	}
	if req.CanDelete != nil {
		perm.CanDelete = *req.CanDelete
	}
	if req.OperateRange != nil {
		perm.OperateRange = *req.OperateRange
	}
	if req.FieldAccess != nil {
		perm.FieldAccess = *req.FieldAccess
	}
	if req.CustomFieldPerm != nil {
		perm.CustomFieldPerm = *req.CustomFieldPerm
	}
	if req.CanOperateView != nil {
		perm.CanOperateView = *req.CanOperateView
	}
	if req.ViewAccess != nil {
		perm.ViewAccess = *req.ViewAccess
	}
	if req.CanCheckViews != nil {
		perm.CanCheckViews = *req.CanCheckViews
	}

	if err := db.Update(db.GetDB(), perm, map[string]any{"id": req.ID}, "data_action", "can_add", "can_delete", "operate_range", "field_access", "custom_field_perm", "can_operate_view", "view_access", "can_check_views"); err != nil {
		log.Error("update project advanced perm error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	return nil
}
