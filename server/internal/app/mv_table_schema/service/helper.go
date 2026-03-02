package service

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	model2 "mvtable/internal/app/mv_field/model"
	model3 "mvtable/internal/app/mv_record/model"
	"mvtable/internal/app/mv_table_schema/model"
	model6 "mvtable/internal/app/mv_view_board/model"
	model5 "mvtable/internal/app/mv_view_table/model"
	userModel "mvtable/internal/app/user/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/expr-lang/expr"
	"go.uber.org/zap"
)

// setDefaultValue 设置指针字段的默认值
func setDefaultValue[T any](ptr **T, defaultValue T) {
	if *ptr == nil {
		*ptr = &defaultValue
	}
}

// setDefaultValueAny 设置any类型的指针字段默认值
func setDefaultValueAny(ptr **any, defaultValue any) {
	if *ptr == nil {
		val := defaultValue
		*ptr = &val
	}
}

// stringPtr 返回字符串指针
func stringPtr(s string) *string {
	return &s
}

// intPtr 返回int指针
func intPtr(i int) *int {
	return &i
}

// newField 通用字段创建函数
func newField[T any](base *model.FieldBase, config []byte, defaultSetter func(*T)) (*T, error) {
	cfg := new(T)
	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}

	// 类型断言为 GeneralFieldConfig 接口
	if fieldConfig, ok := any(cfg).(model.GeneralFieldConfig); ok {
		fieldConfig.GetBase().ID = base.ID
		fieldConfig.GetBase().Type = base.Type
		fieldConfig.GetBase().Title = base.Title

		if fieldConfig.GetBase().Width == 0 {
			fieldConfig.GetBase().Width = base.Width
		}
	}

	defaultSetter(cfg)

	return cfg, nil
}

