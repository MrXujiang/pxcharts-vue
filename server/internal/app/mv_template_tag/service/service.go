package service

import (
	"mvtable/internal/app/mv_template_tag/model"
	"mvtable/internal/app/mv_template_tag/repo"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvTemplateTagService struct{}

func NewMvTemplateTagService() *MvTemplateTagService {
	return &MvTemplateTagService{}
}

// QueryMvTemplateTag 查询模板标签列表
func (s *MvTemplateTagService) QueryMvTemplateTag(req *model.QueryMvTemplateTagReq) (*model.QueryMvTemplateTagRes, error) {
	list, err := repo.QueryTagList(req.SearchWord)
	if err != nil {
		log.Error("查询模板标签列表失败", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	return &model.QueryMvTemplateTagRes{
		List: list,
	}, nil
}

// CreateMvTemplateTag 创建模板标签
func (s *MvTemplateTagService) CreateMvTemplateTag(req *model.CreateMvTemplateTagReq) error {
	// 检查标签名称是否已存在
	existingTag, err := db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"name": req.Name})
	if err != nil {
		log.Error("查询模板标签失败", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}
	if existingTag != nil {
		return errorx.BadRequest("标签名称已存在")
	}

	// 获取最后一个标签的order_index
	var lastTag *model.MvTemplateTag
	tags, _, err := db.List[model.MvTemplateTag](db.GetDB(), 1, 1, map[string]any{}, []string{"order_index DESC"})
	if err != nil {
		log.Error("查询最后一个标签失败", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}
	if len(tags) > 0 {
		lastTag = tags[0]
	}

	// 计算新的order_index
	var orderIndex string
	if lastTag != nil {
		orderIndex = lexorank.Between(lastTag.OrderIndex, lexorank.MaxString)
	} else {
		orderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
	}

	// 创建标签
	tag := &model.MvTemplateTag{
		Name:        req.Name,
		Description: req.Description,
		OrderIndex:  orderIndex,
	}

	if err := db.Create(db.GetDB(), tag); err != nil {
		log.Error("创建模板标签失败", zap.Error(err))
		return errorx.InternalServerError("创建失败")
	}

	return nil
}

// UpdateMvTemplateTag 更新模板标签
func (s *MvTemplateTagService) UpdateMvTemplateTag(req *model.UpdateMvTemplateTagReq) error {
	// 检查标签是否存在
	tag, err := db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("查询模板标签失败", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if tag == nil {
		return errorx.BadRequest("标签不存在")
	}

	// 如果更新名称，检查新名称是否已存在
	if req.Name != nil && *req.Name != tag.Name {
		existingTag, err := db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"name": *req.Name})
		if err != nil {
			log.Error("查询模板标签失败", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
		if existingTag != nil {
			return errorx.BadRequest("标签名称已存在")
		}
	}

	// 构建更新数据
	updateData := &model.MvTemplateTag{}
	if req.Name != nil {
		updateData.Name = *req.Name
	}
	if req.Description != nil {
		updateData.Description = *req.Description
	}

	// 确定要更新的字段
	fields := []string{}
	if req.Name != nil {
		fields = append(fields, "name")
	}
	if req.Description != nil {
		fields = append(fields, "description")
	}

	if len(fields) == 0 {
		return errorx.BadRequest("没有需要更新的字段")
	}

	if err := db.Update(db.GetDB(), updateData, map[string]any{"id": req.ID}, fields...); err != nil {
		log.Error("更新模板标签失败", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

// DeleteMvTemplateTag 删除模板标签
func (s *MvTemplateTagService) DeleteMvTemplateTag(req *model.DeleteMvTemplateTagReq) error {
	// 检查标签是否存在
	tag, err := db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("查询模板标签失败", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}
	if tag == nil {
		return errorx.BadRequest("标签不存在")
	}

	// 软删除标签
	if err := db.Delete[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.ID}); err != nil {
		log.Error("删除模板标签失败", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}

	return nil
}

// UpdateTagSort 调整标签排序
func (s *MvTemplateTagService) UpdateTagSort(req *model.UpdateTagSortReq) error {
	var (
		err        error
		prevTag    *model.MvTemplateTag
		nextTag    *model.MvTemplateTag
		currentTag *model.MvTemplateTag
	)

	// 获取当前标签
	currentTag, err = db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.CurrentTagID})
	if err != nil {
		log.Error("查询当前标签失败", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if currentTag == nil {
		return errorx.BadRequest("当前标签不存在")
	}

	// 获取前一个标签（如果提供了ID）
	if req.PrevTagID != "" {
		prevTag, err = db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.PrevTagID})
		if err != nil {
			log.Error("查询前一个标签失败", zap.Error(err))
			return errorx.InternalServerError("获取失败")
		}
		if prevTag == nil {
			return errorx.New(errorx.ErrNotFound, "前一个标签不存在")
		}
	}

	// 获取下一个标签（如果提供了ID）
	if req.NextTagID != "" {
		nextTag, err = db.Get[model.MvTemplateTag](db.GetDB(), map[string]any{"id": req.NextTagID})
		if err != nil {
			log.Error("查询下一个标签失败", zap.Error(err))
			return errorx.InternalServerError("获取失败")
		}
		if nextTag == nil {
			return errorx.New(errorx.ErrNotFound, "下一个标签不存在")
		}
	}

	// 计算当前标签的排序索引
	if req.PrevTagID != "" && req.NextTagID != "" {
		currentTag.OrderIndex = lexorank.Between(prevTag.OrderIndex, nextTag.OrderIndex)
	}
	if req.PrevTagID != "" && req.NextTagID == "" {
		currentTag.OrderIndex = lexorank.Between(prevTag.OrderIndex, lexorank.MaxString)
	}
	if req.NextTagID != "" && req.PrevTagID == "" {
		currentTag.OrderIndex = lexorank.Between(lexorank.MinString, nextTag.OrderIndex)
	}

	// 更新标签的排序索引
	if err := db.Update(db.GetDB(), currentTag, map[string]any{"id": req.CurrentTagID}, "order_index"); err != nil {
		log.Error("更新标签排序失败", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}
