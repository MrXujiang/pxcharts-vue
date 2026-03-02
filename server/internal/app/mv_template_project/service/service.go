package service

import (
	"encoding/json"
	model2 "mvtable/internal/app/mv_field/model"
	model3 "mvtable/internal/app/mv_folder/model"
	model4 "mvtable/internal/app/mv_project/model"
	model5 "mvtable/internal/app/mv_record/model"
	model6 "mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/app/mv_template_project/model"
	"mvtable/internal/app/mv_template_project/repo"
	model7 "mvtable/internal/app/mv_view/model"
	model8 "mvtable/internal/app/mv_view_form/model"
	model9 "mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvTemplateProjectService struct{}

func NewMvTemplateProjectService() *MvTemplateProjectService {
	return &MvTemplateProjectService{}
}

// QueryMvTemplateProject 查询模板项目列表
func (s *MvTemplateProjectService) QueryMvTemplateProject(req *model.QueryMvTemplateProjectReq) (*model.QueryMvTemplateProjectRes, error) {
	list, total, err := repo.QueryTemplateProjectList(req.Page, req.Size, req.Tag)
	if err != nil {
		log.Error("查询模板项目列表失败", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	return &model.QueryMvTemplateProjectRes{
		List:  list,
		Total: total,
	}, nil
}

// UpdateMvTemplateProject 更新模板项目
func (s *MvTemplateProjectService) UpdateMvTemplateProject(req *model.UpdateMvTemplateProjectReq) error {
	// 检查模板项目是否存在
	templateProject, err := db.Get[model.MvTemplateProject](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("查询模板项目失败", zap.Error(err))
		return errorx.InternalServerError("查询失败")
	}
	if templateProject == nil {
		return errorx.BadRequest("模板项目不存在")
	}

	// 构建更新数据
	updateData := &model.MvTemplateProject{}
	fields := []string{}

	if req.Name != nil {
		updateData.Name = *req.Name
		fields = append(fields, "name")
	}
	if req.Description != nil {
		updateData.Description = *req.Description
		fields = append(fields, "description")
	}
	if req.Cover != nil {
		updateData.Cover = *req.Cover
		fields = append(fields, "cover")
	}
	if req.Tags != nil {
		updateData.Tags = *req.Tags
		fields = append(fields, "tags")
	}

	if len(fields) == 0 {
		return errorx.BadRequest("没有需要更新的字段")
	}

	if err := db.Update(db.GetDB(), updateData, map[string]any{"id": req.ID}, fields...); err != nil {
		log.Error("更新模板项目失败", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

// DeleteMvTemplateProject 删除模板项目及其相关数据
func (s *MvTemplateProjectService) DeleteMvTemplateProject(req *model.DeleteMvTemplateProjectReq) error {
	// 检查模板项目是否存在
	templateProject, err := db.Get[model.MvTemplateProject](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("查询模板项目失败", zap.Error(err))
		return errorx.InternalServerError("查询失败")
	}
	if templateProject == nil {
		return errorx.BadRequest("模板项目不存在")
	}

	// 使用事务确保数据一致性
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. 查询模板项目下的所有表格（通过project_id来查找）
		var tables []*model6.MvTableSchema
		if err := tx.Model(&model6.MvTableSchema{}).
			Where("project_id = ?", req.ID).
			Find(&tables).Error; err != nil {
			log.Error("查询模板项目表格失败", zap.Error(err))
			return errorx.InternalServerError("查询失败")
		}

		tableIDs := make([]string, 0, len(tables))
		for _, table := range tables {
			tableIDs = append(tableIDs, table.ID)
		}

		// 2. 删除视图表格（通过view_id关联）
		if len(tableIDs) > 0 {
			// 先查询这些表格的所有视图
			var views []*model7.MvView
			if err := tx.Model(&model7.MvView{}).
				Where("table_schema_id IN ?", tableIDs).
				Find(&views).Error; err != nil {
				log.Error("查询视图失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}

			viewIDs := make([]string, 0, len(views))
			for _, view := range views {
				viewIDs = append(viewIDs, view.ID)
			}

			// 删除视图表单
			if len(viewIDs) > 0 {
				if err := tx.Where("view_id IN ?", viewIDs).Delete(&model8.MvViewForm{}).Error; err != nil {
					log.Error("删除视图表单失败", zap.Error(err))
					return errorx.InternalServerError("删除失败")
				}

				// 删除视图表格
				if err := tx.Where("view_id IN ?", viewIDs).Delete(&model9.MvViewTable{}).Error; err != nil {
					log.Error("删除视图表格失败", zap.Error(err))
					return errorx.InternalServerError("删除失败")
				}

				// 删除视图
				if err := tx.Where("table_schema_id IN ?", tableIDs).Delete(&model7.MvView{}).Error; err != nil {
					log.Error("删除视图失败", zap.Error(err))
					return errorx.InternalServerError("删除失败")
				}
			}

			// 3. 删除记录
			if err := tx.Where("table_schema_id IN ?", tableIDs).Delete(&model5.MvRecord{}).Error; err != nil {
				log.Error("删除记录失败", zap.Error(err))
				return errorx.InternalServerError("删除失败")
			}

			// 4. 删除字段
			if err := tx.Where("table_schema_id IN ?", tableIDs).Delete(&model2.MvField{}).Error; err != nil {
				log.Error("删除字段失败", zap.Error(err))
				return errorx.InternalServerError("删除失败")
			}

			// 5. 删除表格
			if err := tx.Where("id IN ?", tableIDs).Delete(&model6.MvTableSchema{}).Error; err != nil {
				log.Error("删除表格失败", zap.Error(err))
				return errorx.InternalServerError("删除失败")
			}
		}

		// 6. 删除模板项目
		if err := db.Delete[model.MvTemplateProject](tx, map[string]any{"id": req.ID}); err != nil {
			log.Error("删除模板项目失败", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// SaveAsTemplate 保存项目为模板
func (s *MvTemplateProjectService) SaveAsTemplate(userId string, req *model.SaveAsTemplateReq) (*model.SaveAsTemplateRes, error) {
	// 检查项目是否存在
	project, err := db.Get[model4.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("查询项目失败", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}
	if project == nil {
		return nil, errorx.BadRequest("项目不存在")
	}

	var templateProjectID string

	// 使用事务确保数据一致性
	err = db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建模板项目
		templateProject := &model.MvTemplateProject{
			UserID:      userId,
			Name:        req.Name,
			Description: req.Description,
			Cover:       req.Cover,
			Tags:        req.Tags,
		}
		if err := db.Create(tx, templateProject); err != nil {
			log.Error("创建模板项目失败", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}
		templateProjectID = templateProject.ID

		// 2. 复制表格（需要维护table_schema_id映射）
		// 查询项目下的所有表格（通过folder_id关联到项目的文件夹）
		// 先查询项目的所有文件夹ID
		folders, _, err := db.List[model3.MvFolder](tx, 0, 0, map[string]any{"project_id": req.ProjectID}, []string{"created_at ASC"})
		if err != nil {
			log.Error("查询文件夹失败", zap.Error(err))
			return errorx.InternalServerError("查询失败")
		}

		// 收集所有文件夹ID（包括空字符串，表示根目录）
		folderIDs := make([]string, 0, len(folders)+1)
		folderIDs = append(folderIDs, "")
		for _, folder := range folders {
			folderIDs = append(folderIDs, folder.ID)
		}

		// 查询属于这些文件夹的表格
		projectTables := make([]*model6.MvTableSchema, 0)
		for _, folderID := range folderIDs {
			tables, _, err := db.List[model6.MvTableSchema](tx, 0, 0, map[string]any{"folder_id": folderID}, []string{"created_at ASC"})
			if err != nil {
				log.Error("查询表格失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			projectTables = append(projectTables, tables...)
		}

		tableIDMap := make(map[string]string) // 原table_schema_id -> 新table_schema_id
		for _, table := range projectTables {
			newTable := &model6.MvTableSchema{
				FolderID:    "",                // 所有表格的folder_id都设置为空字符串
				ProjectID:   templateProjectID, // 设置project_id为模板项目ID
				Name:        table.Name,
				Version:     table.Version,
				CreatedBy:   userId,
				Description: table.Description,
				Config:      table.Config,
				Stats:       table.Stats,
				RowName:     table.RowName,
			}
			if err := db.Create(tx, newTable); err != nil {
				log.Error("创建表格失败", zap.Error(err))
				return errorx.InternalServerError("创建失败")
			}
			tableIDMap[table.ID] = newTable.ID
		}

		// 4. 复制字段（需要维护field_id映射）
		fieldIDMap := make(map[string]string) // 原field_id -> 新field_id
		for oldTableID, newTableID := range tableIDMap {
			fields, _, err := db.List[model2.MvField](tx, 0, 0, map[string]any{"table_schema_id": oldTableID}, []string{"order_index ASC"})
			if err != nil {
				log.Error("查询字段失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			for _, field := range fields {
				newField := &model2.MvField{
					TableSchemaID: newTableID,
					Title:         field.Title,
					Type:          field.Type,
					Config:        field.Config,
					OrderIndex:    field.OrderIndex,
				}
				if err := db.Create(tx, newField); err != nil {
					log.Error("创建字段失败", zap.Error(err))
					return errorx.InternalServerError("创建失败")
				}
				fieldIDMap[field.ID] = newField.ID
			}
		}

		// 5. 复制记录（需要更新RowData中的字段ID）
		for oldTableID, newTableID := range tableIDMap {
			records, _, err := db.List[model5.MvRecord](tx, 0, 0, map[string]any{"table_schema_id": oldTableID}, []string{"order_index ASC"})
			if err != nil {
				log.Error("查询记录失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			for _, record := range records {
				// 解析RowData JSON
				var rowDataMap map[string]any
				if err := json.Unmarshal(record.RowData, &rowDataMap); err != nil {
					log.Error("解析RowData失败", zap.Error(err))
					return errorx.InternalServerError("解析记录数据失败")
				}

				// 更新RowData中的字段ID（key从旧字段ID替换为新字段ID）
				newRowDataMap := make(map[string]any)
				for oldFieldID, value := range rowDataMap {
					newRowDataMap[fieldIDMap[oldFieldID]] = value
				}

				// 重新序列化为JSON
				newRowDataBytes, err := json.Marshal(newRowDataMap)
				if err != nil {
					log.Error("序列化RowData失败", zap.Error(err))
					return errorx.InternalServerError("序列化记录数据失败")
				}

				newRecord := &model5.MvRecord{
					TableSchemaID: newTableID,
					CreatedBy:     userId,
					RowData:       newRowDataBytes,
					OrderIndex:    record.OrderIndex,
				}
				if err := db.Create(tx, newRecord); err != nil {
					log.Error("创建记录失败", zap.Error(err))
					return errorx.InternalServerError("创建失败")
				}
			}
		}

		// 6. 复制视图（需要维护view_id映射）
		// 获取所有表格ID
		tableIDs := make([]string, 0, len(tableIDMap))
		for oldTableID := range tableIDMap {
			tableIDs = append(tableIDs, oldTableID)
		}

		// 查询属于这些表格的视图
		projectViews := make([]*model7.MvView, 0)
		for _, oldTableID := range tableIDs {
			views, _, err := db.List[model7.MvView](tx, 0, 0, map[string]any{"table_schema_id": oldTableID}, []string{"created_at ASC"})
			if err != nil {
				log.Error("查询视图失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			projectViews = append(projectViews, views...)
		}

		viewIDMap := make(map[string]string) // 原view_id -> 新view_id
		for _, view := range projectViews {
			newView := &model7.MvView{
				TableSchemaID: tableIDMap[view.TableSchemaID],
				Type:          view.Type,
				Name:          view.Name,
				Description:   view.Description,
				OrderIndex:    view.OrderIndex,
			}
			if err := db.Create(tx, newView); err != nil {
				log.Error("创建视图失败", zap.Error(err))
				return errorx.InternalServerError("创建失败")
			}
			viewIDMap[view.ID] = newView.ID
		}

		// 7. 复制视图表单
		for oldViewID, newViewID := range viewIDMap {
			viewForm, err := db.Get[model8.MvViewForm](tx, map[string]any{"view_id": oldViewID})
			if err != nil {
				log.Error("查询视图表单失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			if viewForm != nil {
				newViewForm := &model8.MvViewForm{
					ViewID:                newViewID,
					Name:                  viewForm.Name,
					Description:           viewForm.Description,
					Cover:                 viewForm.Cover,
					Layout:                viewForm.Layout,
					Stats:                 viewForm.Stats,
					EnableSharing:         viewForm.EnableSharing,
					EnableAnonymous:       viewForm.EnableAnonymous,
					Filter:                viewForm.Filter,
					FilterConfig:          viewForm.FilterConfig,
					EnableNoLogin:         viewForm.EnableNoLogin,
					EnableLimitSubmit:     viewForm.EnableLimitSubmit,
					LimitSubmitType:       viewForm.LimitSubmitType,
					EnableLimitCollect:    viewForm.EnableLimitCollect,
					LimitCollectCount:     viewForm.LimitCollectCount,
					EnableCycleRemind:     viewForm.EnableCycleRemind,
					CycleRemindConfig:     viewForm.CycleRemindConfig,
					EnableEditAfterSubmit: viewForm.EnableEditAfterSubmit,
					Config:                viewForm.Config,
				}
				if err := db.Create(tx, newViewForm); err != nil {
					log.Error("创建视图表单失败", zap.Error(err))
					return errorx.InternalServerError("创建失败")
				}
			}
		}

		// 8. 复制视图表格
		for oldViewID, newViewID := range viewIDMap {
			viewTable, err := db.Get[model9.MvViewTable](tx, map[string]any{"view_id": oldViewID})
			if err != nil {
				log.Error("查询视图表格失败", zap.Error(err))
				return errorx.InternalServerError("查询失败")
			}
			if viewTable != nil {
				newViewTable := &model9.MvViewTable{
					ViewID:       newViewID,
					FilterConfig: viewTable.FilterConfig,
					GroupConfig:  viewTable.GroupConfig,
					SortConfig:   viewTable.SortConfig,
					RowHeight:    viewTable.RowHeight,
					ColorConfig:  viewTable.ColorConfig,
				}
				if err := db.Create(tx, newViewTable); err != nil {
					log.Error("创建视图表格失败", zap.Error(err))
					return errorx.InternalServerError("创建失败")
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.SaveAsTemplateRes{ID: templateProjectID}, nil
}
