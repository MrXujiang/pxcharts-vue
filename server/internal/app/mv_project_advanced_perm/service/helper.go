package service

import (
	"mvtable/internal/app/mv_project_advanced_perm/model"
	"mvtable/internal/pkg/constants"
)

func getDefaultOwnerPerm(tableSchemaId string) *model.MvProjectAdvancedPerm {
	return &model.MvProjectAdvancedPerm{
		TableSchemaID:   tableSchemaId,
		Role:            constants.ProjectActionOwner,
		DataAction:      constants.ActionMange,
		CanAdd:          true,
		CanDelete:       true,
		OperateRange:    constants.AdvancedOperateRangeAll,
		FieldAccess:     constants.AdvancedFieldAccessAll,
		CustomFieldPerm: []model.FieldPerm{},
		CanOperateView:  true,
		ViewAccess:      constants.AdvancedViewAccessAll,
		CanCheckViews:   []string{},
	}
}

func getDefaultAdminPerm(tableSchemaId string) *model.MvProjectAdvancedPerm {
	return &model.MvProjectAdvancedPerm{
		TableSchemaID:   tableSchemaId,
		Role:            constants.ProjectActionAdmin,
		DataAction:      constants.ActionMange,
		CanAdd:          true,
		CanDelete:       true,
		OperateRange:    constants.AdvancedOperateRangeAll,
		FieldAccess:     constants.AdvancedFieldAccessAll,
		CustomFieldPerm: []model.FieldPerm{},
		CanOperateView:  true,
		ViewAccess:      constants.AdvancedViewAccessAll,
		CanCheckViews:   []string{},
	}
}

func getDefaultEditorPerm(tableSchemaId string) *model.MvProjectAdvancedPerm {
	return &model.MvProjectAdvancedPerm{
		TableSchemaID:   tableSchemaId,
		Role:            constants.ProjectActionEditor,
		DataAction:      constants.ActionEdit,
		CanAdd:          true,
		CanDelete:       true,
		OperateRange:    constants.AdvancedOperateRangeAll,
		FieldAccess:     constants.AdvancedFieldAccessAll,
		CustomFieldPerm: []model.FieldPerm{},
		CanOperateView:  true,
		ViewAccess:      constants.AdvancedViewAccessAll,
		CanCheckViews:   []string{},
	}
}

func getDefaultReaderPerm(tableSchemaId string) *model.MvProjectAdvancedPerm {
	return &model.MvProjectAdvancedPerm{
		TableSchemaID:   tableSchemaId,
		Role:            constants.ProjectActionReader,
		DataAction:      constants.ActionRead,
		CanAdd:          false,
		CanDelete:       false,
		OperateRange:    constants.AdvancedOperateRangeAll,
		FieldAccess:     constants.AdvancedFieldAccessAll,
		CustomFieldPerm: []model.FieldPerm{},
		CanOperateView:  false,
		ViewAccess:      constants.AdvancedViewAccessAll,
		CanCheckViews:   []string{},
	}
}
