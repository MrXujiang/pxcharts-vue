package repo

import (
	"errors"
	"mvtable/internal/app/mv_rich_text_content/model"

	"gorm.io/gorm"
)

// GetRichTextContent 获取富文本内容
func GetRichTextContent(db *gorm.DB, recordID, fieldID string) (*model.MvRichTextContent, error) {
	var content model.MvRichTextContent
	err := db.Where("record_id = ? AND field_id = ? AND deleted_at IS NULL", recordID, fieldID).First(&content).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &content, nil
}

// SaveRichTextContent 保存或更新富文本内容
func SaveRichTextContent(db *gorm.DB, recordID, fieldID, content string) error {
	var existing model.MvRichTextContent
	err := db.Where("record_id = ? AND field_id = ? AND deleted_at IS NULL", recordID, fieldID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 创建新记录
		newContent := &model.MvRichTextContent{
			RecordID: recordID,
			FieldID:  fieldID,
			Content:  content,
		}
		return db.Create(newContent).Error
	} else if err != nil {
		return err
	}

	// 更新现有记录
	existing.Content = content
	return db.Save(&existing).Error
}

// DeleteRichTextContent 删除富文本内容（软删除）
func DeleteRichTextContent(db *gorm.DB, recordID, fieldID string) error {
	return db.Where("record_id = ? AND field_id = ?", recordID, fieldID).Delete(&model.MvRichTextContent{}).Error
}

// DeleteRichTextContentByRecordID 根据记录ID删除所有关联的富文本内容（软删除）
func DeleteRichTextContentByRecordID(db *gorm.DB, recordID string) error {
	return db.Where("record_id = ?", recordID).Delete(&model.MvRichTextContent{}).Error
}

// DeleteRichTextContentByRecordIDs 批量删除记录关联的富文本内容（软删除）
func DeleteRichTextContentByRecordIDs(db *gorm.DB, recordIDs []string) error {
	if len(recordIDs) == 0 {
		return nil
	}
	return db.Where("record_id IN ?", recordIDs).Delete(&model.MvRichTextContent{}).Error
}
