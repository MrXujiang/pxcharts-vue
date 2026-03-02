package service

import (
	"encoding/json"
	"fmt"
	"mvtable/internal/app/mv_field/model"
	recordModel "mvtable/internal/app/mv_record/model"
	recordRepo "mvtable/internal/app/mv_record/repo"
	schemaModel "mvtable/internal/app/mv_table_schema/model"
	schemaService "mvtable/internal/app/mv_table_schema/service"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MvFieldService struct{}

func NewMvFieldService() *MvFieldService {
	return &MvFieldService{}
}

func (s *MvFieldService) GetMvFieldList(req *model.GetMvFieldListReq) (*model.GetMvFieldListRes, error) {
	// 按照 GetTableData 的逻辑获取字段列表
	fields, _, err := db.List[model.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get mv field list error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	fieldsItems := make([]schemaModel.GeneralFieldConfig, len(fields))
	for i, v := range fields {
		var fieldConfig schemaModel.GeneralFieldConfig
		fieldConfig, err = schemaService.GetFieldConfig(v.Config, v.Type, v.ID, v.Title)
		if err != nil {
			log.Error("get field config error", zap.Error(err), zap.Any("field", v))
			return nil, errorx.InternalServerError("获取失败")
		}

		fieldsItems[i] = fieldConfig

		// 设置第一个字段为固定列
		if i == 0 {
			fieldsItems[i].GetBase().Fixed = true
		}
	}

	return &model.GetMvFieldListRes{List: fieldsItems}, nil
}

func (s *MvFieldService) CreateMvField(req *model.CreateMvFieldReq) (string, error) {
	config, err := json.Marshal(req.Config)
	if err != nil {
		log.Error("marshal config error", zap.Error(err))
		return "", errorx.InternalServerError("创建失败")
	}

	field := model.MvField{
		TableSchemaID: req.TableSchemaID,
		Title:         req.Title,
		Type:          req.Type,
		Config:        config,
	}

	fields, _, err := db.List[model.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, []string{"order_index ASC"})
	if len(fields) == 0 {
		field.OrderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
	} else {
		field.OrderIndex = lexorank.Between(fields[len(fields)-1].OrderIndex, lexorank.MaxString)
	}
	if err = db.Create(db.GetDB(), &field); err != nil {
		log.Error("create mv field error", zap.Error(err))
		return "", errorx.InternalServerError("创建失败")
	}

	// 根据字段配置设置默认值
	fieldConfig, err := schemaService.GetFieldConfig(config, req.Type, field.ID, req.Title)
	if err != nil {
		log.Error("get field config error", zap.Error(err))
		// 不返回错误，因为字段已经创建成功，只是设置默认值失败
	} else {
		// 根据字段配置类型设置默认值
		switch cfg := fieldConfig.(type) {
		case *schemaModel.FieldCreateTimeConfig:
			// 创建时间字段：为所有现有记录设置记录的创建时间
			if err := s.setDefaultTimeForExistingRecords(req.TableSchemaID, field.ID, true); err != nil {
				log.Error("set create time for existing records error", zap.Error(err))
			}
		case *schemaModel.FieldAutoNumberConfig:
			// 自动编号字段：为所有现有记录按orderIndex顺序设置从1开始的编号
			if err := s.setAutoNumberForExistingRecords(req.TableSchemaID, field.ID); err != nil {
				log.Error("set auto number for existing records error", zap.Error(err))
			}
		case *schemaModel.FieldDateConfig:
			// 日期字段：如果设置为默认创建时间，则设置记录创建时间，否则设置配置的默认值
			if cfg.Settings.IsDefaultCreateTime {
				if err := s.setDefaultTimeForExistingRecords(req.TableSchemaID, field.ID, false); err != nil {
					log.Error("set default time for existing records error", zap.Error(err))
				}
			} else if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldTextConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldRadioConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldNumberConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldSelectMultipleConfig:
			if len(cfg.DefaultValue) > 0 {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldPersonConfig:
			if len(cfg.DefaultValue) > 0 {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldCurrencyConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldImageConfig:
			if len(cfg.DefaultValue) > 0 {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldProgressConfig:
			if cfg.DefaultValue != nil {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldLinkConfig:
			if cfg.DefaultValue.LinkUrl != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldCheckboxConfig:
			if cfg.DefaultValue >= 0 {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldRateConfig:
			if cfg.DefaultValue != nil {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldEmailConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldIdCardConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldRichTextConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldFindRefConfig:
			if cfg.DefaultValue != "" {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, cfg.DefaultValue); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		case *schemaModel.FieldCreatorConfig:
			// 创建人字段：为所有现有记录设置创建人ID为记录的CreatedBy字段值
			if err := s.setCreatorForExistingRecords(req.TableSchemaID, field.ID); err != nil {
				log.Error("set creator for existing records error", zap.Error(err))
			}
		case *schemaModel.FieldSingleAssociationConfig:
			// 单向关联不设置默认值
		case *schemaModel.FieldDoubleAssociationConfig:
			// 双向关联不设置默认值
		case *schemaModel.FieldAttachmentConfig:
			// 附件字段不设置默认值
		case *schemaModel.FieldFormulaConfig:
			// 公式字段不设置默认值
		default:
			// 对于其他字段类型，如果有配置的默认值则设置
			if req.Config != nil {
				if err := s.setDefaultValueForExistingRecords(req.TableSchemaID, field.ID, req.Config); err != nil {
					log.Error("set default value for existing records error", zap.Error(err))
				}
			}
		}
	}

	return field.ID, nil
}

func (s *MvFieldService) UpdateMvField(req *model.UpdateMvFieldReq) error {
	var (
		field = &model.MvField{
			Title: req.Title,
			Type:  req.Type,
		}
		err error
	)

	if req.Config != nil {
		config, err := json.Marshal(req.Config)
		if err != nil {
			log.Error("marshal config error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
		field.Config = config
	}

	if err = db.Update(db.GetDB(), &field, map[string]any{"id": req.ID}, "title", "type", "config"); err != nil {
		log.Error("update mv field error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	// 如果字段类型被更新为创建时间类型，需要为所有现有记录设置该字段的值为记录创建时间
	if req.Type == "createTime" {
		// 获取字段信息以获取 table_schema_id
		updatedField, err := db.Get[model.MvField](db.GetDB(), map[string]any{"id": req.ID})
		if err != nil {
			log.Error("get updated field error", zap.Error(err))
			// 不返回错误，因为字段已经更新成功，只是设置默认值失败
		} else if updatedField != nil {
			if err := s.setDefaultTimeForExistingRecords(updatedField.TableSchemaID, req.ID, true); err != nil {
				log.Error("set create time for existing records error", zap.Error(err))
				// 不返回错误，因为字段已经更新成功，只是设置默认值失败
			}
		}
	} else if req.Type == constants.MvFieldTypePerson {
		// 如果字段类型被更新为人员类型且有自定义默认值，需要为所有现有记录设置默认值
		// 获取字段信息以获取 table_schema_id
		updatedField, err := db.Get[model.MvField](db.GetDB(), map[string]any{"id": req.ID})
		if err != nil {
			log.Error("get updated field error", zap.Error(err))
			// 不返回错误，因为字段已经更新成功，只是设置默认值失败
		} else if updatedField != nil {
			if err := s.setDefaultValueForExistingRecords(updatedField.TableSchemaID, req.ID, req.Config); err != nil {
				log.Error("set default value for existing records error", zap.Error(err))
				// 不返回错误，因为字段已经更新成功，只是设置默认值失败
			}
		}
	}

	return nil
}

func (s *MvFieldService) DeleteMvField(req *model.DeleteMvFieldReq) error {
	// 首先获取要删除的字段信息
	field, err := db.Get[model.MvField](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get field error", zap.Error(err))
		return errorx.InternalServerError("获取字段信息失败")
	}
	if field == nil {
		return errorx.New(errorx.ErrNotFound, "字段不存在")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// 删除字段
		if err := db.Delete[model.MvField](tx, map[string]any{"id": req.ID}); err != nil {
			log.Error("delete mv field error", zap.Error(err))
			return errorx.InternalServerError("删除字段失败")
		}

		// 删除 mv_record 表中与该字段相关的数据
		if err := s.removeFieldDataFromRecords(tx, field); err != nil {
			log.Error("remove field data from records error", zap.Error(err))
			return errorx.InternalServerError("清理记录数据失败")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// SortMvFields 对表格字段进行排序
func (s *MvFieldService) SortMvFields(req *model.SortMvFieldReq) error {
	var (
		fields []*model.MvField
		err    error
	)

	if err = db.GetDB().Where("table_schema_id = ?", req.TableSchemaID).Find(&fields).Error; err != nil {
		log.Error("list mv fields error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return errorx.InternalServerError("获取失败")
	}
	if len(fields) == 0 {
		return errorx.New(errorx.ErrNotFound, "表格字段不存在")
	}
	if len(fields) != len(req.FieldIDs) {
		return errorx.New(errorx.ErrInvalidParam, "字段数量与排序列表不一致")
	}

	fieldMap := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		fieldMap[f.ID] = struct{}{}
	}

	seen := make(map[string]struct{}, len(req.FieldIDs))
	for _, id := range req.FieldIDs {
		if _, ok := fieldMap[id]; !ok {
			return errorx.New(errorx.ErrInvalidParam, "存在不属于该表格的字段")
		}
		if _, ok := seen[id]; ok {
			return errorx.New(errorx.ErrInvalidParam, "字段列表存在重复项")
		}
		seen[id] = struct{}{}
	}

	if err = db.Transaction(func(tx *gorm.DB) error {
		prev := lexorank.MinString
		for idx, fieldID := range req.FieldIDs {
			var orderIndex string
			if idx == 0 {
				orderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
			} else {
				orderIndex = lexorank.Between(prev, lexorank.MaxString)
			}

			if err := db.Update(tx, &model.MvField{OrderIndex: orderIndex}, map[string]any{"id": fieldID}, "order_index"); err != nil {
				log.Error("update mv field order error", zap.Error(err), zap.String("fieldId", fieldID))
				return errorx.InternalServerError("更新失败")
			}
			prev = orderIndex
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// removeFieldDataFromRecords 从所有相关记录中删除指定字段的数据
func (s *MvFieldService) removeFieldDataFromRecords(tx *gorm.DB, field *model.MvField) error {
	const batchSize = 100 // 每批处理100条记录

	offset := 0
	for {
		// 分批查询记录
		var records []*recordModel.MvRecord
		if err := tx.Where("table_schema_id = ?", field.TableSchemaID).
			Offset(offset).
			Limit(batchSize).
			Find(&records).Error; err != nil {
			return fmt.Errorf("query records failed: %w", err)
		}

		// 如果没有更多记录，退出循环
		if len(records) == 0 {
			break
		}

		// 处理当前批次的记录
		for _, record := range records {
			if err := s.removeFieldFromRecord(tx, record, field.ID); err != nil {
				log.Warn("remove field from record failed",
					zap.String("recordId", record.ID),
					zap.String("fieldId", field.ID),
					zap.Error(err))
				// 继续处理其他记录，不因为单个记录失败而中断整个操作
			}
		}

		offset += batchSize
	}

	return nil
}

// removeFieldFromRecord 从单个记录中删除指定字段
func (s *MvFieldService) removeFieldFromRecord(tx *gorm.DB, record *recordModel.MvRecord, fieldID string) error {
	// 解析 row_data JSON
	var rowData map[string]any
	if err := json.Unmarshal(record.RowData, &rowData); err != nil {
		return fmt.Errorf("unmarshal row data failed: %w", err)
	}

	// 删除指定字段
	delete(rowData, fieldID)

	// 重新序列化为 JSON
	updatedRowData, err := json.Marshal(rowData)
	if err != nil {
		return fmt.Errorf("marshal updated row data failed: %w", err)
	}

	// 更新记录
	if err := tx.Model(record).Update("row_data", datatypes.JSON(updatedRowData)).Error; err != nil {
		return fmt.Errorf("update record failed: %w", err)
	}

	return nil
}

// setDefaultTimeForExistingRecords 为所有现有记录设置创建时间字段的值
func (s *MvFieldService) setDefaultTimeForExistingRecords(tableSchemaID, fieldID string, isCreateTimeField bool) error {
	// 查询所有该表格的记录
	records, _, err := db.List[recordModel.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, nil)
	if err != nil {
		return fmt.Errorf("list records failed: %w", err)
	}

	// 为每个记录设置创建时间字段的值
	for _, record := range records {
		// 使用记录的 created_at 作为字段值
		if err := recordRepo.SetFieldValue(db.GetDB(), record.ID, fieldID, record.CreatedAt); err != nil {
			log.Error("set create time field value error", zap.Error(err), zap.String("recordId", record.ID), zap.String("fieldId", fieldID))
			// 继续处理其他记录，不因为单个记录失败而停止
		}
	}

	return nil
}

// setDefaultValueForExistingRecords 为所有现有记录设置字段的默认值
func (s *MvFieldService) setDefaultValueForExistingRecords(tableSchemaID, fieldID string, defaultValue any) error {
	// 检查默认值是否有效
	if defaultValue == nil {
		return nil
	}

	// 查询所有该表格的记录
	records, _, err := db.List[recordModel.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, nil)
	if err != nil {
		return fmt.Errorf("list records failed: %w", err)
	}

	// 为每个记录设置默认值
	for _, record := range records {
		if err := recordRepo.SetFieldValue(db.GetDB(), record.ID, fieldID, defaultValue); err != nil {
			log.Error("set default value field value error", zap.Error(err), zap.String("recordId", record.ID), zap.String("fieldId", fieldID))
			// 继续处理其他记录，不因为单个记录失败而停止
		}
	}

	return nil
}

// setAutoNumberForExistingRecords 为所有现有记录设置自动编号，从1开始按orderIndex排序
func (s *MvFieldService) setAutoNumberForExistingRecords(tableSchemaID, fieldID string) error {
	// 查询所有该表格的记录，按orderIndex升序排序
	records, _, err := db.List[recordModel.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		return fmt.Errorf("list records failed: %w", err)
	}

	// 为每个记录按顺序设置从1开始的编号
	for i, record := range records {
		autoNumber := i + 1 // 从1开始编号
		if err := recordRepo.SetFieldValue(db.GetDB(), record.ID, fieldID, autoNumber); err != nil {
			log.Error("set auto number field value error", zap.Error(err), zap.String("recordId", record.ID), zap.String("fieldId", fieldID), zap.Int("autoNumber", autoNumber))
			// 继续处理其他记录，不因为单个记录失败而停止
		}
	}

	return nil
}

// setCreatorForExistingRecords 为所有现有记录设置创建人字段的值为记录的CreatedBy
func (s *MvFieldService) setCreatorForExistingRecords(tableSchemaID, fieldID string) error {
	// 查询所有该表格的记录
	records, _, err := db.List[recordModel.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, nil)
	if err != nil {
		return fmt.Errorf("list records failed: %w", err)
	}

	// 为每条记录设置创建人字段的值
	for _, record := range records {
		// 解析记录数据
		var rowData map[string]any
		if err := json.Unmarshal(record.RowData, &rowData); err != nil {
			log.Warn("unmarshal record data error", zap.Error(err), zap.String("recordId", record.ID))
			continue
		}

		// 设置创建人字段的值为记录的CreatedBy
		rowData[fieldID] = record.CreatedBy

		// 重新序列化并更新记录
		updatedData, err := json.Marshal(rowData)
		if err != nil {
			log.Warn("marshal updated record data error", zap.Error(err), zap.String("recordId", record.ID))
			continue
		}

		if err := db.Update(db.GetDB(), &recordModel.MvRecord{RowData: datatypes.JSON(updatedData)}, map[string]any{"id": record.ID}); err != nil {
			log.Warn("update record error", zap.Error(err), zap.String("recordId", record.ID))
			continue
		}
	}

	return nil
}
