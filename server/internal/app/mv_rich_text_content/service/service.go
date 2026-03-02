package service

import (
	"mvtable/internal/app/mv_rich_text_content/repo"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvRichTextContentService struct{}

func NewMvRichTextContentService() *MvRichTextContentService {
	return &MvRichTextContentService{}
}

// GetRichTextContent 获取富文本内容
func (s *MvRichTextContentService) GetRichTextContent(recordID, fieldID string) (string, error) {
	content, err := repo.GetRichTextContent(db.GetDB(), recordID, fieldID)
	if err != nil {
		log.Error("get rich text content error", zap.Error(err), zap.String("recordId", recordID), zap.String("fieldId", fieldID))
		return "", errorx.InternalServerError("获取富文本内容失败")
	}
	if content == nil {
		return "", nil
	}
	return content.Content, nil
}

// SaveRichTextContent 保存或更新富文本内容
func (s *MvRichTextContentService) SaveRichTextContent(recordID, fieldID, content string) error {
	if err := repo.SaveRichTextContent(db.GetDB(), recordID, fieldID, content); err != nil {
		log.Error("save rich text content error", zap.Error(err), zap.String("recordId", recordID), zap.String("fieldId", fieldID))
		return errorx.InternalServerError("保存富文本内容失败")
	}
	return nil
}

// DeleteRichTextContent 删除富文本内容
func (s *MvRichTextContentService) DeleteRichTextContent(recordID, fieldID string) error {
	if err := repo.DeleteRichTextContent(db.GetDB(), recordID, fieldID); err != nil {
		log.Error("delete rich text content error", zap.Error(err), zap.String("recordId", recordID), zap.String("fieldId", fieldID))
		return errorx.InternalServerError("删除富文本内容失败")
	}
	return nil
}

// DeleteRichTextContentByRecordID 根据记录ID删除所有关联的富文本内容
func (s *MvRichTextContentService) DeleteRichTextContentByRecordID(recordID string) error {
	if err := repo.DeleteRichTextContentByRecordID(db.GetDB(), recordID); err != nil {
		log.Error("delete rich text content by record id error", zap.Error(err), zap.String("recordId", recordID))
		return errorx.InternalServerError("删除富文本内容失败")
	}
	return nil
}

// DeleteRichTextContentByRecordIDs 批量删除记录关联的富文本内容
func (s *MvRichTextContentService) DeleteRichTextContentByRecordIDs(recordIDs []string) error {
	if err := repo.DeleteRichTextContentByRecordIDs(db.GetDB(), recordIDs); err != nil {
		log.Error("delete rich text content by record ids error", zap.Error(err))
		return errorx.InternalServerError("删除富文本内容失败")
	}
	return nil
}
