package service

import (
	"encoding/json"
	"fmt"
	model2 "mvtable/internal/app/mv_field/model"
	projectModel "mvtable/internal/app/mv_project/model"
	model3 "mvtable/internal/app/mv_record/model"
	richTextModel "mvtable/internal/app/mv_rich_text_content/model"
	richTextService "mvtable/internal/app/mv_rich_text_content/service"
	"mvtable/internal/app/mv_table_schema/model"
	model4 "mvtable/internal/app/mv_view/model"
	model6 "mvtable/internal/app/mv_view_board/model"
	viewFormModel "mvtable/internal/app/mv_view_form/model"
	model5 "mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvTableSchemaService struct {
	richTextContentService *richTextService.MvRichTextContentService
}

func NewMvTableSchemaService() *MvTableSchemaService {
	return &MvTableSchemaService{
		richTextContentService: richTextService.NewMvRichTextContentService(),
	}
}

func (s *MvTableSchemaService) CreateMvTableSchema(userId string, req *model.CreateMvTableSchemaReq) (string, error) {
	var (
		schema = &model.MvTableSchema{
			FolderID:  req.FolderID,
			ProjectID: req.ProjectID,
			Name:      req.Name,
			CreatedBy: userId,
		}
		err error
	)

	// 校验项目是否存在
	project, err := db.Get[projectModel.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get project error", zap.Error(err))
		return "", errorx.InternalServerError("获取项目信息失败")
	}
	if project == nil {
		return "", errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	if err = db.Transaction(func(tx *gorm.DB) error {
		if err = db.Create(tx, schema); err != nil {
			log.Error("create mv table schema error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建表初始字段
		firstOrderIndex := lexorank.Between(lexorank.MinString, lexorank.MaxString)
		secondOrderIndex := lexorank.Between(firstOrderIndex, lexorank.MaxString)
		fields := []*model2.MvField{
			{
				TableSchemaID: schema.ID,
				Title:         "标题",
				Type:          constants.MvFieldTypeText,
				Config:        []byte(`{"width": 120, "fixed": true, "isShow": true}`),
				OrderIndex:    firstOrderIndex,
			},
			{
				TableSchemaID: schema.ID,
				Title:         "内容",
				Type:          constants.MvFieldTypeText,
				Config:        []byte(`{"width": 120, "fixed": false, "isShow": true}`),
				OrderIndex:    secondOrderIndex,
			},
		}
		if err = db.CreateBatch(tx, fields); err != nil {
			log.Error("create initial fields error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建2条表初始数据
		firstOrderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
		secondOrderIndex = lexorank.Between(firstOrderIndex, lexorank.MaxString)
		records := []*model3.MvRecord{
			{
				TableSchemaID: schema.ID,
				RowData:       fmt.Appendf(nil, `{"%s": "标题1", "%s": "内容1"}`, fields[0].ID, fields[1].ID),
				CreatedBy:     userId,
				OrderIndex:    firstOrderIndex,
			},
			{
				TableSchemaID: schema.ID,
				RowData:       fmt.Appendf(nil, `{"%s": "标题2", "%s": "内容2"}`, fields[0].ID, fields[1].ID),
				CreatedBy:     userId,
				OrderIndex:    secondOrderIndex,
			},
		}

		// 创建默认视图
		view := &model4.MvView{
			TableSchemaID: schema.ID,
			Type:          constants.ViewTypeTable,
			Name:          "表格",
		}
		if err = db.Create(tx, view); err != nil {
			log.Error("create default view error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建默认视图配置
		if err = db.Create(tx, &model5.MvViewTable{
			ViewID: view.ID,
		}); err != nil {
			log.Error("create default view table config error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		if err = db.CreateBatch(tx, records); err != nil {
			log.Error("create initial records error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}
		return nil
	}); err != nil {
		return "", err
	}

	return schema.ID, nil
}

func (s *MvTableSchemaService) UpdateMvTableSchema(req *model.UpdateMvTableSchemaReq) error {
	var (
		schema       model.MvTableSchema
		updateFields = make([]string, 0)
		err          error
	)

	if req.Name != nil {
		schema.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.Description != nil {
		schema.Description = *req.Description
		updateFields = append(updateFields, "description")
	}
	if req.RowName != nil {
		schema.RowName = *req.RowName
		updateFields = append(updateFields, "row_name")
	}

	if err = db.Update(db.GetDB(), &schema, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv table schema error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	return nil
}

func (s *MvTableSchemaService) DeleteMvTableSchema(req *model.DeleteMvTableSchemaReq) error {
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Delete[model.MvTableSchema](tx, map[string]any{"id": req.ID}); err != nil {
			log.Error("delete mv table schema error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 删除表字段
		if err := db.Delete[model2.MvField](tx, map[string]any{"table_schema_id": req.ID}); err != nil {
			log.Error("delete mv field error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 删除表数据
		if err := db.Delete[model3.MvRecord](tx, map[string]any{"table_schema_id": req.ID}); err != nil {
			log.Error("delete mv record error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// CopyMvTableSchema 复制数据表
// range: 1=仅复制数据表结构，2=复制数据表结构和所有记录
func (s *MvTableSchemaService) CopyMvTableSchema(userId string, req *model.CopyMvTableSchemaReq) (*model.CopyMvTableSchemaRes, error) {
	if req.Range != 1 && req.Range != 2 {
		return nil, errorx.New(errorx.ErrInvalidParam, "复制范围不合法")
	}

	// 查询源数据表
	sourceSchema, err := db.Get[model.MvTableSchema](db.GetDB(), map[string]any{"id": req.SourceID, "project_id": req.ProjectID})
	if err != nil {
		log.Error("get source table schema error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}
	if sourceSchema == nil {
		return nil, errorx.New(errorx.ErrNotFound, "源数据表不存在")
	}

	// 校验项目下名称唯一
	existSchema, err := db.Get[model.MvTableSchema](db.GetDB(), map[string]any{
		"project_id": req.ProjectID,
		"name":       req.Name,
	})
	if err != nil {
		log.Error("check exist table schema error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}
	if existSchema != nil {
		return nil, errorx.New(errorx.ErrAlreadyExists, "项目下已存在同名数据表")
	}

	// 查询源字段
	sourceFields, _, err := db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": sourceSchema.ID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("list source fields error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}

	// 查询源视图
	sourceViews, _, err := db.List[model4.MvView](db.GetDB(), 0, 0, map[string]any{"table_schema_id": sourceSchema.ID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("list source views error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}

	// 查询源记录（range=2 时）
	var sourceRecords []*model3.MvRecord
	if req.Range == 2 {
		sourceRecords, _, err = db.List[model3.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": sourceSchema.ID}, []string{"order_index ASC"})
		if err != nil {
			log.Error("list source records error", zap.Error(err))
			return nil, errorx.InternalServerError("复制失败")
		}
	}

	var newSchemaID string

	// 使用事务执行复制
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建新表结构
		newSchema := &model.MvTableSchema{
			FolderID:    sourceSchema.FolderID,
			ProjectID:   sourceSchema.ProjectID,
			Name:        req.Name,
			Version:     1,
			CreatedBy:   userId,
			Description: sourceSchema.Description,
			Config:      sourceSchema.Config,
			Stats:       sourceSchema.Stats,
			RowName:     sourceSchema.RowName,
		}
		if err := db.Create(tx, newSchema); err != nil {
			log.Error("create copied table schema error", zap.Error(err))
			return errorx.InternalServerError("复制失败")
		}
		newSchemaID = newSchema.ID

		// 2. 复制字段
		if len(sourceFields) > 0 {
			newFields := make([]*model2.MvField, 0, len(sourceFields))
			for _, f := range sourceFields {
				if f == nil {
					continue
				}
				newField := &model2.MvField{
					TableSchemaID: newSchemaID,
					Title:         f.Title,
					Type:          f.Type,
					Config:        f.Config,
					OrderIndex:    f.OrderIndex,
					IsPrimaryKey:  f.IsPrimaryKey,
				}
				newFields = append(newFields, newField)
			}
			if err := db.CreateBatch(tx, newFields); err != nil {
				log.Error("create copied fields error", zap.Error(err))
				return errorx.InternalServerError("复制失败")
			}
		}

		// 3. 复制视图及其配置
		if len(sourceViews) > 0 {
			for _, v := range sourceViews {
				if v == nil {
					continue
				}
				newView := &model4.MvView{
					TableSchemaID: newSchemaID,
					Type:          v.Type,
					Name:          v.Name,
					Description:   v.Description,
					OrderIndex:    v.OrderIndex,
				}
				if err := db.Create(tx, newView); err != nil {
					log.Error("create copied view error", zap.Error(err))
					return errorx.InternalServerError("复制失败")
				}

				// 复制视图配置
				if v.Type == constants.ViewTypeTable {
					sourceViewTable, err := db.Get[model5.MvViewTable](tx, map[string]any{"view_id": v.ID})
					if err != nil {
						log.Error("get source view table config error", zap.Error(err))
						return errorx.InternalServerError("复制失败")
					}
					if sourceViewTable != nil {
						newViewTable := &model5.MvViewTable{
							ViewID:       newView.ID,
							FilterConfig: sourceViewTable.FilterConfig,
							GroupConfig:  sourceViewTable.GroupConfig,
							SortConfig:   sourceViewTable.SortConfig,
							RowHeight:    sourceViewTable.RowHeight,
							ColorConfig:  sourceViewTable.ColorConfig,
						}
						if err := db.Create(tx, newViewTable); err != nil {
							log.Error("create copied view table config error", zap.Error(err))
							return errorx.InternalServerError("复制失败")
						}
					}
				}

				if v.Type == constants.ViewTypeForm {
					sourceViewForm, err := db.Get[viewFormModel.MvViewForm](tx, map[string]any{"view_id": v.ID})
					if err != nil {
						log.Error("get source view form config error", zap.Error(err))
						return errorx.InternalServerError("复制失败")
					}
					if sourceViewForm != nil {
						newViewForm := &viewFormModel.MvViewForm{
							ViewID:                newView.ID,
							Name:                  sourceViewForm.Name,
							Description:           sourceViewForm.Description,
							Cover:                 sourceViewForm.Cover,
							Layout:                sourceViewForm.Layout,
							Stats:                 sourceViewForm.Stats,
							EnableSharing:         sourceViewForm.EnableSharing,
							EnableAnonymous:       sourceViewForm.EnableAnonymous,
							Filter:                sourceViewForm.Filter,
							FilterConfig:          sourceViewForm.FilterConfig,
							EnableNoLogin:         sourceViewForm.EnableNoLogin,
							EnableLimitSubmit:     sourceViewForm.EnableLimitSubmit,
							LimitSubmitType:       sourceViewForm.LimitSubmitType,
							EnableLimitCollect:    sourceViewForm.EnableLimitCollect,
							LimitCollectCount:     sourceViewForm.LimitCollectCount,
							EnableCycleRemind:     sourceViewForm.EnableCycleRemind,
							CycleRemindConfig:     sourceViewForm.CycleRemindConfig,
							EnableEditAfterSubmit: sourceViewForm.EnableEditAfterSubmit,
							Config:                sourceViewForm.Config,
						}
						if err := db.Create(tx, newViewForm); err != nil {
							log.Error("create copied view form config error", zap.Error(err))
							return errorx.InternalServerError("复制失败")
						}
					}
				}
			}
		}

		// 4. 复制记录及富文本（range=2）
		if req.Range == 2 && len(sourceRecords) > 0 {
			newRecords := make([]*model3.MvRecord, 0, len(sourceRecords))
			recordIDMap := make(map[string]string, len(sourceRecords)) // oldID -> newID

			for _, r := range sourceRecords {
				if r == nil {
					continue
				}
				newRecord := &model3.MvRecord{
					TableSchemaID: newSchemaID,
					CreatedBy:     r.CreatedBy,
					RowData:       r.RowData,
					OrderIndex:    r.OrderIndex,
				}
				newRecords = append(newRecords, newRecord)
			}

			if err := db.CreateBatch(tx, newRecords); err != nil {
				log.Error("create copied records error", zap.Error(err))
				return errorx.InternalServerError("复制失败")
			}

			// 构建 oldID -> newID 映射
			j := 0
			for _, r := range sourceRecords {
				if r == nil {
					continue
				}
				// 跳过对应的 nil newRecord（理论上不会为 nil）
				for j < len(newRecords) && newRecords[j] == nil {
					j++
				}
				if j >= len(newRecords) {
					break
				}
				recordIDMap[r.ID] = newRecords[j].ID
				j++
			}

			// 复制富文本内容
			for oldID, newID := range recordIDMap {
				contents, _, err := db.List[richTextModel.MvRichTextContent](tx, 0, 0, map[string]any{"record_id": oldID}, nil)
				if err != nil {
					log.Error("list rich text content error", zap.Error(err), zap.String("recordId", oldID))
					return errorx.InternalServerError("复制失败")
				}
				if len(contents) == 0 {
					continue
				}

				newContents := make([]*richTextModel.MvRichTextContent, 0, len(contents))
				for _, c := range contents {
					if c == nil {
						continue
					}
					newContent := &richTextModel.MvRichTextContent{
						RecordID: newID,
						FieldID:  c.FieldID,
						Content:  c.Content,
					}
					newContents = append(newContents, newContent)
				}
				if len(newContents) > 0 {
					if err := db.CreateBatch(tx, newContents); err != nil {
						log.Error("create copied rich text content error", zap.Error(err))
						return errorx.InternalServerError("复制失败")
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.CopyMvTableSchemaRes{ID: newSchemaID}, nil
}

func (s *MvTableSchemaService) GetTableData(req *model.GetTableDataReq) (*model.GetTableDataRes, error) {
	var (
		tableSchema *model.MvTableSchema
		view        *model4.MvView
		viewTable   *model5.MvViewTable
		viewBoard   *model6.MvViewBoard
		fields      []*model2.MvField
		records     []*model3.MvRecord
		err         error
	)

	view, err = db.Get[model4.MvView](db.GetDB(), map[string]any{"id": req.ViewID})
	if err != nil {
		log.Error("get view error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if view == nil {
		return nil, errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	tableSchema, err = db.Get[model.MvTableSchema](db.GetDB(), map[string]any{"id": view.TableSchemaID})
	if err != nil {
		log.Error("get table schema error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if tableSchema == nil {
		return nil, errorx.New(errorx.ErrNotFound, "表格不存在")
	}

	fields, _, err = db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchema.ID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get fields error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 获取视图配置（用于筛选和排序）
	switch view.Type {
	case constants.ViewTypeTable:
		viewTable, err = db.Get[model5.MvViewTable](db.GetDB(), map[string]any{"view_id": req.ViewID})
		if err != nil {
			log.Warn("get view table config error", zap.Error(err))
			// 视图配置不存在时不中断，继续获取所有记录
		}
	case constants.ViewTypeBoard:
		viewBoard, err = db.Get[model6.MvViewBoard](db.GetDB(), map[string]any{"view_id": req.ViewID})
		if err != nil {
			log.Warn("get view board config error", zap.Error(err))
			// 视图配置不存在时不中断，继续获取所有记录
		}
	}

	// 获取所有记录
	records, _, err = db.List[model3.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchema.ID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get records error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 应用筛选和排序
	records, err = s.filterAndSortRecords(records, fields, viewTable, viewBoard)
	if err != nil {
		log.Warn("filter and sort records error", zap.Error(err))
		// 筛选排序失败时使用原始记录
	}

	fieldsItems := make([]model.GeneralFieldConfig, len(fields))
	for i, v := range fields {
		var fieldConfig model.GeneralFieldConfig
		fieldConfig, err = GetFieldConfig(v.Config, v.Type, v.ID, v.Title)
		if err != nil {
			log.Error("get field config error", zap.Error(err), zap.Any("field", v))
			return nil, errorx.InternalServerError("获取失败")
		}
		fieldsItems[i] = fieldConfig

		if i == 0 {
			fieldsItems[i].GetBase().Fixed = true
		}
	}

	recordsItems := make([]model.RecordItem, len(records))
	for i, v := range records {
		rowData, err := s.GetRowData(v, fields)
		if err != nil {
			log.Error("get row data error", zap.Error(err), zap.Any("record", v))
			return nil, errorx.InternalServerError("获取失败")
		}
		recordsItems[i] = rowData
	}

	updatedBy := ""
	if tableSchema.UpdatedBy != nil {
		updatedBy = *tableSchema.UpdatedBy
	}

	var tableConfig model.TableConfig
	if err := json.Unmarshal(tableSchema.Config, &tableConfig); err != nil {
		log.Error("unmarshal table config error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	return &model.GetTableDataRes{
		TableSchemaID: tableSchema.ID,
		Name:          tableSchema.Name,
		Version:       tableSchema.Version,
		Ct:            tableSchema.CreatedAt,
		Ut:            tableSchema.UpdatedAt,
		CreatedBy:     tableSchema.CreatedBy,
		UpdatedBy:     updatedBy,
		Config:        tableConfig,
		Fields:        fieldsItems,
		Records:       recordsItems,
	}, nil
}

// GetRowData 处理行数据
func (s *MvTableSchemaService) GetRowData(record *model3.MvRecord, fields []*model2.MvField) (model.RecordItem, error) {
	var rowData model.RecordItem
	if err := json.Unmarshal(record.RowData, &rowData); err != nil {
		return nil, err
	}

	// 🔸 重要：先设置审计信息，让公式字段能访问到
	// 设置记录审计信息
	s.setRecordAuditInfo(&rowData, record)

	// 设置记录 ID
	rowData["rowId"] = record.ID

	for _, field := range fields {
		if err := s.processFieldValue(field, &rowData, record.ID); err != nil {
			return nil, err
		}
	}

	return rowData, nil
}