// GetFieldConfig 获取字段配置
func GetFieldConfig(config []byte, fieldType constants.FieldType, id, title string) (model.GeneralFieldConfig, error) {
	baseConfig := &model.FieldBase{
		ID:       id,
		Type:     fieldType,
		Title:    title,
		Width:    120,
		Fixed:    false,
		IsShow:   true,
		Label:    "",
		Desc:     "",
		IconName: "",
	}
	switch fieldType {
	case constants.MvFieldTypeText:
		return newField(baseConfig, config, func(cfg *model.FieldTextConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "文本"
			cfg.IconName = "formatVerticalAlignLeft"
		})
	case constants.MvFieldTypeRadio:
		return newField(baseConfig, config, func(cfg *model.FieldRadioConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "单选"
			cfg.IconName = "chevronDownCircle"
			//cfg.Settings.IsReferenced = false
			//cfg.Options = []model.RadioOption{
			//	{ID: "1", Label: "选项一", Color: "#FF6B6B"},
			//	{ID: "2", Label: "选项二", Color: "#4ECDC4"},
			//}
		})
	case constants.MvFieldTypeNumber:
		return newField(baseConfig, config, func(cfg *model.FieldNumberConfig) {
			//cfg.DefaultValue = "1234.0"
			cfg.Label = "数字"
			cfg.IconName = "artboard"
			//cfg.Settings.UseThousandSeparator = false
			//cfg.Settings.DisplayFormat = "decimal"
			//cfg.Settings.DecimalPlaces = 1
			//cfg.Settings.ThousandSeparator = "comma"
			//cfg.Settings.LargeNumberAbbreviation = "none"
			//cfg.Settings.DisallowNegative = false
		})
	case constants.MvFieldTypeSelectMultiple:
		return newField(baseConfig, config, func(cfg *model.FieldSelectMultipleConfig) {
			//cfg.DefaultValue = []string{}
			cfg.Label = "多选"
			cfg.IconName = "componentCheckbox"
			//cfg.Settings.IsReferenced = false
			//cfg.Options = []model.RadioOption{
			//	{ID: "1", Label: "选项一", Color: "#FF6B6B"},
			//	{ID: "2", Label: "选项二", Color: "#4ECDC4"},
			//}
		})
	case constants.MvFieldTypeDate:
		return newField(baseConfig, config, func(cfg *model.FieldDateConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "日期"
			cfg.IconName = "time"
			//cfg.Format = "YYYY-MM-DD"
			//cfg.Settings.IsDefaultCreateTime = false
		})
	case constants.MvFieldTypePerson:
		return newField(baseConfig, config, func(cfg *model.FieldPersonConfig) {
			//cfg.DefaultValue = []any{}
			cfg.Label = "人员"
			cfg.IconName = "user1"
			//cfg.Settings.ShowNickname = true
			//cfg.Settings.ShowExternalTag = true
			//cfg.Settings.AllowMultiple = true
		})
	case constants.MvFieldTypeCurrency:
		return newField(baseConfig, config, func(cfg *model.FieldCurrencyConfig) {
			//cfg.DefaultValue = "1234.0"
			cfg.Label = "货币"
			cfg.IconName = "currencyExchange"
			//cfg.Settings.UseThousandSeparator = false
			//cfg.Settings.Currency = "RMB"
			//cfg.Settings.DecimalPlaces = 1
			//cfg.Settings.ThousandSeparator = "comma"
			//cfg.Settings.LargeNumberAbbreviation = "none"
			//cfg.Settings.DisallowNegative = false
		})
	case constants.MvFieldTypeImage:
		return newField(baseConfig, config, func(cfg *model.FieldImageConfig) {
			//cfg.DefaultValue = []any{}
			cfg.Label = "图片"
			cfg.IconName = "image"
		})
	case constants.MvFieldTypeAttachment:
		return newField(baseConfig, config, func(cfg *model.FieldAttachmentConfig) {
			cfg.Label = "附件"
			cfg.IconName = "fileAttachment"
		})
	case constants.MvFieldTypeProgress:
		return newField(baseConfig, config, func(cfg *model.FieldProgressConfig) {
			//cfg.Label = "进度"
			cfg.IconName = "componentInput"
			cfg.Settings.NumberFormat = "numVal"
			//cfg.Settings.DecimalPlaces = 0
			//cfg.Settings.Color = "1"
			//cfg.Settings.EnableCustomProgress = false
		})
	case constants.MvFieldTypeLink:
		return newField(baseConfig, config, func(cfg *model.FieldLinkConfig) {
			//cfg.DefaultValue = model.FieldLinkDefaultValue{
			//	LinkTitle: "",
			//	LinkUrl:   "",
			//}
			cfg.Label = "链接"
			cfg.IconName = "link"
			//cfg.LinkTitle = ""
			//cfg.Settings.IsLinkTitle = false
		})
	case constants.MvFieldTypeFormula:
		return newField(baseConfig, config, func(cfg *model.FieldFormulaConfig) {
			cfg.Label = "公式"
			cfg.IconName = "functions1"
		})
	case constants.MvFieldTypeCheckbox:
		return newField(baseConfig, config, func(cfg *model.FieldCheckboxConfig) {
			//cfg.DefaultValue = 0
			cfg.Label = "复选框"
			cfg.IconName = "checkRectangle"
		})
	case constants.MvFieldTypeRate:
		return newField(baseConfig, config, func(cfg *model.FieldRateConfig) {
			//cfg.DefaultValue = nil
			cfg.Label = "评分"
			cfg.IconName = "star"
			//cfg.Settings.Icon = "star"
			//cfg.Settings.ScoreStart = 1
			//cfg.Settings.ScoreEnd = 5
		})
	case constants.MvFieldTypeEmail:
		return newField(baseConfig, config, func(cfg *model.FieldEmailConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "邮箱"
			cfg.IconName = "email"
		})
	case constants.MvFieldTypeIdCard:
		return newField(baseConfig, config, func(cfg *model.FieldIdCardConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "身份证"
			cfg.IconName = "idCard"
		})
	case constants.MvFieldTypeRichText:
		return newField(baseConfig, config, func(cfg *model.FieldRichTextConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "富文本"
			cfg.IconName = "textInitial"
		})
	case constants.MvFieldTypeFindRef:
		return newField(baseConfig, config, func(cfg *model.FieldFindRefConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "查找引用"
			cfg.IconName = "dataSearch"
			//cfg.Settings.MatchCondition = "all"
			//cfg.Settings.SourceRef = model.FieldSourceRef{
			//	TableSchemaID: "",
			//	FieldID:       "",
			//}
			//cfg.Settings.Conditions = []model.FieldFindCondition{
			//	{
			//		SourceTableFieldID: "",
			//		Operator:           "equal",
			//		TargetTableFieldID: "",
			//	},
			//}
		})
	case constants.MvFieldTypeCreator:
		return newField(baseConfig, config, func(cfg *model.FieldCreatorConfig) {
			cfg.Label = "创建人"
			cfg.IconName = "userAdd"
			//cfg.Settings.ShowNickname = true
			//cfg.Settings.ShowExternalTag = true
		})
	case constants.MvFieldTypeSingleAssociation:
		return newField(baseConfig, config, func(cfg *model.FieldSingleAssociationConfig) {
			cfg.Label = "单向关联"
			cfg.IconName = "swapRight"
			//cfg.Settings.SourceRef = model.FieldAssociationSourceRef{
			//	TableSchemaID: "",
			//}
			//cfg.Settings.Conditions.AllowMultiSelect = true
		})
	case constants.MvFieldTypeDoubleAssociation:
		return newField(baseConfig, config, func(cfg *model.FieldDoubleAssociationConfig) {
			cfg.Label = "双向关联"
			cfg.IconName = "swap"
			//cfg.Settings.SourceRef = model.FieldAssociationSourceRef{
			//	TableSchemaID: "",
			//}
			//cfg.Settings.Conditions.AllowMultiSelect = true
		})
	case constants.MvFieldTypeCreateTime:
		return newField(baseConfig, config, func(cfg *model.FieldCreateTimeConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "创建时间"
			cfg.IconName = "calendarPlus"
			//cfg.Format = "YYYY-MM-DD"
		})
	case constants.MvFieldTypeAutoNumber:
		return newField(baseConfig, config, func(cfg *model.FieldAutoNumberConfig) {
			//cfg.DefaultValue = ""
			cfg.Label = "自动编号"
			cfg.IconName = "arrowUp01"
			//cfg.Settings.IDType = "autoIncrement"
			//cfg.Settings.Rules = model.FieldAutoNumberRules{
			//	Number: 3,
			//	Text:   "",
			//}
		})
	default:
		return nil, fmt.Errorf("unsupported field type: %s", fieldType)

	}
}

