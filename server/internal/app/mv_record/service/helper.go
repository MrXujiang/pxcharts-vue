package service

import (
	"encoding/json"
	"fmt"
	model2 "mvtable/internal/app/mv_field/model"
	"mvtable/internal/app/mv_record/model"
	"mvtable/internal/app/mv_record/repo"
	schemaModel "mvtable/internal/app/mv_table_schema/model"
	schemaService "mvtable/internal/app/mv_table_schema/service"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// handleDoubleAssociation 处理双向引用，保持双向一致性
func (s *MvRecordService) handleDoubleAssociation(tx *gorm.DB, field *model2.MvField, record *model.MvRecord, value any) error {
	currentRecordID := record.ID

	// 解析当前字段配置
	srcCfg, err := parseDoubleAssocConfig(field.Config)
	if err != nil {
		log.Error("parse double association config error", zap.Error(err), zap.String("fieldId", field.ID))
		return errorx.BadRequest("字段配置错误")
	}

	// 解析目标记录ID列表
	newIDs, err := normalizeAssocIDs(value)
	if err != nil {
		return errorx.BadRequest(err.Error())
	}
	if !srcCfg.AllowMultiSelect && len(newIDs) > 1 {
		return errorx.BadRequest("仅允许选择1条记录")
	}

	// 读取当前记录旧值，计算差异
	var rowData map[string]any
	if err := json.Unmarshal(record.RowData, &rowData); err != nil {
		return errorx.InternalServerError("解析记录数据失败")
	}
	oldIDs := extractIDs(rowData[field.Title])

	// 确保目标字段存在，如缺失则创建
	targetFieldID := srcCfg.TargetFieldID
	if targetFieldID == "" {
		targetFieldID, err = s.ensureTargetField(tx, field, srcCfg)
		if err != nil {
			return err
		}
		// 写回当前字段配置的 targetFieldId
		srcCfg.TargetFieldID = targetFieldID
		if err := s.updateFieldConfig(tx, field.ID, srcCfg); err != nil {
			return err
		}
	}

	// 获取目标字段，验证类型
	targetField, err := db.Get[model2.MvField](tx, map[string]any{"id": targetFieldID})
	if err != nil {
		return errorx.InternalServerError("获取关联字段失败")
	}
	if targetField == nil {
		return errorx.BadRequest("关联字段不存在")
	}
	if targetField.Type != constants.MvFieldTypeDoubleAssociation {
		return errorx.BadRequest("关联字段类型不匹配")
	}

	// 解析目标字段配置
	targetCfg, err := parseDoubleAssocConfig(targetField.Config)
	if err != nil {
		return errorx.BadRequest("关联字段配置错误")
	}
	// 若目标字段未记录反向字段ID，补充写回（双向缓存）
	if targetCfg.TargetFieldID == "" {
		targetCfg.TargetFieldID = field.ID
		if err := s.updateFieldConfig(tx, targetField.ID, targetCfg); err != nil {
			return err
		}
	}

	// 校验目标记录存在性
	if len(newIDs) > 0 {
		records, _, err := db.List[model.MvRecord](tx, 0, 0, map[string]any{"id": []any{"IN", newIDs}, "table_schema_id": srcCfg.TableSchemaID}, nil)
		if err != nil {
			return errorx.InternalServerError("查询关联记录失败")
		}
		if len(records) != len(newIDs) {
			return errorx.BadRequest("存在无效的关联记录ID")
		}
	}

	// 更新当前记录的值（始终数组，即便单选）
	if err := repo.SetFieldValue(tx, currentRecordID, field.ID, newIDs); err != nil {
		log.Error("set current record double assoc value error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	// 计算增量
	addSet := toSet(newIDs)
	oldSet := toSet(oldIDs)
	removeSet := difference(oldSet, addSet)
	// 需要处理的目标记录集合
	unionIDs := unionKeys(addSet, oldSet)

	if len(unionIDs) == 0 {
		return nil
	}

	// 拉取目标记录
	targetRecords, _, err := db.List[model.MvRecord](tx, 0, 0, map[string]any{"id": []any{"IN", unionIDs}, "table_schema_id": srcCfg.TableSchemaID}, nil)
	if err != nil {
		return errorx.InternalServerError("查询关联记录失败")
	}

	// 更新目标记录的反向字段
	for _, r := range targetRecords {
		if r == nil {
			continue
		}
		var tr map[string]any
		if err := json.Unmarshal(r.RowData, &tr); err != nil {
			return errorx.InternalServerError("解析关联记录数据失败")
		}
		existing := extractIDs(tr[targetField.Title])

		_, needAdd := addSet[r.ID]
		needRemove := removeSet[r.ID]
		updated := applyAddRemove(existing, targetCfg.AllowMultiSelect, currentRecordID, needAdd, needRemove)

		if err := repo.SetFieldValue(tx, r.ID, targetField.ID, updated); err != nil {
			log.Error("set target record double assoc value error", zap.Error(err), zap.String("targetRecordId", r.ID))
			return errorx.InternalServerError("更新失败")
		}
	}

	return nil
}

// ensureTargetField 确保目标表存在双向字段，若不存在则创建并返回字段ID
func (s *MvRecordService) ensureTargetField(tx *gorm.DB, srcField *model2.MvField, cfg doubleAssocConfig) (string, error) {
	// 获取源表信息，用于生成标题
	srcSchema, err := db.Get[schemaModel.MvTableSchema](tx, map[string]any{"id": srcField.TableSchemaID})
	if err != nil {
		return "", errorx.InternalServerError("获取源表失败")
	}
	if srcSchema == nil {
		return "", errorx.BadRequest("源表不存在")
	}
	baseTitle := srcSchema.Name
	if baseTitle == "" {
		baseTitle = "关联"
	}

	// 生成可用标题 A1/A2...
	var title string
	for i := 1; i < 1000; i++ {
		title = fmt.Sprintf("%s%d", baseTitle, i)
		exist, err := db.Get[model2.MvField](tx, map[string]any{"table_schema_id": cfg.TableSchemaID, "title": title})
		if err != nil {
			return "", errorx.InternalServerError("检查字段失败")
		}
		if exist == nil {
			break
		}
		title = ""
	}
	if title == "" {
		return "", errorx.InternalServerError("生成字段标题失败")
	}

	// 构造目标字段配置
	targetCfg := doubleAssocConfig{
		TableSchemaID:    srcField.TableSchemaID, // 反向指向源表
		AllowMultiSelect: cfg.AllowMultiSelect,
	}
	cfgBytes, _ := json.Marshal(targetCfg)

	newField := &model2.MvField{
		TableSchemaID: cfg.TableSchemaID,
		Title:         title,
		Type:          constants.MvFieldTypeDoubleAssociation,
		Config:        cfgBytes,
		OrderIndex:    lexorank.Between(lexorank.MinString, lexorank.MaxString),
	}

	if err := db.Create(tx, newField); err != nil {
		log.Error("create target double association field error", zap.Error(err))
		return "", errorx.InternalServerError("创建关联字段失败")
	}

	return newField.ID, nil
}

// updateFieldConfig 更新字段配置
func (s *MvRecordService) updateFieldConfig(tx *gorm.DB, fieldID string, cfg doubleAssocConfig) error {
	bytesCfg, err := json.Marshal(cfg)
	if err != nil {
		return errorx.InternalServerError("序列化字段配置失败")
	}
	update := &model2.MvField{
		Config: bytesCfg,
	}
	if err := db.Update(tx, update, map[string]any{"id": fieldID}, "config"); err != nil {
		return errorx.InternalServerError("更新字段配置失败")
	}
	return nil
}

// normalizeAssocIDs 将值规范为字符串数组
func normalizeAssocIDs(v any) ([]string, error) {
	switch val := v.(type) {
	case nil:
		return []string{}, nil
	case string:
		if val == "" {
			return []string{}, nil
		}
		return []string{val}, nil
	case []any:
		res := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok && s != "" {
				res = append(res, s)
			}
		}
		return unique(res), nil
	default:
		return nil, fmt.Errorf("关联值格式错误")
	}
}

// extractIDs 从任意值提取字符串数组
func extractIDs(v any) []string {
	switch val := v.(type) {
	case nil:
		return []string{}
	case string:
		if val == "" {
			return []string{}
		}
		return []string{val}
	case []string:
		return unique(val)
	case []any:
		res := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok && s != "" {
				res = append(res, s)
			}
		}
		return unique(res)
	default:
		return []string{}
	}
}

func unique(list []string) []string {
	seen := make(map[string]struct{}, len(list))
	out := make([]string, 0, len(list))
	for _, v := range list {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// applyAddRemove 根据是否需要添加/移除当前记录ID，返回更新后的ID列表
func applyAddRemove(existing []string, allowMulti bool, currentID string, needAdd bool, needRemove bool) []string {
	set := toSet(existing)
	if needRemove {
		delete(set, currentID)
	}
	if needAdd {
		set[currentID] = struct{}{}
	}

	// 单选：最多一个，优先保留 currentID
	if !allowMulti {
		if _, ok := set[currentID]; ok {
			return []string{currentID}
		}
		for k := range set {
			return []string{k}
		}
		return []string{}
	}

	// 保持原有顺序，追加 currentID
	res := make([]string, 0, len(set))
	for _, id := range existing {
		if _, ok := set[id]; ok {
			res = append(res, id)
			delete(set, id)
		}
	}
	if _, ok := set[currentID]; ok {
		res = append(res, currentID)
		delete(set, currentID)
	}
	for id := range set {
		res = append(res, id)
	}
	return res
}

// parseDoubleAssocConfig 解析双向引用配置
type doubleAssocConfig struct {
	TableSchemaID    string `json:"tableSchemaId"`
	AllowMultiSelect bool   `json:"allowMultiSelect"`
	Description      string `json:"description,omitempty"`
	TargetFieldID    string `json:"targetFieldId,omitempty"`
}

func parseDoubleAssocConfig(cfgBytes []byte) (doubleAssocConfig, error) {
	var cfg doubleAssocConfig
	if len(cfgBytes) == 0 {
		return cfg, fmt.Errorf("配置缺失")
	}
	if err := json.Unmarshal(cfgBytes, &cfg); err != nil {
		return cfg, err
	}
	// 默认值
	// AllowMultiSelect 默认 false（保持数组，但长度限制由调用方控制）
	return cfg, nil
}

// 辅助集合操作
func toSet(list []string) map[string]struct{} {
	res := make(map[string]struct{}, len(list))
	for _, v := range list {
		res[v] = struct{}{}
	}
	return res
}

func difference(a, b map[string]struct{}) map[string]bool {
	res := make(map[string]bool)
	for k := range a {
		if _, ok := b[k]; !ok {
			res[k] = true
		}
	}
	return res
}

func unionKeys(a map[string]struct{}, b map[string]struct{}) []string {
	res := make([]string, 0, len(a)+len(b))
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
		res = append(res, k)
	}
	for k := range b {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			res = append(res, k)
		}
	}
	return res
}

// setFieldDefaultValues 为记录设置字段默认值
func (s *MvRecordService) setFieldDefaultValues(fields []*model2.MvField, tableSchemaID string, userId string) (datatypes.JSON, error) {
	rowData := make(map[string]any)

	for _, field := range fields {
		// 获取字段配置
		fieldConfig, err := schemaService.GetFieldConfig(field.Config, field.Type, field.ID, field.Title)
		if err != nil {
			log.Warn("get field config error", zap.Error(err), zap.String("fieldId", field.ID), zap.Any("fieldType", field.Type))
			// 跳过有问题的字段，继续处理其他字段
			continue
		}

		// 设置字段默认值
		s.setDefaultValueFromConfig(fieldConfig, field, rowData, userId, tableSchemaID)
	}

	// 转换为 JSON
	jsonData, err := json.Marshal(rowData)
	if err != nil {
		return datatypes.JSON{}, fmt.Errorf("marshal row data failed: %w", err)
	}

	return datatypes.JSON(jsonData), nil
}

// setFieldValues 为记录设置字段值，支持传入rowData
func (s *MvRecordService) setFieldValues(fields []*model2.MvField, tableSchemaID string, userId string, inputRowData *datatypes.JSON) (datatypes.JSON, error) {
	rowData := make(map[string]any)

	// 解析传入的rowData（如果有的话）
	var inputData map[string]any
	if inputRowData != nil {
		if err := json.Unmarshal(*inputRowData, &inputData); err != nil {
			return datatypes.JSON{}, fmt.Errorf("unmarshal input row data failed: %w", err)
		}
	}

	for _, field := range fields {
		// 获取字段配置
		fieldConfig, err := schemaService.GetFieldConfig(field.Config, field.Type, field.ID, field.Title)
		if err != nil {
			log.Warn("get field config error", zap.Error(err), zap.String("fieldId", field.ID), zap.Any("fieldType", field.Type))
			// 跳过有问题的字段，继续处理其他字段
			continue
		}

		// 检查传入的rowData中是否有这个字段的值
		fieldID := field.ID
		if inputData != nil {
			if value, exists := inputData[fieldID]; exists {
				// 使用传入的值
				rowData[fieldID] = value
				continue
			}
		}

		// 没有传入值，使用默认值
		s.setDefaultValueFromConfig(fieldConfig, field, rowData, userId, tableSchemaID)
	}

	// 转换为 JSON
	jsonData, err := json.Marshal(rowData)
	if err != nil {
		return datatypes.JSON{}, fmt.Errorf("marshal row data failed: %w", err)
	}

	return datatypes.JSON(jsonData), nil
}

// setDefaultValueFromConfig 从字段配置中设置默认值
func (s *MvRecordService) setDefaultValueFromConfig(fieldConfig schemaModel.GeneralFieldConfig, field *model2.MvField, rowData map[string]any, tableSchemaId string, userId string) {
	switch cfg := fieldConfig.(type) {
	case *schemaModel.FieldTextConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldRadioConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldNumberConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldSelectMultipleConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldDateConfig:
		if dateConfig, ok := fieldConfig.(*schemaModel.FieldDateConfig); ok && dateConfig.Settings.IsDefaultCreateTime {
			// 日期字段且设置为默认创建时间，设置为当前时间
			rowData[field.ID] = time.Now()
		} else if field.Type == constants.MvFieldTypeCreator {
			// 设置创建人id
			rowData[field.ID] = userId
		} else {
			rowData[field.ID] = cfg.DefaultValue
		}
	case *schemaModel.FieldPersonConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldCurrencyConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldImageConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldAttachmentConfig:
		rowData[field.ID] = nil
	case *schemaModel.FieldProgressConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldLinkConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldFormulaConfig:
		rowData[field.ID] = nil
	case *schemaModel.FieldCheckboxConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldRateConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldEmailConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldIdCardConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldRichTextConfig:
		rowData[field.ID] = cfg.DefaultValue
	case *schemaModel.FieldFindRefConfig:
		// 查找引用不设置默认值
	case *schemaModel.FieldCreatorConfig:
		rowData[field.ID] = userId
	case *schemaModel.FieldSingleAssociationConfig:
		// 单向引用不设置默认值
	case *schemaModel.FieldDoubleAssociationConfig:
		// 双向引用不设置默认值
	case *schemaModel.FieldCreateTimeConfig:
		rowData[field.ID] = time.Now()
	case *schemaModel.FieldAutoNumberConfig:
		// 自动编号字段：查询该字段的最大值，然后+1
		maxValue, err := s.getMaxAutoNumberValue(tableSchemaId, field.ID)
		if err != nil {
			log.Warn("get max auto number value error", zap.Error(err), zap.String("tableSchemaId", tableSchemaId), zap.String("fieldId", field.ID))
			// 如果查询失败，使用默认值1
			rowData[field.ID] = 1
		} else {
			rowData[field.ID] = maxValue + 1
		}
	default:
		log.Error("no support field config type")
	}
}

// getMaxAutoNumberValue 获取指定自动编号字段的最大值
func (s *MvRecordService) getMaxAutoNumberValue(tableSchemaID, fieldID string) (int, error) {
	// 查询所有该表格的记录
	records, _, err := db.List[model.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, nil)
	if err != nil {
		return 0, fmt.Errorf("list records failed: %w", err)
	}

	maxValue := 0
	for _, record := range records {
		var rowData map[string]any
		if err := json.Unmarshal(record.RowData, &rowData); err != nil {
			continue // 跳过无法解析的记录
		}

		if value, exists := rowData[fieldID]; exists {
			if intValue, ok := value.(float64); ok {
				// JSON unmarshal 会将数字解析为float64
				if int(intValue) > maxValue {
					maxValue = int(intValue)
				}
			} else if intValue, ok := value.(int); ok {
				if intValue > maxValue {
					maxValue = intValue
				}
			}
		}
	}

	return maxValue, nil
}
