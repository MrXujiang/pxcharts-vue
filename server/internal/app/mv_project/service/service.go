package service

import (
	"encoding/json"
	dashboardModel "mvtable/internal/app/mv_dashboard/model"
	dashboardService "mvtable/internal/app/mv_dashboard/service"
	docModel "mvtable/internal/app/mv_doc/model"
	docService "mvtable/internal/app/mv_doc/service"
	fieldModel "mvtable/internal/app/mv_field/model"
	folderModel "mvtable/internal/app/mv_folder/model"
	"mvtable/internal/app/mv_project/model"
	"mvtable/internal/app/mv_project/repo"
	model4 "mvtable/internal/app/mv_project_favorite/model"
	model3 "mvtable/internal/app/mv_project_perm/model"
	model2 "mvtable/internal/app/mv_project_state/model"
	model8 "mvtable/internal/app/mv_record/model"
	model5 "mvtable/internal/app/mv_table_schema/model"
	tableSchemaService "mvtable/internal/app/mv_table_schema/service"
	model6 "mvtable/internal/app/mv_view/model"
	viewService "mvtable/internal/app/mv_view/service"
	model7 "mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvProjectService struct{}

func NewMvProjectService() *MvProjectService {
	return &MvProjectService{}
}

func (s *MvProjectService) CreateMvProject(userId string, req *model.CreateMvProjectReq) (*model.CreateMvProjectRes, error) {
	project := &model.MvProject{
		Name:        req.Name,
		Description: req.Description,
		UserID:      userId,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Create(tx, project); err != nil {
			log.Error("create mv project error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 初始化对应的分享配置和权限
		if err := db.Create(tx, &model2.MvProjectState{
			ProjectID:  project.ID,
			ShareRange: constants.ShareRangePrivate,
		}); err != nil {
			log.Error("create mv project state error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		if err := db.Create(tx, &model3.MvProjectPerm{
			ProjectID: project.ID,
			TargetID:  userId,
			Role:      constants.ProjectActionOwner,
		}); err != nil {
			log.Error("create mv project permission error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建一个默认表格、一个默认字段、一条默认数据、一个默认表格视图
		tableSchema := &model5.MvTableSchema{
			ProjectID: project.ID,
			Name:      "数据表",
			CreatedBy: userId,
		}
		if err := db.Create(tx, tableSchema); err != nil {
			log.Error("create default table schema error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		view := &model6.MvView{
			TableSchemaID: tableSchema.ID,
			Type:          constants.ViewTypeTable,
			Name:          "表格",
		}
		if err := db.Create(tx, view); err != nil {
			log.Error("create default view error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		viewTable := &model7.MvViewTable{
			ViewID: view.ID,
		}
		if err := db.Create(tx, viewTable); err != nil {
			log.Error("create default view table error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建一个默认字段
		fieldConfig, _ := json.Marshal(map[string]any{"width": 120, "isShow": true, "fixed": true})
		field := &fieldModel.MvField{
			TableSchemaID: tableSchema.ID,
			Title:         "标题",
			Type:          constants.MvFieldTypeText,
			Config:        fieldConfig,
			OrderIndex:    lexorank.MinString,
		}
		if err := db.Create(tx, field); err != nil {
			log.Error("create default field error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 创建一条默认数据
		rowData, _ := json.Marshal(map[string]any{
			field.Title: "示例数据",
		})
		record := &model8.MvRecord{
			TableSchemaID: tableSchema.ID,
			CreatedBy:     userId,
			RowData:       rowData,
			OrderIndex:    lexorank.MinString,
		}
		if err := db.Create(tx, record); err != nil {
			log.Error("create default record error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.CreateMvProjectRes{ID: project.ID}, nil
}

func (s *MvProjectService) UpdateMvProject(userId string, req *model.UpdateMvProjectReq) error {
	var (
		project      model.MvProject
		updateFields = make([]string, 0)
		err          error
	)

	if req.Name != nil {
		project.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.Description != nil {
		project.Description = *req.Description
		updateFields = append(updateFields, "description")
	}
	if req.EnableAdvancedPerm != nil {
		project.EnableAdvancedPerm = *req.EnableAdvancedPerm
		updateFields = append(updateFields, "enable_advanced_perm")
	}

	if !hasPermission(userId, req.ID) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}

	if err = db.Update(db.GetDB(), &project, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv project error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

func (s *MvProjectService) DeleteMvProject(userId string, req *model.DeleteMvProjectReq) error {
	if !hasPermission(userId, req.ID) {
		return errorx.New(errorx.ErrNoPermission, "无权限")
	}
	if err := db.Delete[model.MvProject](db.GetDB(), map[string]any{"id": req.ID}); err != nil {
		log.Error("delete mv project error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}
	return nil
}

func (s *MvProjectService) GetProject(userId string, req *model.GetProjectReq) (*model.GetProjectRes, error) {
	var (
		project    *model.MvProject
		isFavorite bool
		err        error
	)
	project, err = db.Get[model.MvProject](db.GetDB(), map[string]any{"id": req.ProjectID})
	if err != nil {
		log.Error("get mv project error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}
	if project == nil {
		return nil, errorx.New(errorx.ErrNotFound, "项目不存在")
	}

	// 判断用户是否收藏该项目
	favorite, err := db.Get[model4.MvProjectFavorite](db.GetDB(), map[string]any{"user_id": userId, "project_id": req.ProjectID})
	if err != nil {
		log.Error("get mv project favorite error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}
	if favorite != nil {
		isFavorite = true
	}

	// 设置最近访问时间
	if err = repo.UpsertRecentProject(db.GetDB(), userId, req.ProjectID); err != nil {
		log.Error("upsert recent project error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	return &model.GetProjectRes{
		MvProject:  *project,
		IsFavorite: isFavorite,
	}, nil
}

func (s *MvProjectService) QueryProject(userId string, req *model.QueryProjectReq) (*model.QueryProjectRes, error) {
	var (
		projects []*model.QueryProjectItem
		total    int64
		err      error
	)
	if req.Type == constants.QueryProjectRecently {
		projects, total, err = repo.GetRecentProjects(db.GetDB(), userId)
		if err != nil {
			log.Error("query recent projects error", zap.Error(err))
			return nil, errorx.InternalServerError("查询失败")
		}
	}

	if req.Type == constants.QueryProjectCreated {
		projects, total, err = repo.GetCreatedProjects(db.GetDB(), userId, req.Page, req.Size)
		if err != nil {
			log.Error("query created projects error", zap.Error(err))
			return nil, errorx.InternalServerError("查询失败")
		}
	}

	if req.Type == constants.QueryProjectShared {
		projects, total, err = repo.GetSharedProjects(db.GetDB(), userId, req.Page, req.Size)
		if err != nil {
			log.Error("query shared projects error", zap.Error(err))
			return nil, errorx.InternalServerError("查询失败")
		}
	}

	if req.Type == constants.QueryProjectFavorite {
		projects, total, err = repo.GetFavoriteProjects(db.GetDB(), userId, req.Page, req.Size)
		if err != nil {
			log.Error("query favorite projects error", zap.Error(err))
			return nil, errorx.InternalServerError("查询失败")
		}
	}

	// 判断用户是否收藏该项目
	for _, project := range projects {
		favorite, err := db.Get[model4.MvProjectFavorite](db.GetDB(), map[string]any{"user_id": userId, "project_id": project.ID})
		if err != nil {
			log.Error("get mv project favorite error", zap.Error(err))
			return nil, errorx.InternalServerError("查询失败")
		}
		if favorite != nil {
			project.IsFavorite = true
		}
	}

	return &model.QueryProjectRes{
		List:  projects,
		Total: total,
	}, nil
}

func (s *MvProjectService) SetFavoriteProject(userId string, req *model.SetFavoriteProjectReq) error {
	favorite, err := db.Get[model4.MvProjectFavorite](db.GetDB(), map[string]any{"user_id": userId, "project_id": req.ProjectID})

	if err != nil {
		log.Error("get mv project favorite error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	if req.IsFavorite && favorite != nil {
		return errorx.New(errorx.ErrOperationFailed, "已收藏")
	}

	if !req.IsFavorite && favorite == nil {
		return errorx.New(errorx.ErrOperationFailed, "已取消收藏")
	}

	if req.IsFavorite {
		// 收藏项目
		if err = db.Create(db.GetDB(), &model4.MvProjectFavorite{
			ProjectID: req.ProjectID,
			UserID:    userId,
		}); err != nil {
			log.Error("create mv project favorite error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
	} else {
		// 取消收藏
		if err = db.Delete[model4.MvProjectFavorite](db.GetDB(), map[string]any{"project_id": req.ProjectID, "user_id": userId}); err != nil {
			log.Error("create mv project favorite error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
	}

	return nil
}

// RenameProjectNode 通用重命名接口：对 table / form / dashboard / folder 进行重命名
func (s *MvProjectService) RenameProjectNode(userId string, req *model.RenameProjectNodeReq) error {
	switch req.Type {
	case "folder":
		// 先查询原文件夹，检查名称是否变化
		folder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get mv folder error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if folder == nil {
			return errorx.New(errorx.ErrNotFound, "文件夹不存在")
		}

		// 如果名称未变化，直接返回
		if folder.Name == req.Name {
			return nil
		}

		// 校验项目下是否已存在同名文件夹（排除自身）
		existFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{
			"project_id": req.ProjectID,
			"name":       req.Name,
		})
		if err != nil {
			log.Error("check exist folder error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if existFolder != nil && existFolder.ID != req.TargetID {
			return errorx.New(errorx.ErrAlreadyExists, "项目下已存在同名文件夹")
		}

		// 更新文件夹名称
		update := &folderModel.MvFolder{
			Name: req.Name,
		}
		if err := db.Update(db.GetDB(), update, map[string]any{"id": req.TargetID, "project_id": req.ProjectID}, "name"); err != nil {
			log.Error("rename folder node error", zap.Error(err), zap.String("targetId", req.TargetID))
			return errorx.InternalServerError("重命名失败")
		}
		return nil
	case "table", "form":
		// table 和 form 都是重命名对应的数据表（mv_table_schema）
		// 先查询原表，检查名称是否变化
		tableSchema, err := db.Get[model5.MvTableSchema](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get mv table schema error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if tableSchema == nil {
			return errorx.New(errorx.ErrNotFound, "数据表不存在")
		}

		// 如果名称未变化，直接返回
		if tableSchema.Name == req.Name {
			return nil
		}

		// 校验项目下是否已存在同名数据表（排除自身）
		existTable, err := db.Get[model5.MvTableSchema](db.GetDB(), map[string]any{
			"project_id": req.ProjectID,
			"name":       req.Name,
		})
		if err != nil {
			log.Error("check exist table schema error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if existTable != nil && existTable.ID != req.TargetID {
			return errorx.New(errorx.ErrAlreadyExists, "项目下已存在同名数据表")
		}

		// 更新数据表名称
		tableSvc := tableSchemaService.NewMvTableSchemaService()
		name := req.Name
		if err := tableSvc.UpdateMvTableSchema(&model5.UpdateMvTableSchemaReq{
			ID:   req.TargetID,
			Name: &name,
		}); err != nil {
			return err
		}
		return nil
	case "dashboard":
		// 先查询原仪表盘，检查名称是否变化
		dashboard, err := db.Get[dashboardModel.MvDashboard](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get mv dashboard error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if dashboard == nil {
			return errorx.New(errorx.ErrNotFound, "仪表盘不存在")
		}

		// 如果名称未变化，直接返回
		if dashboard.Name == req.Name {
			return nil
		}

		// 校验项目下是否已存在同名仪表盘（排除自身）
		existDashboard, err := db.Get[dashboardModel.MvDashboard](db.GetDB(), map[string]any{
			"project_id": req.ProjectID,
			"name":       req.Name,
		})
		if err != nil {
			log.Error("check exist dashboard error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if existDashboard != nil && existDashboard.ID != req.TargetID {
			return errorx.New(errorx.ErrAlreadyExists, "项目下已存在同名仪表盘")
		}

		// 更新仪表盘名称
		dashboardSvc := dashboardService.NewMvDashboardService()
		name := req.Name
		if err := dashboardSvc.UpdateMvDashboard(&dashboardModel.UpdateMvDashboardReq{
			ID:   req.TargetID,
			Name: &name,
		}); err != nil {
			return err
		}
		return nil
	case "doc":
		// 先查询原文档，检查名称是否变化
		doc, err := db.Get[docModel.MvDoc](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get mv doc error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if doc == nil {
			return errorx.New(errorx.ErrNotFound, "文档不存在")
		}

		// 如果名称未变化，直接返回
		if doc.Name == req.Name {
			return nil
		}

		// 校验项目下是否已存在同名文档（排除自身）
		existDoc, err := db.Get[docModel.MvDoc](db.GetDB(), map[string]any{
			"project_id": req.ProjectID,
			"name":       req.Name,
		})
		if err != nil {
			log.Error("check exist doc error", zap.Error(err))
			return errorx.InternalServerError("重命名失败")
		}
		if existDoc != nil && existDoc.ID != req.TargetID {
			return errorx.New(errorx.ErrAlreadyExists, "项目下已存在同名文档")
		}

		// 更新文档名称
		docSvc := docService.NewMvDocService()
		name := req.Name
		if err := docSvc.UpdateMvDoc(userId, &docModel.UpdateMvDocReq{
			ID:   req.TargetID,
			Name: &name,
		}); err != nil {
			return err
		}
		return nil
	default:
		return errorx.New(errorx.ErrInvalidParam, "重命名类型不合法")
	}
}

// CreateProjectNode 在项目下创建节点（文件夹 / 表格 / 表单 / 仪表盘）
// type: table（创建数据表及表格视图）、form（创建数据表及表单视图）、dashboard（创建仪表盘）、folder（创建文件夹）
func (s *MvProjectService) CreateProjectNode(userId string, req *model.CreateProjectNodeReq) (*model.CreateProjectNodeRes, error) {
	switch req.Type {
	case "folder":
		// 直接在当前项目下创建文件夹
		folder := &folderModel.MvFolder{
			ProjectID: req.ProjectID,
			Name:      req.Name,
			ParentID:  req.FolderID,
		}
		if err := db.Create(db.GetDB(), folder); err != nil {
			log.Error("create project node folder error", zap.Error(err))
			return nil, errorx.InternalServerError("创建失败")
		}
		return &model.CreateProjectNodeRes{
			ID:   folder.ID,
			Type: req.Type,
		}, nil
	case "table":
		// 创建数据表和默认表格视图
		tableSvc := tableSchemaService.NewMvTableSchemaService()
		tableID, err := tableSvc.CreateMvTableSchema(userId, &model5.CreateMvTableSchemaReq{
			ProjectID: req.ProjectID,
			FolderID:  req.FolderID,
			Name:      req.Name,
		})
		if err != nil {
			return nil, err
		}
		return &model.CreateProjectNodeRes{
			ID:   tableID,
			Type: req.Type,
		}, nil
	case "form":
		// 创建数据表 + 表单视图
		tableSvc := tableSchemaService.NewMvTableSchemaService()
		viewSvc := viewService.NewMvViewFormService()

		tableID, err := tableSvc.CreateMvTableSchema(userId, &model5.CreateMvTableSchemaReq{
			ProjectID: req.ProjectID,
			FolderID:  req.FolderID,
			Name:      req.Name,
		})
		if err != nil {
			return nil, err
		}

		viewID, err := viewSvc.CreateMvView(&model6.CreateMvViewReq{
			TableSchemaID: tableID,
			Type:          constants.ViewTypeForm,
		})
		if err != nil {
			return nil, err
		}

		return &model.CreateProjectNodeRes{
			ID:   viewID,
			Type: req.Type,
		}, nil
	case "dashboard":
		// 创建仪表盘
		dashboardSvc := dashboardService.NewMvDashboardService()
		dashboardID, err := dashboardSvc.CreateMvDashboard(&dashboardModel.CreateMvDashboardReq{
			ProjectID: req.ProjectID,
			FolderID:  req.FolderID,
			Name:      req.Name,
		})
		if err != nil {
			return nil, err
		}

		return &model.CreateProjectNodeRes{
			ID:   dashboardID,
			Type: req.Type,
		}, nil
	case "doc":
		// 创建文档
		docSvc := docService.NewMvDocService()
		docID, err := docSvc.CreateMvDoc(userId, &docModel.CreateMvDocReq{
			ProjectID: req.ProjectID,
			FolderID:  req.FolderID,
			Name:      req.Name,
		})
		if err != nil {
			return nil, err
		}

		return &model.CreateProjectNodeRes{
			ID:   docID,
			Type: req.Type,
		}, nil
	default:
		return nil, errorx.New(errorx.ErrInvalidParam, "类型不合法")
	}
}

// GetProjectTables 获取项目下的所有数据表
func (s *MvProjectService) GetProjectTables(userId string, req *model.GetProjectTablesReq) (*model.GetProjectTablesRes, error) {
	// 检查用户是否有权限访问该项目
	if !hasPermission(userId, req.ProjectID) {
		return nil, errorx.New(errorx.ErrNoPermission, "无权限")
	}

	// 查询项目下的所有数据表
	tables, total, err := db.List[model5.MvTableSchema](db.GetDB(), 0, 0, map[string]any{
		"project_id": req.ProjectID,
	}, []string{"created_at DESC"})

	if err != nil {
		log.Error("get project tables error", zap.Error(err), zap.String("projectId", req.ProjectID))
		return nil, errorx.InternalServerError("查询失败")
	}

	// 转换为响应格式
	tableItems := make([]*model.TableSchemaItem, len(tables))
	for i, table := range tables {
		tableItems[i] = &model.TableSchemaItem{
			ID:          table.ID,
			Name:        table.Name,
			Description: table.Description,
			FolderID:    table.FolderID,
		}
	}

	return &model.GetProjectTablesRes{
		List:  tableItems,
		Total: total,
	}, nil
}