// processFieldValue 处理单个字段的值，根据字段类型进行特殊处理
func (s *MvTableSchemaService) processFieldValue(field *model2.MvField, rowData *model.RecordItem, recordID string) error {
	// 如果是富文本类型，需要特殊处理
	switch field.Type {
	case constants.MvFieldTypeRichText:
		// 从富文本内容表获取完整内容
		content, err := s.richTextContentService.GetRichTextContent(recordID, field.ID)
		if err != nil {
			// 如果获取失败，设置为空字符串
			content = ""
		}
		// 去除 HTML 标签
		content = stripHtmlTags(content)
		// 截取前10个字符（考虑UTF-8编码）
		if utf8.RuneCountInString(content) > 10 {
			runes := []rune(content)
			content = string(runes[:10])
		}
		(*rowData)[field.ID] = content
	case constants.MvFieldTypeSingleAssociation, constants.MvFieldTypeDoubleAssociation:
		assocValue, err := s.resolveAssociationField(field, *rowData)
		if err != nil {
			log.Error("resolve association field error", zap.Error(err), zap.String("fieldId", field.ID), zap.Any("fieldType", field.Type))
			(*rowData)[field.ID] = nil
			return nil
		}
		(*rowData)[field.ID] = assocValue
	case constants.MvFieldTypeFindRef:
		fieldConfig := model.FieldFindRefConfig{}
		if err := json.Unmarshal(field.Config, &fieldConfig); err != nil {
			return err
		}
		findRefValue, err := s.resolveFindRefField(fieldConfig, *rowData)
		if err != nil {
			log.Error("resolve find ref field error", zap.Error(err), zap.String("fieldId", field.ID))
			(*rowData)[field.ID] = nil
			return nil
		}
		(*rowData)[field.ID] = findRefValue
	case constants.MvFieldTypeCreator:
		creatorValue, err := s.resolveCreatorField(*rowData, field.ID)
		if err != nil {
			log.Error("resolve creator field error", zap.Error(err), zap.String("fieldId", field.ID))
			(*rowData)[field.ID] = nil
			return nil
		}
		(*rowData)[field.ID] = creatorValue
	case constants.MvFieldTypeFormula:
		formulaValue, err := s.resolveFormulaField(field, *rowData)
		if err != nil {
			log.Error("resolve formula field error", zap.Error(err), zap.String("fieldId", field.ID))
			(*rowData)[field.ID] = nil
			return nil
		}
		(*rowData)[field.ID] = formulaValue
	}
	return nil
}

// resolveFindRefField 处理查找引用字段，根据配置的条件查找匹配的记录
func (s *MvTableSchemaService) resolveFindRefField(fieldConfig model.FieldFindRefConfig, rowData model.RecordItem) (any, error) {
	// 检查 SourceRef 是否配置
	if fieldConfig.Settings.SourceRef.TableSchemaID == "" {
		return nil, nil
	}

	targetTableID := fieldConfig.Settings.SourceRef.TableSchemaID
	targetFieldID := fieldConfig.Settings.SourceRef.FieldID

	// 检查是否有有效的查找条件
	if len(fieldConfig.Settings.Conditions) == 0 {
		return nil, nil
	}

	// 构建查询条件
	query := db.GetDB().Model(&model3.MvRecord{}).Where("table_schema_id = ?", targetTableID)

	// 获取匹配条件，默认是AND逻辑
	matchCondition := fieldConfig.Settings.MatchCondition

	// 收集有效的查询条件
	var whereConditions []string
	var whereArgs []any
	conditionsApplied := false

	for _, condition := range fieldConfig.Settings.Conditions {
		if condition.SourceTableFieldID == "" || condition.TargetTableFieldID == "" {
			continue // 跳过无效条件
		}

		sourceFieldID := condition.SourceTableFieldID
		targetTableFieldID := condition.TargetTableFieldID
		operator := condition.Operator

		// 获取源字段的值
		sourceValue, exists := rowData[sourceFieldID]
		if !exists {
			continue // 源字段不存在，跳过此条件
		}

		conditionsApplied = true

		// 根据操作符构建查询条件
		switch operator {
		case "equal":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' = ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "notEqual":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' != ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "contains":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' LIKE ?", targetTableFieldID))
			whereArgs = append(whereArgs, "%"+fmt.Sprintf("%v", sourceValue)+"%")
		case "notContains":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' NOT LIKE ?", targetTableFieldID))
			whereArgs = append(whereArgs, "%"+fmt.Sprintf("%v", sourceValue)+"%")
		case "null":
			whereConditions = append(whereConditions, fmt.Sprintf("(row_data->'%s') IS NULL OR (row_data->'%s' = 'null')", targetTableFieldID, targetTableFieldID))
		case "notNull":
			whereConditions = append(whereConditions, fmt.Sprintf("(row_data->'%s') IS NOT NULL AND (row_data->'%s' != 'null')", targetTableFieldID, targetTableFieldID))
		case "greaterThan":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' > ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "lessThan":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' < ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "greaterThanOrEqual":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' >= ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "lessThanOrEqual":
			whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' <= ?", targetTableFieldID))
			whereArgs = append(whereArgs, sourceValue)
		case "before", "after", "beforeOrEqual", "afterOrEqual":
			// 日期时间比较，特殊处理
			if dateValue, ok := s.toTime(sourceValue); ok {
				dateStr := dateValue.Format("2006-01-02 15:04:05")
				switch operator {
				case "before":
					whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' < ?", targetTableFieldID))
					whereArgs = append(whereArgs, dateStr)
				case "after":
					whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' > ?", targetTableFieldID))
					whereArgs = append(whereArgs, dateStr)
				case "beforeOrEqual":
					whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' <= ?", targetTableFieldID))
					whereArgs = append(whereArgs, dateStr)
				case "afterOrEqual":
					whereConditions = append(whereConditions, fmt.Sprintf("row_data->>'%s' >= ?", targetTableFieldID))
					whereArgs = append(whereArgs, dateStr)
				}
			}
		}
	}

	// 如果没有有效的条件，返回空
	if !conditionsApplied {
		return nil, nil
	}

	// 根据匹配条件构建最终查询
	if matchCondition == "any" && len(whereConditions) > 1 {
		// OR逻辑：使用复杂的WHERE子句
		fullCondition := strings.Join(whereConditions, " OR ")
		query = query.Where(fullCondition, whereArgs...)
	} else {
		// AND逻辑（默认）：逐个添加WHERE条件
		for i, condition := range whereConditions {
			query = query.Where(condition, whereArgs[i])
		}
	}

	// 查询所有匹配的记录
	var records []*model3.MvRecord
	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("query find ref records failed: %w", err)
	}

	// 如果没有找到匹配的记录，返回空
	if len(records) == 0 {
		return nil, nil
	}

	// 如果指定了目标字段，返回该字段的值列表
	if targetFieldID != "" {
		var values []any
		for _, record := range records {
			var rowDataMap map[string]any
			if err := json.Unmarshal(record.RowData, &rowDataMap); err != nil {
				continue
			}
			if value, exists := rowDataMap[targetFieldID]; exists && value != nil {
				values = append(values, value)
			}
		}
		return values, nil
	}

	// 如果没有指定目标字段，返回记录ID列表
	var recordIDs []string
	for _, record := range records {
		recordIDs = append(recordIDs, record.ID)
	}
	return recordIDs, nil
}

