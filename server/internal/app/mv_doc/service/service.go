package service

import (
	"mvtable/internal/app/mv_doc/model"
	userModel "mvtable/internal/app/user/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvDocService struct{}

func NewMvDocService() *MvDocService {
	return &MvDocService{}
}

// CreateMvDoc 创建文档
func (s *MvDocService) CreateMvDoc(userId string, req *model.CreateMvDocReq) (string, error) {
	doc := &model.MvDoc{
		ProjectID: req.ProjectID,
		FolderID:  req.FolderID,
		Name:      req.Name,
		Content:   req.Content,
		CreatedBy: userId,
		UpdatedBy: userId,
	}

	if err := db.Create(db.GetDB(), doc); err != nil {
		log.Error("create mv doc error", zap.Error(err))
		return "", errorx.InternalServerError("创建失败")
	}

	return doc.ID, nil
}

// UpdateMvDoc 更新文档
func (s *MvDocService) UpdateMvDoc(userId string, req *model.UpdateMvDocReq) error {
	doc, err := db.Get[model.MvDoc](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv doc error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if doc == nil {
		return errorx.New(errorx.ErrNotFound, "文档不存在")
	}

	updateFields := make([]string, 0)

	if req.Name != nil {
		doc.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.Content != nil {
		doc.Content = *req.Content
		updateFields = append(updateFields, "content")
	}
	if req.FolderID != nil {
		doc.FolderID = *req.FolderID
		updateFields = append(updateFields, "folder_id")
	}

	// 更新修改人
	doc.UpdatedBy = userId
	updateFields = append(updateFields, "updated_by")

	if len(updateFields) == 0 {
		return nil
	}

	if err = db.Update(db.GetDB(), doc, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv doc error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

// DeleteMvDoc 删除文档
func (s *MvDocService) DeleteMvDoc(req *model.DeleteMvDocReq) error {
	if err := db.Delete[model.MvDoc](db.GetDB(), map[string]any{"id": req.ID}); err != nil {
		log.Error("delete mv doc error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}

	return nil
}

// GetMvDoc 获取文档详情
func (s *MvDocService) GetMvDoc(req *model.GetMvDocReq) (*model.GetMvDocRes, error) {
	doc, err := db.Get[model.MvDoc](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv doc error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if doc == nil {
		return nil, errorx.New(errorx.ErrNotFound, "文档不存在")
	}

	// 获取创建者和更新者信息
	var creator, updater *model.UserInfo
	if doc.CreatedBy != "" {
		creator, err = s.getUserInfo(doc.CreatedBy)
		if err != nil {
			log.Error("get creator info error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
	}
	if doc.UpdatedBy != "" {
		updater, err = s.getUserInfo(doc.UpdatedBy)
		if err != nil {
			log.Error("get updater info error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
	}

	docItem := model.MvDocItem{
		ID:        doc.ID,
		ProjectID: doc.ProjectID,
		FolderID:  doc.FolderID,
		Name:      doc.Name,
		Content:   doc.Content,
		Creator:   creator,
		Updater:   updater,
		Ct:        doc.CreatedAt.Format("2006-01-02 15:04:05"),
		Ut:        doc.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return &model.GetMvDocRes{
		MvDocItem: docItem,
	}, nil
}

// ListMvDocs 查询文档列表
func (s *MvDocService) ListMvDocs(req *model.ListMvDocReq) (*model.ListMvDocRes, error) {
	query := make(map[string]any)

	if req.ProjectID != "" {
		query["project_id"] = req.ProjectID
	}
	if req.FolderID != "" {
		query["folder_id"] = req.FolderID
	}

	// 分页参数
	offset := (req.Page - 1) * req.Page
	limit := req.Size

	// 排序
	orderBy := []string{"created_at DESC"}

	docs, total, err := db.List[model.MvDoc](db.GetDB(), offset, limit, query, orderBy)
	if err != nil {
		log.Error("list mv docs error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 收集所有用户ID
	userIDSet := make(map[string]struct{})
	for _, doc := range docs {
		if doc != nil {
			if doc.CreatedBy != "" {
				userIDSet[doc.CreatedBy] = struct{}{}
			}
			if doc.UpdatedBy != "" {
				userIDSet[doc.UpdatedBy] = struct{}{}
			}
		}
	}

	// 批量查询用户信息
	userIDs := make([]string, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	users, _, err := db.List[userModel.User](db.GetDB(), 0, 0, map[string]any{}, []string{})
	if err != nil {
		log.Error("list users error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 构建用户ID到用户信息的映射
	userMap := make(map[string]*model.UserInfo)
	for _, user := range users {
		if user != nil {
			if _, exists := userIDSet[user.ID]; exists {
				userMap[user.ID] = &model.UserInfo{
					ID:       user.ID,
					Email:    user.Email,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
				}
			}
		}
	}

	items := make([]model.MvDocItem, len(docs))
	for i, doc := range docs {
		if doc == nil {
			continue
		}
		items[i] = model.MvDocItem{
			ID:        doc.ID,
			ProjectID: doc.ProjectID,
			FolderID:  doc.FolderID,
			Name:      doc.Name,
			Content:   doc.Content,
			Creator:   userMap[doc.CreatedBy],
			Updater:   userMap[doc.UpdatedBy],
			Ct:        doc.CreatedAt.Format("2006-01-02 15:04:05"),
			Ut:        doc.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &model.ListMvDocRes{
		List:  items,
		Total: total,
	}, nil
}

// getUserInfo 获取用户信息
func (s *MvDocService) getUserInfo(userID string) (*model.UserInfo, error) {
	if userID == "" {
		return nil, nil
	}

	// 查询用户信息
	user, err := db.Get[userModel.User](db.GetDB(), map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	// 构建返回结果
	return &model.UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}
