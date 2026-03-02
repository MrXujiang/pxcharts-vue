package service

import (
	"mvtable/internal/app/mv_form_submit/model"
	"mvtable/internal/app/mv_form_submit/repo"
	model4 "mvtable/internal/app/mv_record/model"
	model2 "mvtable/internal/app/mv_view_form/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvFormSubmitService struct{}

func NewMvFormSubmitService() *MvFormSubmitService {
	return &MvFormSubmitService{}
}

func (s *MvFormSubmitService) SubmitForm(userId string, req *model.SubmitFormReq) error {
	var insertData *model.MvFormSubmit
	err := db.Transaction(func(tx *gorm.DB) error {
		viewForm, err := db.Get[model2.MvViewForm](tx, map[string]any{"id": req.FormID})
		if err != nil {
			log.Error("get mv view form error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
		if viewForm == nil {
			return errorx.New(errorx.ErrNotFound, "表单不存在或已被删除")
		}

		// 是否需要登录
		if !viewForm.EnableNoLogin && userId == "" {
			return errorx.New(errorx.ErrNoPermission, "请先登录后再提交表单")
		}

		// 提交是否达到总数限制
		if viewForm.EnableLimitCollect {
			submitForms, _, err := db.List[model.MvFormSubmit](tx, 0, 0, map[string]any{"view_form_id": viewForm.ID}, nil)
			if err != nil {
				log.Error("query mv form submit error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}

			if int64(len(submitForms)) >= viewForm.LimitCollectCount {
				return errorx.New(errorx.ErrOperationFailed, "表单提交已达上限，无法继续提交")
			}
		}

		// 提交是否达到限制
		if viewForm.EnableLimitSubmit && !viewForm.EnableNoLogin {
			switch viewForm.LimitSubmitType {
			case constants.FormSubmitOnce:
				submitForm, err := db.Get[model.MvFormSubmit](tx, map[string]any{"user_id": userId, "view_form_id": viewForm.ID})
				if err != nil {
					log.Error("get mv form submit error", zap.Error(err))
					return errorx.InternalServerError("操作失败")
				}
				if submitForm != nil {
					return errorx.New(errorx.ErrOperationFailed, "表单只允许提交一次，无法重复提交")
				}
			case constants.FormSubmitPerDay:
				// 查看今天是否已经提交过表单，如果已经提交过则不允许再次提交
				now := time.Now()
				zero := now.Truncate(24 * time.Hour).In(now.Location())
				submitFor, err := db.Get[model.MvFormSubmit](tx, map[string]any{"user_id": userId, "view_form_id": viewForm.ID, "created_at": []any{">=", zero}})
				if err != nil {
					log.Error("get mv form submit error", zap.Error(err))
					return errorx.InternalServerError("操作失败")
				}
				if submitFor != nil {
					return errorx.New(errorx.ErrOperationFailed, "表单每天只允许提交一次，无法重复提交")
				}
			default:
				return errorx.New(errorx.ErrOperationFailed, "表单配置错误，无法提交")
			}
		}

		// 插入表单数据到 mv_record 表中
		tableSchema, err := repo.GetTableSchemaByViewFormID(tx, viewForm.ID)
		if err != nil {
			log.Error("get table schema by view form id error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
		if tableSchema == nil {
			return errorx.New(errorx.ErrOperationFailed, "表单关联的表格不存在")
		}

		// 获取表单字段
		//fields, _, err := db.List[model3.MvField](tx, 0, 0, map[string]any{"table_schema_id": tableSchema.ID}, nil)
		//if err != nil {
		//	log.Error("get mv field error", zap.Error(err))
		//	return errorx.InternalServerError("操作失败")
		//}

		//formConfig := viewForm.Config

		if err := db.Create(tx, &model4.MvRecord{
			TableSchemaID: tableSchema.ID,
			CreatedBy:     userId,
			RowData:       req.FormData,
		}); err != nil {
			log.Error("create mv record error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}

		insertData = &model.MvFormSubmit{
			UserID:     userId,
			ViewFormID: viewForm.ID,
			RecordID:   "",
		}
		if err := db.Create(tx, insertData); err != nil {
			log.Error("create mv form submit error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