// resolveAssociationField 处理单向/双向引用字段，返回包含 recordId 和主键值的结果
func (s *MvTableSchemaService) resolveAssociationField(field *model2.MvField, rowData model.RecordItem) (any, error) {
	referencedRecordID, ok := rowData[field.ID]
	if !ok || referencedRecordID == nil {
		return nil, nil
	}

	// 获取字段配置
	fieldConfig, err := GetFieldConfig(field.Config, field.Type, field.ID, field.Title)
	if err != nil {
		return nil, err
	}

	// 根据字段类型获取引用的表 ID
	var tableSchemaID string
	switch cfg := fieldConfig.(type) {
	case *model.FieldSingleAssociationConfig:
		tableSchemaID = cfg.Settings.SourceRef.TableSchemaID
	case *model.FieldDoubleAssociationConfig:
		tableSchemaID = cfg.Settings.SourceRef.TableSchemaID
	default:
		return nil, fmt.Errorf("unsupported association field type: %s", field.Type)
	}

	if tableSchemaID == "" {
		return nil, fmt.Errorf("tableSchemaId not found in field config, fieldId: %s", field.ID)
	}

	// 获取被引用的记录ID（可能是字符串或数组）
	recordIDs := extractRecordIDs(referencedRecordID)
	if len(recordIDs) == 0 {
		return nil, nil
	}

	// 查询被引用表的主键字段
	referencedFields, _, err := db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		return nil, err
	}
	if len(referencedFields) == 0 {
		return nil, fmt.Errorf("no primary key field in referenced table, tableSchemaId: %s", tableSchemaID)
	}
	primaryKeyField := referencedFields[0]

	// 查询被引用的记录
	referencedRecords, _, err := db.List[model3.MvRecord](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", recordIDs}, "table_schema_id": tableSchemaID}, nil)
	if err != nil {
		return nil, err
	}

	// 记录 ID -> 主键值 映射
	valueByRecordID := make(map[string]any, len(referencedRecords))
	for _, record := range referencedRecords {
		if record == nil {
			continue
		}
		var recordRowData model.RecordItem
		if err := json.Unmarshal(record.RowData, &recordRowData); err != nil {
			log.Error("unmarshal referenced record row data error", zap.Error(err), zap.String("recordId", record.ID))
			continue
		}
		if primaryKeyValue, ok := recordRowData[primaryKeyField.ID]; ok {
			valueByRecordID[record.ID] = primaryKeyValue
		}
	}

	// 按传入的 recordIDs 顺序组装返回值，包含 recordId 与主键值
	assocResults := make([]map[string]any, 0, len(recordIDs))
	for _, id := range recordIDs {
		val, ok := valueByRecordID[id]
		if !ok || val == "" || val == nil {
			assocResults = append(assocResults, map[string]any{
				"recordId": id,
				"value":    "未命名记录",
			})
			continue
		}
		assocResults = append(assocResults, map[string]any{
			"recordId": id,
			"value":    val,
		})
	}

	// 设置结果（单选返回单个对象，多选返回数组）
	if len(assocResults) == 0 {
		return []any{}, nil
	}

	return assocResults, nil
}

// extractRecordIDs 将 row_data 中的引用值提取为字符串 ID 列表
func extractRecordIDs(referencedRecordID any) []string {
	var recordIDs []string

	switch val := referencedRecordID.(type) {
	case string:
		if val != "" {
			recordIDs = []string{val}
		}
	case []string:
		for _, id := range val {
			if id != "" {
				recordIDs = append(recordIDs, id)
			}
		}
	case []any:
		for _, id := range val {
			if idStr, ok := id.(string); ok && idStr != "" {
				recordIDs = append(recordIDs, idStr)
			}
		}
	}

	return recordIDs
}

// filterAndSortRecords 筛选和排序记录
func (s *MvTableSchemaService) filterAndSortRecords(
	records []*model3.MvRecord,
	fields []*model2.MvField,
	viewTable *model5.MvViewTable,
	viewBoard *model6.MvViewBoard,
) ([]*model3.MvRecord, error) {

	if viewTable == nil && viewBoard == nil {
		return records, nil
	}

	var err error
	filteredRecords := records

	// 根据视图类型获取配置
	var filterConfig, sortConfig []byte

	if viewTable != nil {
		// 表格视图配置
		filterConfig, _ = viewTable.FilterConfig.MarshalJSON()
		sortConfig, _ = viewTable.SortConfig.MarshalJSON()
	} else if viewBoard != nil {
		// 看板视图配置
		filterConfig, _ = viewBoard.FilterConfig.MarshalJSON()
		sortConfig, _ = viewBoard.SortConfig.MarshalJSON()
	}

	// 应用筛选条件
	if len(filterConfig) > 2 { // 非空数组
		filteredRecords, err = s.applyFilters(filteredRecords, fields, filterConfig)
		if err != nil {
			log.Warn("apply filters failed", zap.Error(err))
			// 筛选失败时返回所有记录，不中断流程
		}
	}

	// 应用排序
	if len(sortConfig) > 2 { // 非空数组
		err = s.applySorting(filteredRecords, fields, sortConfig)
		if err != nil {
			log.Warn("apply sorting failed", zap.Error(err))
			// 排序失败时保持原有顺序
		}
	}

	return filteredRecords, nil
}

