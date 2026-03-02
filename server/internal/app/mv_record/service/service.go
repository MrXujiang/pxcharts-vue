package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	model2 "mvtable/internal/app/mv_field/model"
	"mvtable/internal/app/mv_record/model"
	"mvtable/internal/app/mv_record/repo"
	richTextRepo "mvtable/internal/app/mv_rich_text_content/repo"
	richTextService "mvtable/internal/app/mv_rich_text_content/service"
	schemaModel "mvtable/internal/app/mv_table_schema/model"
	schemaService "mvtable/internal/app/mv_table_schema/service"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/lexorank"
	"mvtable/pkg/log"
	"strings"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvRecordService struct {
	richTextContentService *richTextService.MvRichTextContentService
	tableSchemaService     *schemaService.MvTableSchemaService
}

func NewMvRecordService() *MvRecordService {
	return &MvRecordService{
		richTextContentService: richTextService.NewMvRichTextContentService(),
		tableSchemaService:     schemaService.NewMvTableSchemaService(),
	}
}

func (s *MvRecordService) InsertMvRecord(userId string, req *model.InsertMvRecordReq) (any, error) {
	var (
		err        error
		prevRecord *model.MvRecord
		fields     []*model2.MvField
		record     = &model.MvRecord{
			TableSchemaID: req.TableSchemaID,
			CreatedBy:     userId,
		}
	)

	// 获取上一行record
	prevRecord, err = db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.PrevID})

	if err != nil {
		log.Error("get prev record error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 设置排序索引
	if prevRecord != nil {
		record.OrderIndex = lexorank.Between(prevRecord.OrderIndex, lexorank.MaxString)
	} else {
		// 查询当前的第一条记录
		firstRecords, _, err := db.List[model.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, []string{"order_index ASC"})
		if err != nil {
			log.Error("get first record error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		if len(firstRecords) > 0 {
			record.OrderIndex = lexorank.Between(lexorank.MinString, firstRecords[0].OrderIndex)
		} else {
			record.OrderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
		}
	}

	// 查询所有字段，并设置默认值
	fields, _, err = db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, nil)

	if err != nil {
		log.Error("get fields error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	// 为记录设置字段值（支持传入rowData）
	rowData, err := s.setFieldValues(fields, req.TableSchemaID, userId, req.RowData)
	if err != nil {
		log.Error("set field values error", zap.Error(err))
		return nil, errorx.InternalServerError("设置字段值失败")
	}
	record.RowData = rowData

	log.Info("fields", zap.Any("fields", fields))

	if err = db.Create(db.GetDB(), record); err != nil {
		log.Error("create record error", zap.Error(err))
		return nil, errorx.InternalServerError("创建失败")
	}

	// 获取完整的记录数据
	fullRecordData, err := s.tableSchemaService.GetRowData(record, fields)
	if err != nil {
		log.Error("get full record data error", zap.Error(err), zap.String("recordId", record.ID))
		return nil, errorx.InternalServerError("获取记录数据失败")
	}

	return fullRecordData, nil
}

// GetRecord 获取单条完整记录数据
func (s *MvRecordService) GetRecord(req *model.GetRecordReq) (any, error) {
	// 校验表格是否存在
	tableSchema, err := db.Get[schemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": req.TableSchemaID})
	if err != nil {
		log.Error("get table schema error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return nil, errorx.InternalServerError("获取表格信息失败")
	}
	if tableSchema == nil {
		return nil, errorx.New(errorx.ErrNotFound, "表格不存在")
	}

	// 获取记录
	record, err := db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.RecordID})
	if err != nil {
		log.Error("get record error", zap.Error(err), zap.String("recordId", req.RecordID))
		return nil, errorx.InternalServerError("获取记录失败")
	}
	if record == nil {
		return nil, errorx.New(errorx.ErrNotFound, "记录不存在")
	}

	// 验证记录是否属于指定的表格
	if record.TableSchemaID != req.TableSchemaID {
		return nil, errorx.New(errorx.ErrNotFound, "记录不属于指定的表格")
	}

	// 获取字段列表
	fields, _, err := db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, nil)
	if err != nil {
		log.Error("get fields error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return nil, errorx.InternalServerError("获取字段列表失败")
	}

	// 获取完整的记录数据
	fullRecordData, err := s.tableSchemaService.GetRowData(record, fields)
	if err != nil {
		log.Error("get full record data error", zap.Error(err), zap.String("recordId", record.ID))
		return nil, errorx.InternalServerError("获取记录数据失败")
	}

	return fullRecordData, nil
}

// GetRecords 获取表格所有完整记录数据
func (s *MvRecordService) GetRecords(req *model.GetRecordsReq) (*model.GetRecordsRes, error) {
	// 校验表格是否存在
	tableSchema, err := db.Get[schemaModel.MvTableSchema](db.GetDB(), map[string]any{"id": req.TableSchemaID})
	if err != nil {
		log.Error("get table schema error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return nil, errorx.InternalServerError("获取表格信息失败")
	}
	if tableSchema == nil {
		return nil, errorx.New(errorx.ErrNotFound, "表格不存在")
	}

	// 获取字段列表
	fields, _, err := db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get fields error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return nil, errorx.InternalServerError("获取字段列表失败")
	}

	// 处理字段信息，参考GetMvFieldList的逻辑
	fieldItems := make([]schemaModel.GeneralFieldConfig, len(fields))
	for i, field := range fields {
		// 获取字段配置
		fieldConfig, err := schemaService.GetFieldConfig(field.Config, field.Type, field.ID, field.Title)
		if err != nil {
			log.Error("get field config error", zap.Error(err), zap.Any("field", field))
			return nil, errorx.InternalServerError("获取失败")
		}

		fieldItems[i] = fieldConfig

		// 设置第一个字段为固定列
		if i == 0 {
			fieldItems[i].GetBase().Fixed = true
		}
	}

	// 获取所有记录，按order_index排序
	records, _, err := db.List[model.MvRecord](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get records error", zap.Error(err), zap.String("tableSchemaId", req.TableSchemaID))
		return nil, errorx.InternalServerError("获取记录列表失败")
	}

	// 处理每条记录的完整数据
	recordItems := make([]any, len(records))
	for i, record := range records {
		// 获取完整的记录数据
		fullRecordData, err := s.tableSchemaService.GetRowData(record, fields)
		if err != nil {
			log.Error("get full record data error", zap.Error(err), zap.String("recordId", record.ID))
			return nil, errorx.InternalServerError("获取记录数据失败")
		}
		recordItems[i] = fullRecordData
	}

	return &model.GetRecordsRes{
		Fields:  fieldItems,
		Records: recordItems,
	}, nil
}

// GetRecordOptions 获取表格记录选项列表
func (s *MvRecordService) GetRecordOptions(req *model.GetRecordOptionsReq) (*model.GetRecordOptionsRes, error) {
	// 调用 GetRecords 获取完整数据
	recordsRes, err := s.GetRecords(&model.GetRecordsReq{TableSchemaID: req.TableSchemaID})
	if err != nil {
		return nil, err
	}

	// 如果没有字段或记录，返回空列表
	if len(recordsRes.Fields) == 0 || len(recordsRes.Records) == 0 {
		return &model.GetRecordOptionsRes{Options: []model.RecordOption{}}, nil
	}

	// 获取第一个字段的 ID
	firstFieldID := recordsRes.Fields[0].GetBase().ID

	// 映射记录数据
	options := make([]model.RecordOption, 0, len(recordsRes.Records))
	for _, record := range recordsRes.Records {
		// 断言为 RecordItem 类型
		recordItem, ok := record.(schemaModel.RecordItem)
		if !ok {
			log.Error("invalid record format", zap.String("type", fmt.Sprintf("%T", record)), zap.Any("record", record))
			continue
		}

		rowId, ok := recordItem["rowId"].(string)
		if !ok {
			log.Error("rowId not found in record", zap.Any("record", record))
			continue
		}

		value := recordItem[firstFieldID]

		if value == "" || value == nil {
			continue
		}

		options = append(options, model.RecordOption{
			ID:    rowId,
			Value: value,
		})
	}

	return &model.GetRecordOptionsRes{
		Options: options,
	}, nil
}

func (s *MvRecordService) SaveRowSort(userId string, req *model.SaveRowSortReq) error {
	var (
		err           error
		prevRecord    *model.MvRecord
		nextRecord    *model.MvRecord
		currentRecord *model.MvRecord
	)

	if req.PrevRecordID != "" {
		prevRecord, err = db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.PrevRecordID})
		if err != nil {
			log.Error("get prev record error", zap.Error(err))
			return errorx.InternalServerError("获取失败")
		}

		if prevRecord == nil {
			return errorx.New(errorx.ErrNotFound, "上一行记录不存在")
		}
	}

	if req.NextRecordID != "" {
		nextRecord, err = db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.NextRecordID})
		if err != nil {
			log.Error("get next record error", zap.Error(err))
			return errorx.InternalServerError("获取失败")
		}

		if nextRecord == nil {
			return errorx.New(errorx.ErrNotFound, "下一行记录不存在")
		}
	}

	currentRecord, err = db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.CurrentRecordID})
	if err != nil {
		log.Error("get current record error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}

	if currentRecord == nil {
		return errorx.New(errorx.ErrNotFound, "当前行记录不存在")
	}

	// 计算当前记录的排序索引
	if req.PrevRecordID == "" && req.NextRecordID != "" {
		currentRecord.OrderIndex = lexorank.Between(lexorank.MinString, nextRecord.OrderIndex)
	}
	if req.PrevRecordID != "" && req.NextRecordID == "" {
		currentRecord.OrderIndex = lexorank.Between(prevRecord.OrderIndex, lexorank.MaxString)
	}
	if req.PrevRecordID != "" && req.NextRecordID != "" {
		currentRecord.OrderIndex = lexorank.Between(prevRecord.OrderIndex, nextRecord.OrderIndex)
	}

	if err := db.Update(db.GetDB(), &currentRecord, map[string]any{"id": req.CurrentRecordID}, "order_index"); err != nil {
		log.Error("update current record error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

func (s *MvRecordService) UpdateCell(userId string, req *model.UpdateCellReq) error {
	var (
		field *model2.MvField
		err   error
	)

	field, err = db.Get[model2.MvField](db.GetDB(), map[string]any{"id": req.FieldID})
	if err != nil {
		log.Error("get field error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}

	if field == nil {
		return errorx.New(errorx.ErrNotFound, "字段不存在")
	}

	record, err := db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": req.RecordID})
	if err != nil {
		log.Error("get record error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if record == nil {
		return errorx.New(errorx.ErrNotFound, "记录不存在")
	}

	// 使用事务确保字段值更新和UpdatedBy设置的原子性
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 如果是富文本类型，需要特殊处理
		switch field.Type {
		case constants.MvFieldTypeRichText:
			if err := s.handleRichTextField(tx, record, field, req.Value); err != nil {
				return err
			}
		case constants.MvFieldTypeDoubleAssociation:
			if err := s.handleDoubleAssociation(tx, field, record, req.Value); err != nil {
				return err
			}

		default:
			if err := repo.SetFieldValue(tx, record.ID, field.ID, req.Value); err != nil {
				log.Error("set field value error", zap.Error(err))
				return errorx.InternalServerError("更新失败")
			}
		}

		// 设置更新人ID
		if err := db.Update(tx, &model.MvRecord{UpdatedBy: userId}, map[string]any{"id": record.ID}, "updated_by"); err != nil {
			log.Error("set updated by error", zap.Error(err))
			return errorx.InternalServerError("设置更新人失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *MvRecordService) UpdateRow(userId string, req model.UpdateRowReq) error {
	var (
		err error
	)

	// 从请求中提取 rowId
	rowIDValue, exists := req["rowId"]
	if !exists {
		return errorx.BadRequest("缺少 rowId 字段")
	}
	rowID, ok := rowIDValue.(string)
	if !ok {
		return errorx.BadRequest("rowId 必须为字符串类型")
	}

	// 获取记录
	record, err := db.Get[model.MvRecord](db.GetDB(), map[string]any{"id": rowID})
	if err != nil {
		log.Error("get record error", zap.Error(err))
		return errorx.InternalServerError("获取失败")
	}
	if record == nil {
		return errorx.New(errorx.ErrNotFound, "记录不存在")
	}

	// 使用事务确保所有字段更新和UpdatedBy设置的原子性
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 遍历所有要更新的字段（除了 rowId）
		for fieldID, value := range req {
			if fieldID == "rowId" {
				continue // 跳过 rowId 字段
			}
			// 获取字段信息
			field, err := db.Get[model2.MvField](tx, map[string]any{"id": fieldID})
			if err != nil {
				log.Error("get field error", zap.Error(err), zap.String("fieldId", fieldID))
				return errorx.InternalServerError("获取字段失败")
			}

			if field == nil {
				log.Error("field not found", zap.String("fieldId", fieldID))
				return errorx.New(errorx.ErrNotFound, "字段不存在")
			}

			// 根据字段类型进行不同的处理
			switch field.Type {
			case constants.MvFieldTypeRichText:
				if err := s.handleRichTextField(tx, record, field, value); err != nil {
					return err
				}
			case constants.MvFieldTypeDoubleAssociation:
				if err := s.handleDoubleAssociation(tx, field, record, value); err != nil {
					return err
				}
			default:
				if err := repo.SetFieldValue(tx, record.ID, field.ID, value); err != nil {
					log.Error("set field value error", zap.Error(err), zap.String("fieldId", fieldID))
					return errorx.InternalServerError("更新失败")
				}
			}
		}

		// 设置更新人ID
		if err := db.Update(tx, &model.MvRecord{UpdatedBy: userId}, map[string]any{"id": record.ID}, "updated_by"); err != nil {
			log.Error("set updated by error", zap.Error(err))
			return errorx.InternalServerError("设置更新人失败")
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *MvRecordService) DeleteRecords(userId string, req *model.DeleteRecordsReq) error {
	if len(req.RecordIDs) == 0 {
		return errorx.New(errorx.ErrInvalidParam, "记录ID列表不能为空")
	}

	// 使用事务确保数据一致性
	err := db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的富文本内容
		if err := s.richTextContentService.DeleteRichTextContentByRecordIDs(req.RecordIDs); err != nil {
			log.Error("delete rich text content error", zap.Error(err))
			return err
		}

		// 删除记录
		if err := tx.Where("id IN ?", req.RecordIDs).Delete(&model.MvRecord{}).Error; err != nil {
			log.Error("delete records error", zap.Error(err), zap.Strings("recordIds", req.RecordIDs))
			return err
		}

		return nil
	})

	if err != nil {
		return errorx.InternalServerError("删除失败")
	}

	return nil
}

// ImportRecords 从 CSV 或 Excel 文件导入数据
func (s *MvRecordService) ImportRecords(userId string, tableSchemaID string, file *multipart.FileHeader) (*model.ImportRecordsRes, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		log.Error("open file error", zap.Error(err))
		return nil, errorx.BadRequest("无法打开文件")
	}
	defer src.Close()

	// 获取文件扩展名
	fileName := file.Filename
	ext := strings.ToLower(fileName[strings.LastIndex(fileName, ".")+1:])

	// 获取表格的所有字段
	fields, _, err := db.List[model2.MvField](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchemaID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("get fields error", zap.Error(err))
		return nil, errorx.InternalServerError("获取字段失败")
	}

	// 创建字段标题到字段ID的映射
	fieldMap := make(map[string]*model2.MvField)
	for _, field := range fields {
		fieldMap[field.Title] = field
	}

	// 解析文件
	var headers []string
	var rows [][]string
	var parseErr error

	switch ext {
	case "csv":
		headers, rows, parseErr = parseCSV(src)
	case "xlsx", "xls":
		headers, rows, parseErr = parseExcel(src)
	default:
		return nil, errorx.BadRequest("不支持的文件格式，仅支持 CSV 和 Excel 文件")
	}

	if parseErr != nil {
		log.Error("parse file error", zap.Error(parseErr))
		return nil, errorx.BadRequest(fmt.Sprintf("解析文件失败: %v", parseErr))
	}

	if len(headers) == 0 {
		return nil, errorx.BadRequest("文件为空或格式不正确")
	}

	// 验证标题行，检查哪些字段可以匹配
	matchedFields := make(map[int]*model2.MvField) // 列索引 -> 字段
	for i, header := range headers {
		if field, exists := fieldMap[strings.TrimSpace(header)]; exists {
			matchedFields[i] = field
		}
	}

	if len(matchedFields) == 0 {
		return nil, errorx.BadRequest("文件中没有匹配的字段，请检查标题行")
	}

	// 获取最后一条记录的排序索引（按 order_index 降序排列的第一条）
	lastRecords, _, err := db.List[model.MvRecord](db.GetDB(), 1, 1, map[string]any{"table_schema_id": tableSchemaID}, []string{"order_index DESC"})
	if err != nil {
		log.Error("get last record error", zap.Error(err))
		return nil, errorx.InternalServerError("获取排序索引失败")
	}

	var startOrderIndex string
	if len(lastRecords) > 0 && lastRecords[0] != nil {
		startOrderIndex = lastRecords[0].OrderIndex
	} else {
		startOrderIndex = lexorank.MinString
	}

	// 批量导入记录
	res := &model.ImportRecordsRes{
		TotalCount: len(rows),
		Errors:     []model.ImportError{},
	}

	// 使用事务批量创建
	err = db.Transaction(func(tx *gorm.DB) error {
		for rowIndex, row := range rows {
			// 构建 row_data
			rowData := make(map[string]any)
			hasData := false
			richTextFields := make(map[string]string) // fieldID -> content

			for colIndex, value := range row {
				if field, exists := matchedFields[colIndex]; exists {
					trimmedValue := strings.TrimSpace(value)
					if trimmedValue != "" {
						// 如果是富文本类型，保存到富文本内容表
						if field.Type == constants.MvFieldTypeRichText {
							richTextFields[field.ID] = trimmedValue
							// 在 row_data 中保存标记
							rowData[field.ID] = map[string]any{
								"type":    constants.MvFieldTypeRichText,
								"fieldId": field.ID,
							}
						} else {
							rowData[field.ID] = trimmedValue
						}
						hasData = true
					}
				}
			}

			// 如果该行没有任何数据，跳过
			if !hasData {
				res.FailCount++
				res.Errors = append(res.Errors, model.ImportError{
					Row:   rowIndex + 2, // +2 因为标题行是第1行，数据从第2行开始
					Error: "该行没有有效数据",
				})
				continue
			}

			// 创建记录
			record := &model.MvRecord{
				TableSchemaID: tableSchemaID,
				CreatedBy:     userId,
			}

			// 设置排序索引
			if startOrderIndex == lexorank.MinString {
				record.OrderIndex = lexorank.Between(lexorank.MinString, lexorank.MaxString)
				startOrderIndex = record.OrderIndex
			} else {
				record.OrderIndex = lexorank.Between(startOrderIndex, lexorank.MaxString)
				startOrderIndex = record.OrderIndex
			}

			// 将 rowData 转换为 JSON
			rowDataJSON, err := json.Marshal(rowData)
			if err != nil {
				res.FailCount++
				res.Errors = append(res.Errors, model.ImportError{
					Row:   rowIndex + 2,
					Error: fmt.Sprintf("数据格式错误: %v", err),
				})
				continue
			}
			record.RowData = rowDataJSON

			// 创建记录
			if err := db.Create(tx, record); err != nil {
				log.Error("create record error", zap.Error(err), zap.Int("row", rowIndex+2))
				res.FailCount++
				res.Errors = append(res.Errors, model.ImportError{
					Row:   rowIndex + 2,
					Error: fmt.Sprintf("创建记录失败: %v", err),
				})
				continue
			}

			// 保存富文本内容
			for fieldID, content := range richTextFields {
				if err := s.richTextContentService.SaveRichTextContent(record.ID, fieldID, content); err != nil {
					log.Error("save rich text content error", zap.Error(err), zap.Int("row", rowIndex+2))
					// 记录错误但不中断导入流程
					res.Errors = append(res.Errors, model.ImportError{
						Row:   rowIndex + 2,
						Error: fmt.Sprintf("保存富文本内容失败: %v", err),
					})
				}
			}

			res.SuccessCount++
		}

		return nil
	})

	if err != nil {
		log.Error("import records transaction error", zap.Error(err))
		return nil, errorx.InternalServerError("导入失败")
	}

	return res, nil
}

// parseCSV 解析 CSV 文件
func parseCSV(reader io.Reader) ([]string, [][]string, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	// 读取标题行
	headers, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil, fmt.Errorf("文件为空")
		}
		return nil, nil, fmt.Errorf("读取CSV文件失败: %w", err)
	}

	// 读取所有数据行
	var rows [][]string
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		rows = append(rows, row)
	}

	return headers, rows, nil
}

// parseExcel 解析 Excel 文件
func parseExcel(reader io.Reader) ([]string, [][]string, error) {
	// 读取文件内容到内存
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 打开 Excel 文件
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, fmt.Errorf("打开Excel文件失败: %w", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, nil, fmt.Errorf("Excel文件没有工作表")
	}

	// 获取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("Excel文件为空")
	}

	// 第一行为标题
	headers := rows[0]

	// 剩余行为数据
	dataRows := rows[1:]

	return headers, dataRows, nil
}

// handleRichTextField 处理富文本字段的更新逻辑
func (s *MvRecordService) handleRichTextField(tx *gorm.DB, record *model.MvRecord, field *model2.MvField, value any) error {
	// 将值转换为字符串
	var content string
	if value != nil {
		if str, ok := value.(string); ok {
			content = str
		} else {
			// 如果不是字符串，尝试转换为JSON字符串
			contentBytes, err := json.Marshal(value)
			if err != nil {
				log.Error("marshal rich text value error", zap.Error(err))
				return errorx.BadRequest("富文本内容格式错误")
			}
			content = string(contentBytes)
		}
	}

	// 如果内容为空，删除富文本内容
	if content == "" {
		if err := richTextRepo.DeleteRichTextContent(tx, record.ID, field.ID); err != nil {
			log.Error("delete rich text content error", zap.Error(err))
			// 即使删除失败也继续，因为可能是内容不存在
		}
		// 在 row_data 中设置为 null
		if err := repo.SetFieldValue(tx, record.ID, field.ID, nil); err != nil {
			log.Error("set field value error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
	} else {
		// 保存富文本内容到专门的表
		if err := richTextRepo.SaveRichTextContent(tx, record.ID, field.ID, content); err != nil {
			log.Error("save rich text content error", zap.Error(err))
			return errorx.InternalServerError("保存富文本内容失败")
		}

		// 在 row_data 中保存一个标记，表示这是富文本字段
		// 可以保存字段ID作为引用，或者保存一个特殊标记
		marker := map[string]any{
			"type":    constants.MvFieldTypeRichText,
			"fieldId": field.ID,
		}
		if err := repo.SetFieldValue(tx, record.ID, field.ID, marker); err != nil {
			log.Error("set field value error", zap.Error(err))
			return errorx.InternalServerError("更新失败")
		}
	}

	return nil
}
