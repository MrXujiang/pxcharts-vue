package service

import (
	model4 "mvtable/internal/app/mv_dashboard/model"
	"mvtable/internal/app/mv_doc/model"
	folderModel "mvtable/internal/app/mv_folder/model"
	"mvtable/internal/app/mv_folder/repo"
	model2 "mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvFolderService struct{}

func NewMvFolderService() *MvFolderService {
	return &MvFolderService{}
}

func (s *MvFolderService) CreateMvFolder(req *folderModel.CreateMvFolderReq) error {
	if req.ParentID != "" {
		parentFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": req.ParentID})
		if err != nil {
			log.Error("get parent folder error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}
		if parentFolder == nil {
			return errorx.New(errorx.ErrNotFound, "父文件夹不存在")
		}
	}

	// 校验同一项目下是否已存在同名文件夹
	existFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{
		"project_id": req.ProjectID,
		"name":       req.Name,
	})
	if err != nil {
		log.Error("check exist mv folder error", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}
	if existFolder != nil {
		return errorx.New(errorx.ErrAlreadyExists, "同级文件夹名称已存在")
	}

	if err := db.Create(db.GetDB(), &folderModel.MvFolder{
		ProjectID: req.ProjectID,
		Name:      req.Name,
		ParentID:  req.ParentID,
	}); err != nil {
		log.Error("create mv folder error", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}
	return nil
}

func (s *MvFolderService) UpdateMvFolder(req *folderModel.UpdateMvFolderReq) error {
	// 先查出原文件夹，拿到 projectId / parentId 做同级重名校验
	folder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv folder error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if folder == nil {
		return errorx.New(errorx.ErrNotFound, "文件夹不存在")
	}

	// 如果名称未变化，直接返回
	if folder.Name == req.Name {
		return nil
	}

	// 校验同一项目下是否已存在同名文件夹（排除自身）
	existFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{
		"project_id": folder.ProjectID,
		"name":       req.Name,
	})
	if err != nil {
		log.Error("check exist mv folder error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if existFolder != nil && existFolder.ID != folder.ID {
		return errorx.New(errorx.ErrAlreadyExists, "同级文件夹名称已存在")
	}

	update := &folderModel.MvFolder{
		Name: req.Name,
	}
	if err := db.Update(db.GetDB(), update, map[string]any{"id": req.ID}, "name"); err != nil {
		log.Error("update mv folder error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	return nil
}

func (s *MvFolderService) DeleteMvFolder(req *folderModel.DeleteMvFolderReq) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Delete[folderModel.MvFolder](tx, map[string]any{"id": req.ID}); err != nil {
			log.Error("delete mv folder error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 查询所有子文件夹ID
		childrenIds, err := repo.GetChildrenIdsByParentId(tx, req.ID)
		if err != nil {
			log.Error("get children folder ids error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 当前文件夹 + 所有子文件夹
		folderIds := append(childrenIds, req.ID)

		// 删除所有相关数据（记录、字段、表结构等）
		if err := repo.DeleteDataByFolderIds(tx, folderIds); err != nil {
			log.Error("delete data by folder ids error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		// 删除所有视图及其配置，以及仪表盘和仪表盘图表
		if err := repo.DeleteViewAndDashboardByFolderIds(tx, folderIds); err != nil {
			log.Error("delete view_and_dashboard by folder ids error", zap.Error(err))
			return errorx.InternalServerError("删除失败")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *MvFolderService) Subquery(req *folderModel.QuerySubNodeReq) (*folderModel.QuerySubNodeRes, error) {
	var err error
	subNodes := make([]*folderModel.SubNode, 0)

	folders, _, err := db.List[folderModel.MvFolder](db.GetDB(), 0, 0, map[string]any{"parent_id": req.ParentID, "project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query folders error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	tableSchemas, _, err := db.List[model2.MvTableSchema](db.GetDB(), 0, 0, map[string]any{"folder_id": req.ParentID, "project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query table schemas error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	dashboards, _, err := db.List[model4.MvDashboard](db.GetDB(), 0, 0, map[string]any{"folder_id": req.ParentID, "project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query dashboards error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	docs, _, err := db.List[model.MvDoc](db.GetDB(), 0, 0, map[string]any{"folder_id": req.ParentID, "project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query docs error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	for _, folder := range folders {
		subNodes = append(subNodes, &folderModel.SubNode{
			ID:       folder.ID,
			Name:     folder.Name,
			Type:     "folder",
			ParentID: folder.ParentID,
			Children: make([]*folderModel.SubNode, 0),
		})
	}
	for _, tableSchema := range tableSchemas {
		subNodes = append(subNodes, &folderModel.SubNode{
			ID:       tableSchema.ID,
			Name:     tableSchema.Name,
			Type:     "table",
			ParentID: tableSchema.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		})
	}
	for _, dashboard := range dashboards {
		subNodes = append(subNodes, &folderModel.SubNode{
			ID:       dashboard.ID,
			Name:     dashboard.Name,
			Type:     "dashboard",
			ParentID: dashboard.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		})
	}
	for _, doc := range docs {
		subNodes = append(subNodes, &folderModel.SubNode{
			ID:       doc.ID,
			Name:     doc.Name,
			Type:     "doc",
			ParentID: doc.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		})
	}

	return &folderModel.QuerySubNodeRes{List: subNodes}, nil
}

// QueryAllNodes 查询项目下所有节点（树状）
func (s *MvFolderService) QueryAllNodes(req *folderModel.QueryAllNodeReq) (*folderModel.QuerySubNodeRes, error) {
	var (
		folders      []*folderModel.MvFolder
		tableSchemas []*model2.MvTableSchema
		dashboards   []*model4.MvDashboard
		docs         []*model.MvDoc
		err          error
	)

	// 查询项目下所有文件夹
	folders, _, err = db.List[folderModel.MvFolder](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query folders error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	// 查询项目下所有表格
	tableSchemas, _, err = db.List[model2.MvTableSchema](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query table schemas error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	// 查询项目下所有仪表盘
	dashboards, _, err = db.List[model4.MvDashboard](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query dashboards error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	// 查询项目下所有文档
	docs, _, err = db.List[model.MvDoc](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query docs error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	// 构建 parentId -> []*SubNode 映射
	childrenMap := make(map[string][]*folderModel.SubNode)

	for _, folder := range folders {
		node := &folderModel.SubNode{
			ID:       folder.ID,
			Name:     folder.Name,
			Type:     "folder",
			ParentID: folder.ParentID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := folder.ParentID
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	for _, ts := range tableSchemas {
		node := &folderModel.SubNode{
			ID:       ts.ID,
			Name:     ts.Name,
			Type:     "table",
			ParentID: ts.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := ts.FolderID
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	for _, dashboard := range dashboards {
		node := &folderModel.SubNode{
			ID:       dashboard.ID,
			Name:     dashboard.Name,
			Type:     "dashboard",
			ParentID: dashboard.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := dashboard.FolderID
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	for _, doc := range docs {
		node := &folderModel.SubNode{
			ID:       doc.ID,
			Name:     doc.Name,
			Type:     "doc",
			ParentID: doc.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := doc.FolderID
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 递归填充 Children 字段
	var attachChildren func(nodes []*folderModel.SubNode)
	attachChildren = func(nodes []*folderModel.SubNode) {
		for _, n := range nodes {
			if childList, ok := childrenMap[n.ID]; ok {
				n.Children = childList
				attachChildren(childList)
			}
		}
	}

	// 根节点：parent_id / folder_id 为空字符串
	rootNodes := childrenMap[""]
	if rootNodes == nil {
		rootNodes = make([]*folderModel.SubNode, 0)
	}
	attachChildren(rootNodes)

	return &folderModel.QuerySubNodeRes{List: rootNodes}, nil
}

func (s *MvFolderService) Search(req *folderModel.SearchReq) (*folderModel.SearchRes, error) {
	keyword := strings.TrimSpace(req.Keyword)
	if keyword == "" {
		return &folderModel.SearchRes{List: make([]*folderModel.SubNode, 0)}, nil
	}

	dbCli := db.GetDB()

	var (
		folders      []*folderModel.MvFolder
		tableSchemas []*model2.MvTableSchema
		dashboards   []*model4.MvDashboard
		docs         []*model.MvDoc
	)

	// 搜索匹配的文件夹
	if err := dbCli.Model(&folderModel.MvFolder{}).
		Where("project_id = ?", req.ProjectID).
		Where("deleted_at IS NULL").
		Where("name ILIKE ?", "%"+keyword+"%").
		Order("created_at DESC").
		Find(&folders).Error; err != nil {
		log.Error("search folders error", zap.Error(err))
		return nil, errorx.InternalServerError("搜索失败")
	}

	// 搜索匹配的表格
	if err := dbCli.Model(&model2.MvTableSchema{}).
		Where("project_id = ?", req.ProjectID).
		Where("deleted_at IS NULL").
		Where("name ILIKE ?", "%"+keyword+"%").
		Order("created_at DESC").
		Find(&tableSchemas).Error; err != nil {
		log.Error("search table schemas error", zap.Error(err))
		return nil, errorx.InternalServerError("搜索失败")
	}

	// 搜索匹配的仪表盘
	if err := dbCli.Model(&model4.MvDashboard{}).
		Where("project_id = ?", req.ProjectID).
		Where("deleted_at IS NULL").
		Where("name ILIKE ?", "%"+keyword+"%").
		Order("created_at DESC").
		Find(&dashboards).Error; err != nil {
		log.Error("search dashboards error", zap.Error(err))
		return nil, errorx.InternalServerError("搜索失败")
	}

	// 搜索匹配的文档
	if err := dbCli.Model(&model.MvDoc{}).
		Where("project_id = ?", req.ProjectID).
		Where("deleted_at IS NULL").
		Where("name ILIKE ?", "%"+keyword+"%").
		Order("created_at DESC").
		Find(&docs).Error; err != nil {
		log.Error("search docs error", zap.Error(err))
		return nil, errorx.InternalServerError("搜索失败")
	}

	// 收集所有相关的文件夹ID（包括匹配的文件夹、表格、仪表盘和文档所在的文件夹）
	folderIDSet := make(map[string]bool)
	for _, folder := range folders {
		folderIDSet[folder.ID] = true
		if folder.ParentID != "" {
			folderIDSet[folder.ParentID] = true
		}
	}
	for _, ts := range tableSchemas {
		if ts.FolderID != "" {
			folderIDSet[ts.FolderID] = true
		}
	}
	for _, dashboard := range dashboards {
		if dashboard.FolderID != "" {
			folderIDSet[dashboard.FolderID] = true
		}
	}
	for _, doc := range docs {
		if doc.FolderID != "" {
			folderIDSet[doc.FolderID] = true
		}
	}

	// 递归收集所有父文件夹ID
	var collectParentFolders func(folderID string)
	collectParentFolders = func(folderID string) {
		if folderID == "" {
			return
		}
		folder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": folderID})
		if err != nil || folder == nil {
			return
		}
		if folder.ParentID != "" {
			folderIDSet[folder.ParentID] = true
			collectParentFolders(folder.ParentID)
		}
	}

	// 收集所有父文件夹
	for folderID := range folderIDSet {
		collectParentFolders(folderID)
	}

	// 查询所有相关的文件夹（包括父文件夹路径上的所有文件夹）
	var folderIDs []string
	for id := range folderIDSet {
		folderIDs = append(folderIDs, id)
	}

	var allFolders []*folderModel.MvFolder
	if len(folderIDs) > 0 {
		if err := dbCli.Model(&folderModel.MvFolder{}).
			Where("project_id = ?", req.ProjectID).
			Where("deleted_at IS NULL").
			Where("id IN ?", folderIDs).
			Find(&allFolders).Error; err != nil {
			log.Error("query folders error", zap.Error(err))
			return nil, errorx.InternalServerError("搜索失败")
		}
	}

	// 构建 parentId -> []*SubNode 映射
	childrenMap := make(map[string][]*folderModel.SubNode)

	// 添加所有相关文件夹节点（包括匹配的文件夹和路径上的父文件夹）
	for _, folder := range allFolders {
		node := &folderModel.SubNode{
			ID:       folder.ID,
			Name:     folder.Name,
			Type:     "folder",
			ParentID: folder.ParentID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := folder.ParentID
		if parentID == "" {
			parentID = ""
		}
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 添加匹配的表格节点
	for _, ts := range tableSchemas {
		node := &folderModel.SubNode{
			ID:       ts.ID,
			Name:     ts.Name,
			Type:     "table",
			ParentID: ts.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := ts.FolderID
		if parentID == "" {
			parentID = ""
		}
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 添加匹配的仪表盘节点
	for _, dashboard := range dashboards {
		node := &folderModel.SubNode{
			ID:       dashboard.ID,
			Name:     dashboard.Name,
			Type:     "dashboard",
			ParentID: dashboard.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := dashboard.FolderID
		if parentID == "" {
			parentID = ""
		}
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 添加匹配的文档节点
	for _, doc := range docs {
		node := &folderModel.SubNode{
			ID:       doc.ID,
			Name:     doc.Name,
			Type:     "doc",
			ParentID: doc.FolderID,
			Children: make([]*folderModel.SubNode, 0),
		}
		parentID := doc.FolderID
		if parentID == "" {
			parentID = ""
		}
		childrenMap[parentID] = append(childrenMap[parentID], node)
	}

	// 递归填充 Children 字段
	var attachChildren func(nodes []*folderModel.SubNode)
	attachChildren = func(nodes []*folderModel.SubNode) {
		for _, n := range nodes {
			if childList, ok := childrenMap[n.ID]; ok && len(childList) > 0 {
				n.Children = childList
				attachChildren(childList)
			}
		}
	}

	// 根节点：parent_id / folder_id 为空字符串
	rootNodes := childrenMap[""]
	if rootNodes == nil {
		rootNodes = make([]*folderModel.SubNode, 0)
	}
	attachChildren(rootNodes)

	return &folderModel.SearchRes{List: rootNodes}, nil
}

// ListProjectFolders 查询项目下所有文件夹（扁平列表）
func (s *MvFolderService) ListProjectFolders(req *folderModel.QueryFolderListReq) (*folderModel.QueryFolderListRes, error) {
	folders, _, err := db.List[folderModel.MvFolder](db.GetDB(), 0, 0, map[string]any{"project_id": req.ProjectID}, nil)
	if err != nil {
		log.Error("query project folders error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	items := make([]*folderModel.FolderItem, 0, len(folders))
	for _, f := range folders {
		if f == nil {
			continue
		}
		items = append(items, &folderModel.FolderItem{
			ID:        f.ID,
			ProjectID: f.ProjectID,
			Name:      f.Name,
			ParentID:  f.ParentID,
		})
	}

	return &folderModel.QueryFolderListRes{List: items}, nil
}

// MoveNode 移动节点到其他文件夹
func (s *MvFolderService) MoveNode(req *folderModel.MoveNodeReq) error {
	// 校验目标文件夹（如果提供了）
	if req.TargetFolderID != "" {
		targetFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": req.TargetFolderID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get target folder error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		if targetFolder == nil {
			return errorx.New(errorx.ErrNotFound, "目标文件夹不存在")
		}
	}

	switch req.Type {
	case "folder":
		// 移动文件夹：更新 ParentID
		sourceFolder, err := db.Get[folderModel.MvFolder](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get source folder error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		if sourceFolder == nil {
			return errorx.New(errorx.ErrNotFound, "源文件夹不存在")
		}

		// 如果目标文件夹ID就是源文件夹ID，无需移动
		if req.TargetFolderID == sourceFolder.ParentID {
			return nil
		}

		// 防止移动到自己的子文件夹中（防止循环引用）
		if req.TargetFolderID != "" {
			// 检查目标文件夹是否是源文件夹的子文件夹
			childrenIds, err := repo.GetChildrenIdsByParentId(db.GetDB(), req.TargetID)
			if err != nil {
				log.Error("get children folder ids error", zap.Error(err))
				return errorx.InternalServerError("移动失败")
			}
			for _, childID := range childrenIds {
				if childID == req.TargetFolderID {
					return errorx.New(errorx.ErrInvalidParam, "不能移动到自己的子文件夹中")
				}
			}
			// 检查目标文件夹是否是源文件夹本身
			if req.TargetFolderID == req.TargetID {
				return errorx.New(errorx.ErrInvalidParam, "不能移动到自身")
			}
		}

		// 更新文件夹的 ParentID
		update := &folderModel.MvFolder{
			ParentID: req.TargetFolderID,
		}
		if err := db.Update(db.GetDB(), update, map[string]any{"id": req.TargetID, "project_id": req.ProjectID}, "parent_id"); err != nil {
			log.Error("move folder error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		return nil

	case "table", "form":
		// 移动数据表：更新 FolderID（table 和 form 都是 mv_table_schema）
		sourceTable, err := db.Get[model2.MvTableSchema](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get source table schema error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		if sourceTable == nil {
			return errorx.New(errorx.ErrNotFound, "源数据表不存在")
		}

		// 如果目标文件夹ID就是当前文件夹ID，无需移动
		if req.TargetFolderID == sourceTable.FolderID {
			return nil
		}

		// 更新数据表的 FolderID
		update := &model2.MvTableSchema{
			FolderID: req.TargetFolderID,
		}
		if err := db.Update(db.GetDB(), update, map[string]any{"id": req.TargetID, "project_id": req.ProjectID}, "folder_id"); err != nil {
			log.Error("move table schema error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		return nil

	case "dashboard":
		// 移动仪表盘：更新 FolderID
		sourceDashboard, err := db.Get[model4.MvDashboard](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get source dashboard error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		if sourceDashboard == nil {
			return errorx.New(errorx.ErrNotFound, "源仪表盘不存在")
		}

		// 如果目标文件夹ID就是当前文件夹ID，无需移动
		if req.TargetFolderID == sourceDashboard.FolderID {
			return nil
		}

		// 更新仪表盘的 FolderID
		update := &model4.MvDashboard{
			FolderID: req.TargetFolderID,
		}
		if err := db.Update(db.GetDB(), update, map[string]any{"id": req.TargetID, "project_id": req.ProjectID}, "folder_id"); err != nil {
			log.Error("move dashboard error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		return nil

	case "doc":
		// 移动文档：更新 FolderID
		sourceDoc, err := db.Get[model.MvDoc](db.GetDB(), map[string]any{"id": req.TargetID, "project_id": req.ProjectID})
		if err != nil {
			log.Error("get source doc error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		if sourceDoc == nil {
			return errorx.New(errorx.ErrNotFound, "源文档不存在")
		}

		// 如果目标文件夹ID就是当前文件夹ID，无需移动
		if req.TargetFolderID == sourceDoc.FolderID {
			return nil
		}

		// 更新文档的 FolderID
		update := &model.MvDoc{
			FolderID: req.TargetFolderID,
		}
		if err := db.Update(db.GetDB(), update, map[string]any{"id": req.TargetID, "project_id": req.ProjectID}, "folder_id"); err != nil {
			log.Error("move doc error", zap.Error(err))
			return errorx.InternalServerError("移动失败")
		}
		return nil

	default:
		return errorx.New(errorx.ErrInvalidParam, "节点类型不合法")
	}
}