// applyFilters 应用筛选条件
func (s *MvTableSchemaService) applyFilters(
	records []*model3.MvRecord,
	fields []*model2.MvField,
	filterConfig json.RawMessage,
) ([]*model3.MvRecord, error) {

	var filterConditions []model.FilterCondition
	if err := json.Unmarshal(filterConfig, &filterConditions); err != nil {
		return records, fmt.Errorf("unmarshal filter config failed: %w", err)
	}

	// 创建字段ID到字段的映射，提高查找效率
	fieldMap := make(map[string]*model2.MvField)
	for _, field := range fields {
		fieldMap[field.ID] = field
	}

	var filteredRecords []*model3.MvRecord

	for _, record := range records {
		match := true

		for _, condition := range filterConditions {
			field, exists := fieldMap[condition.FieldID]
			if !exists {
				// 字段不存在，跳过此条件
				continue
			}

			// 获取字段值
			rowDataBytes, _ := record.RowData.MarshalJSON()
			fieldValue, err := s.getFieldValue(rowDataBytes, field.ID)
			if err != nil {
				// 获取字段值失败，跳过此条件
				continue
			}

			// 检查是否匹配筛选条件
			if !s.matchFilterCondition(fieldValue, condition, field.Type) {
				match = false
				break
			}
		}

		if match {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords, nil
}

// applySorting 应用排序条件
func (s *MvTableSchemaService) applySorting(
	records []*model3.MvRecord,
	fields []*model2.MvField,
	sortConfig json.RawMessage,
) error {

	var sortConditions []model.SortCondition
	if err := json.Unmarshal(sortConfig, &sortConditions); err != nil {
		return fmt.Errorf("unmarshal sort config failed: %w", err)
	}

	if len(sortConditions) == 0 {
		return nil
	}

	// 创建字段ID到字段的映射
	fieldMap := make(map[string]*model2.MvField)
	for _, field := range fields {
		fieldMap[field.ID] = field
	}

	// 使用sort.SliceStable进行稳定排序
	sort.SliceStable(records, func(i, j int) bool {
		for _, condition := range sortConditions {
			field, exists := fieldMap[condition.FieldID]
			if !exists {
				continue
			}

			rowDataI, _ := records[i].RowData.MarshalJSON()
			rowDataJ, _ := records[j].RowData.MarshalJSON()
			valI, errI := s.getFieldValue(rowDataI, field.ID)
			valJ, errJ := s.getFieldValue(rowDataJ, field.ID)

			if errI != nil && errJ != nil {
				continue
			}
			if errI != nil {
				return condition.Order == "desc" // 有错误的值排在后面
			}
			if errJ != nil {
				return condition.Order == "asc" // 有错误的值排在后面
			}

			comparison := s.compareValues(valI, valJ, field.Type)
			if comparison != 0 {
				if condition.Order == "desc" {
					return comparison > 0
				}
				return comparison < 0
			}
		}
		// 如果所有条件都相等，保持原有顺序
		return false
	})

	return nil
}

// getFieldValue 从记录的rowData中获取指定字段的值
func (s *MvTableSchemaService) getFieldValue(rowData []byte, fieldID string) (any, error) {
	var data map[string]any
	if err := json.Unmarshal(rowData, &data); err != nil {
		return nil, err
	}

	value, exists := data[fieldID]
	if !exists {
		return nil, nil // 字段不存在，返回nil
	}

	return value, nil
}

// matchFilterCondition 检查字段值是否匹配筛选条件
func (s *MvTableSchemaService) matchFilterCondition(fieldValue any, condition model.FilterCondition, fieldType constants.FieldType) bool {
	// 先检查null和notNull操作符
	if condition.Operator == "null" {
		return s.isEmptyValue(fieldValue)
	}
	if condition.Operator == "notNull" {
		return !s.isEmptyValue(fieldValue)
	}

	// 对于空值，其他操作符都返回false
	if s.isEmptyValue(fieldValue) {
		return false
	}

	switch condition.Operator {
	case "equal":
		return s.equals(fieldValue, condition.Value)
	case "notEqual":
		return !s.equals(fieldValue, condition.Value)
	case "contains":
		return s.contains(fieldValue, condition.Value)
	case "notContains":
		return !s.contains(fieldValue, condition.Value)
	case "greaterThan":
		return s.compareNumbers(fieldValue, condition.Value) > 0
	case "lessThan":
		return s.compareNumbers(fieldValue, condition.Value) < 0
	case "greaterThanOrEqual":
		result := s.compareNumbers(fieldValue, condition.Value)
		return result > 0 || result == 0
	case "lessThanOrEqual":
		result := s.compareNumbers(fieldValue, condition.Value)
		return result < 0 || result == 0
	case "before", "after", "beforeOrEqual", "afterOrEqual":
		return s.compareDateTime(fieldValue, condition.Value, condition.Operator)
	default:
		return false
	}
}

// equals 比较两个值是否相等
func (s *MvTableSchemaService) equals(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// 尝试转换为相同类型进行比较
	aStr := s.toString(a)
	bStr := s.toString(b)
	return aStr == bStr
}

// contains 检查a是否包含b
func (s *MvTableSchemaService) contains(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	aStr := s.toString(a)
	bStr := s.toString(b)
	return strings.Contains(aStr, bStr)
}

// isEmptyValue 检查值是否为空（nil、空字符串或空数组）
func (s *MvTableSchemaService) isEmptyValue(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []int64:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	default:
		return false
	}
}

// getRelativeTimeRange 获取相对时间的时间范围
func (s *MvTableSchemaService) getRelativeTimeRange(relativeTime string, now time.Time) (start, end time.Time, ok bool) {
	switch relativeTime {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	case "tomorrow":
		start = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	case "yesterday":
		start = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1).Add(-time.Nanosecond)
	case "thisWeek":
		// 计算本周的开始（周一）
		daysFromMonday := int(now.Weekday()) - int(time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "lastWeek":
		// 计算上周的开始（周一）
		daysFromMonday := int(now.Weekday()) - int(time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday-7, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "nextWeek":
		// 计算下周的开始（周一）
		daysFromMonday := int(now.Weekday()) - int(time.Monday)
		if daysFromMonday < 0 {
			daysFromMonday += 7
		}
		start = time.Date(now.Year(), now.Month(), now.Day()-daysFromMonday+7, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case "thisMonth":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "lastMonth":
		start = time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "nextMonth":
		start = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case "thisYear":
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	case "lastYear":
		start = time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	case "nextYear":
		start = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	default:
		return time.Time{}, time.Time{}, false
	}
	return start, end, true
}

// resolveFormulaField 计算公式字段的值
func (s *MvTableSchemaService) resolveFormulaField(field *model2.MvField, rowData model.RecordItem) (any, error) {
	fieldConfig := model.FieldFormulaConfig{}
	if err := json.Unmarshal(field.Config, &fieldConfig); err != nil {
		return nil, fmt.Errorf("unmarshal formula config error: %w", err)
	}

	// 如果没有表达式，返回空值
	if fieldConfig.Expression == "" {
		return nil, nil
	}

	// 配置自定义函数
	env := map[string]any{}
	maps.Copy(env, rowData)

	// 添加数学函数
	env["SUM"] = s.formulaSum
	env["AVERAGE"] = s.formulaAverage
	env["MAX"] = s.formulaMax
	env["MIN"] = s.formulaMin
	env["ROUND"] = s.formulaRound

	// 添加文本函数
	env["CONCATENATE"] = s.formulaConcatenate
	env["UPPER"] = s.formulaUpper
	env["LOWER"] = s.formulaLower
	env["LEN"] = s.formulaLen
	env["TRIM"] = s.formulaTrim

	// 添加逻辑函数
	env["IF"] = s.formulaIf
	env["AND"] = s.formulaAnd
	env["OR"] = s.formulaOr
	env["NOT"] = s.formulaNot

	// 添加日期函数
	env["NOW"] = s.formulaNow
	env["TODAY"] = s.formulaToday
	env["YEAR"] = s.formulaYear
	env["MONTH"] = s.formulaMonth
	env["DATEIF"] = s.formulaDateif

	// 使用 expr 库计算表达式
	program, err := expr.Compile(fieldConfig.Expression, expr.Env(env))
	if err != nil {
		return nil, fmt.Errorf("compile formula expression error: %w", err)
	}

	result, err := expr.Run(program, env)
	if err != nil {
		return nil, fmt.Errorf("run formula expression error: %w", err)
	}

	return result, nil
}

// 数学函数
func (s *MvTableSchemaService) formulaSum(args ...any) any {
	sum := 0.0
	for _, arg := range args {
		switch v := arg.(type) {
		case []any:
			for _, item := range v {
				if num, ok := s.toFloat64(item); ok {
					sum += num
				}
			}
		default:
			if num, ok := s.toFloat64(arg); ok {
				sum += num
			}
		}
	}
	return sum
}

func (s *MvTableSchemaService) formulaAverage(args ...any) any {
	var sum float64
	var count int
	for _, arg := range args {
		switch v := arg.(type) {
		case []any:
			for _, item := range v {
				if num, ok := s.toFloat64(item); ok {
					sum += num
					count++
				}
			}
		default:
			if num, ok := s.toFloat64(arg); ok {
				sum += num
				count++
			}
		}
	}
	if count == 0 {
		return 0.0
	}
	return sum / float64(count)
}

func (s *MvTableSchemaService) formulaMax(args ...any) any {
	var max float64
	var hasValue bool
	for _, arg := range args {
		switch v := arg.(type) {
		case []any:
			for _, item := range v {
				if num, ok := s.toFloat64(item); ok {
					if !hasValue || num > max {
						max = num
						hasValue = true
					}
				}
			}
		default:
			if num, ok := s.toFloat64(arg); ok {
				if !hasValue || num > max {
					max = num
					hasValue = true
				}
			}
		}
	}
	if !hasValue {
		return nil
	}
	return max
}

func (s *MvTableSchemaService) formulaMin(args ...any) any {
	var min float64
	var hasValue bool
	for _, arg := range args {
		switch v := arg.(type) {
		case []any:
			for _, item := range v {
				if num, ok := s.toFloat64(item); ok {
					if !hasValue || num < min {
						min = num
						hasValue = true
					}
				}
			}
		default:
			if num, ok := s.toFloat64(arg); ok {
				if !hasValue || num < min {
					min = num
					hasValue = true
				}
			}
		}
	}
	if !hasValue {
		return nil
	}
	return min
}

func (s *MvTableSchemaService) formulaRound(value any, decimals any) any {
	val, valOk := s.toFloat64(value)
	if !valOk {
		return nil
	}

	dec := int64(0)
	if decVal, ok := s.toFloat64(decimals); ok {
		dec = int64(decVal)
	}

	multiplier := math.Pow(10, float64(dec))
	return math.Round(val*multiplier) / multiplier
}

// 文本函数
func (s *MvTableSchemaService) formulaConcatenate(args ...any) any {
	var result strings.Builder
	for _, arg := range args {
		if str := s.toString(arg); str != "" {
			result.WriteString(str)
		}
	}
	return result.String()
}

func (s *MvTableSchemaService) formulaUpper(value any) any {
	str := s.toString(value)
	return strings.ToUpper(str)
}

func (s *MvTableSchemaService) formulaLower(value any) any {
	str := s.toString(value)
	return strings.ToLower(str)
}

func (s *MvTableSchemaService) formulaLen(value any) any {
	str := s.toString(value)
	return len([]rune(str)) // 使用 rune 计数支持 Unicode
}

func (s *MvTableSchemaService) formulaTrim(value any) any {
	str := s.toString(value)
	return strings.TrimSpace(str)
}

// 逻辑函数
func (s *MvTableSchemaService) formulaIf(condition any, trueValue any, falseValue any) any {
	if s.toBool(condition) {
		return trueValue
	}
	return falseValue
}

func (s *MvTableSchemaService) formulaAnd(args ...any) any {
	for _, arg := range args {
		if !s.toBool(arg) {
			return false
		}
	}
	return true
}

func (s *MvTableSchemaService) formulaOr(args ...any) any {
	for _, arg := range args {
		if s.toBool(arg) {
			return true
		}
	}
	return false
}

func (s *MvTableSchemaService) formulaNot(value any) any {
	return !s.toBool(value)
}

// 日期函数
func (s *MvTableSchemaService) formulaNow() any {
	return time.Now()
}

func (s *MvTableSchemaService) formulaToday() any {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

func (s *MvTableSchemaService) formulaYear(date any) any {
	if t, ok := s.toTime(date); ok {
		return t.Year()
	}
	return nil
}

func (s *MvTableSchemaService) formulaMonth(date any) any {
	if t, ok := s.toTime(date); ok {
		return int(t.Month())
	}
	return nil
}

func (s *MvTableSchemaService) formulaDateif(startDate any, endDate any, unit string) any {
	start, startOk := s.toTime(startDate)
	end, endOk := s.toTime(endDate)
	if !startOk || !endOk {
		return nil
	}

	duration := end.Sub(start)

	switch strings.ToUpper(unit) {
	case "Y": // 年
		years := end.Year() - start.Year()
		if end.YearDay() < start.YearDay() {
			years--
		}
		return years
	case "M": // 月
		months := (end.Year()-start.Year())*12 + int(end.Month()-start.Month())
		if end.Day() < start.Day() {
			months--
		}
		return months
	case "D": // 天
		return int(duration.Hours() / 24)
	case "H": // 小时
		return int(duration.Hours())
	case "N": // 分钟
		return int(duration.Minutes())
	case "S": // 秒
		return int(duration.Seconds())
	default:
		return nil
	}
}

// 辅助函数

func (s *MvTableSchemaService) toBool(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return v != "" && v != "false" && v != "0"
	case int, int64, int32, float64, float32:
		if num, ok := s.toFloat64(v); ok {
			return num != 0
		}
	}
	return false
}

// compareDateTime 比较日期时间
func (s *MvTableSchemaService) compareDateTime(a, b any, operator string) bool {
	aTime, aOk := s.toTime(a)

	if !aOk {
		return false
	}

	// 尝试将 b 转换为时间
	bTime, bOk := s.toTime(b)
	if bOk {
		// b 是普通的时间值
		switch operator {
		case "before":
			return aTime.Before(bTime)
		case "after":
			return aTime.After(bTime)
		case "beforeOrEqual":
			return aTime.Before(bTime) || aTime.Equal(bTime)
		case "afterOrEqual":
			return aTime.After(bTime) || aTime.Equal(bTime)
		default:
			return false
		}
	}

	// b 可能是相对时间值，尝试解析
	bStr, ok := b.(string)
	if !ok {
		return false
	}

	now := time.Now()
	rangeStart, rangeEnd, ok := s.getRelativeTimeRange(bStr, now)
	if !ok {
		return false
	}

	// 根据操作符比较时间范围
	switch operator {
	case "before":
		// a 在范围结束之前
		return aTime.Before(rangeEnd)
	case "after":
		// a 在范围开始之后
		return aTime.After(rangeStart)
	case "beforeOrEqual":
		// a 在范围结束之前或等于范围结束
		return aTime.Before(rangeEnd) || aTime.Equal(rangeEnd)
	case "afterOrEqual":
		// a 在范围开始之后或等于范围开始
		return aTime.After(rangeStart) || aTime.Equal(rangeStart)
	default:
		return false
	}
}

// compareValues 比较两个值，用于排序
func (s *MvTableSchemaService) compareValues(a, b any, fieldType constants.FieldType) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// 根据字段类型进行比较
	switch fieldType {
	case constants.MvFieldTypeNumber, constants.MvFieldTypeCurrency:
		return s.compareNumbers(a, b)
	case constants.MvFieldTypeDate:
		return s.compareDates(a, b)
	default:
		// 默认按字符串比较
		return s.compareStrings(a, b)
	}
}

// compareNumbers 比较数字
func (s *MvTableSchemaService) compareNumbers(a, b any) int {
	aFloat, aOk := s.toFloat64(a)
	bFloat, bOk := s.toFloat64(b)

	if !aOk && !bOk {
		return 0
	}
	if !aOk {
		return -1
	}
	if !bOk {
		return 1
	}

	if aFloat < bFloat {
		return -1
	}
	if aFloat > bFloat {
		return 1
	}
	return 0
}

// compareDates 比较日期
func (s *MvTableSchemaService) compareDates(a, b any) int {
	aTime, aOk := s.toTime(a)
	bTime, bOk := s.toTime(b)

	if !aOk && !bOk {
		return 0
	}
	if !aOk {
		return -1
	}
	if !bOk {
		return 1
	}

	if aTime.Before(bTime) {
		return -1
	}
	if aTime.After(bTime) {
		return 1
	}
	return 0
}

// compareStrings 比较字符串
func (s *MvTableSchemaService) compareStrings(a, b any) int {
	aStr := s.toString(a)
	bStr := s.toString(b)
	return strings.Compare(aStr, bStr)
}

// toString 转换为字符串
func (s *MvTableSchemaService) toString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		// 对于复杂类型，转换为JSON字符串
		if data, err := json.Marshal(val); err == nil {
			return string(data)
		}
		return fmt.Sprintf("%v", val)
	}
}

// toFloat64 转换为float64
func (s *MvTableSchemaService) toFloat64(v any) (float64, bool) {
	if v == nil {
		return 0, false
	}

	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}

	return 0, false
}

// toTime 转换为时间
func (s *MvTableSchemaService) toTime(v any) (time.Time, bool) {
	if v == nil {
		return time.Time{}, false
	}

	switch val := v.(type) {
	case time.Time:
		return val, true
	case string:
		// 尝试多种时间格式
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02",
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t, true
			}
		}
	}

	return time.Time{}, false
}

// resolvePersonField 解析人员字段，返回用户信息数组
func (s *MvTableSchemaService) resolvePersonField(rowData model.RecordItem, fieldID string) ([]map[string]any, error) {
	fieldValue, ok := rowData[fieldID]
	if !ok || fieldValue == nil {
		return []map[string]any{}, nil
	}

	// 尝试将值转换为用户ID数组
	var userIDs []string
	switch v := fieldValue.(type) {
	case []any:
		for _, item := range v {
			if id, ok := item.(string); ok {
				userIDs = append(userIDs, id)
			}
		}
	case []string:
		userIDs = v
	case string:
		userIDs = []string{v}
	default:
		return []map[string]any{}, nil
	}

	if len(userIDs) == 0 {
		return []map[string]any{}, nil
	}

	// 查询用户信息
	users, _, err := db.List[userModel.User](db.GetDB(), 0, 0, map[string]any{"id": []any{"IN", userIDs}}, nil)
	if err != nil {
		return nil, err
	}

	// 构建返回结果
	result := make([]map[string]any, 0, len(users))
	for _, user := range users {
		userInfo := map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		}
		result = append(result, userInfo)
	}

	return result, nil
}

// resolveCreatorField 解析创建者字段，返回用户信息
func (s *MvTableSchemaService) resolveCreatorField(rowData model.RecordItem, fieldID string) (map[string]any, error) {
	fieldValue, ok := rowData[fieldID]
	if !ok || fieldValue == nil {
		return nil, nil
	}

	// 尝试将值转换为用户ID字符串
	var userID string
	switch v := fieldValue.(type) {
	case string:
		userID = v
	default:
		return nil, nil
	}

	return s.getUserInfo(userID)
}

// setRecordAuditInfo 设置记录的审计信息（创建时间、创建者、更新时间、更新者）
func (s *MvTableSchemaService) setRecordAuditInfo(rowData *model.RecordItem, record *model3.MvRecord) {
	// 设置创建时间
	(*rowData)["createdAt"] = record.CreatedAt

	// 设置创建者信息
	if record.CreatedBy != "" {
		creatorInfo, err := s.getUserInfo(record.CreatedBy)
		if err != nil {
			log.Error("get creator info error", zap.Error(err), zap.String("userId", record.CreatedBy))
			(*rowData)["creator"] = nil
		} else {
			(*rowData)["creator"] = creatorInfo
		}
	} else {
		(*rowData)["creator"] = nil
	}

	// 设置更新时间
	(*rowData)["updatedAt"] = record.UpdatedAt

	// 设置更新者信息，从rowData中的updatedBy字段获取
	if updatedBy, ok := (*rowData)["updatedBy"]; ok && updatedBy != nil {
		if updatedByStr, ok := updatedBy.(string); ok && updatedByStr != "" {
			updaterInfo, err := s.getUserInfo(updatedByStr)
			if err != nil {
				log.Error("get updater info error", zap.Error(err), zap.String("userId", updatedByStr))
				(*rowData)["updater"] = nil
			} else {
				(*rowData)["updater"] = updaterInfo
			}
		} else {
			(*rowData)["updater"] = nil
		}
	} else {
		(*rowData)["updater"] = nil
	}
}

// getUserInfo 获取用户信息，返回包含id、email、nickname、avatar的map
func (s *MvTableSchemaService) getUserInfo(userID string) (map[string]any, error) {
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
	result := map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
	}

	return result, nil
}

// stripHtmlTags 去除HTML标签，只保留纯文本内容
func stripHtmlTags(html string) string {
	// 去除HTML标签的正则表达式
	re := regexp.MustCompile(`<[^>]*>`)
	// 去除标签
	text := re.ReplaceAllString(html, "")
	// 去除多余的空白字符
	text = strings.TrimSpace(text)
	// 去除连续的空白字符
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	return text
}
