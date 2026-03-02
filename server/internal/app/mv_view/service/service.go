package service

import (
	"encoding/json"
	"fmt"
	model4 "mvtable/internal/app/mv_table_schema/model"
	"mvtable/internal/app/mv_view/model"
	model6 "mvtable/internal/app/mv_view_active/model"
	model7 "mvtable/internal/app/mv_view_board/model"
	model3 "mvtable/internal/app/mv_view_form/model"
	model2 "mvtable/internal/app/mv_view_table/model"
	"mvtable/internal/pkg/constants"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"
	"regexp"
	"slices"
	"strconv"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MvViewFormService struct{}

func NewMvViewFormService() *MvViewFormService {
	return &MvViewFormService{}
}

func (s *MvViewFormService) QueryMvView(userId string, req *model.QueryMvViewReq) (*model.QueryMvViewRes, error) {
	var (
		tableSchema *model4.MvTableSchema
		views       []*model.MvView
		items       []model.MvViewItem
		err         error
	)

	// 判断表格是否存在
	tableSchema, err = db.Get[model4.MvTableSchema](db.GetDB(), map[string]any{"id": req.TableSchemaID})
	if err != nil {
		log.Error("get table schema error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}
	if tableSchema == nil {
		return nil, errorx.New(errorx.ErrNotFound, "表格不存在")
	}

	// 获取表格所有视图
	views, _, err = db.List[model.MvView](db.GetDB(), 0, 0, map[string]any{"table_schema_id": tableSchema.ID}, []string{"order_index ASC"})
	if err != nil {
		log.Error("query mv view error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}

	for _, view := range views {
		items = append(items, model.MvViewItem{
			ID:          view.ID,
			Type:        view.Type,
			Name:        view.Name,
			IsDefault:   false,
			Description: view.Description,
		})
	}

	// 查询默认激活视图
	viewActives, _, err := db.List[model6.MvViewActive](db.GetDB(), 0, 0, map[string]any{"table_schema_id": req.TableSchemaID, "user_id": userId}, nil)
	if err != nil {
		log.Error("query mv view active error", zap.Error(err))
		return nil, errorx.InternalServerError("查询失败")
	}
	if len(viewActives) == 0 {
		items[0].IsDefault = true
	} else {
		for i, v := range items {
			if v.ID == viewActives[0].ViewID {
				items[i].IsDefault = true
				break
			}
		}
	}

	return &model.QueryMvViewRes{
		List: items,
	}, nil
}

func (s *MvViewFormService) CreateMvView(req *model.CreateMvViewReq) (string, error) {
	if !slices.Contains([]string{constants.ViewTypeTable, constants.ViewTypeForm, constants.ViewTypeBoard}, req.Type) {
		return "", errorx.New(errorx.ErrInvalidParam, "视图类型不合法")
	}

	// 确定视图名称前缀
	var namePrefix string
	switch req.Type {
	case constants.ViewTypeTable:
		namePrefix = "表格"
	case constants.ViewTypeForm:
		namePrefix = "表单"
	case constants.ViewTypeBoard:
		namePrefix = "看板"
	}

	var newDashboardID string

	err := db.Transaction(func(tx *gorm.DB) error {
		// 在事务内查询同一表格下同类型的所有视图，确保并发安全
		existingViews, _, err := db.List[model.MvView](tx, 0, 0, map[string]any{
			"table_schema_id": req.TableSchemaID,
			"type":            req.Type,
		}, nil)
		if err != nil {
			log.Error("query existing views error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}

		// 生成唯一名称
		newName := s.generateUniqueViewName(namePrefix, existingViews)

		// 再次检查名称是否唯一（防止并发冲突）
		existView, err := db.Get[model.MvView](tx, map[string]any{
			"table_schema_id": req.TableSchemaID,
			"type":            req.Type,
			"name":            newName,
		})
		if err != nil {
			log.Error("check view name exists error", zap.Error(err))
			return errorx.InternalServerError("创建失败")
		}
		// 如果名称已存在，重新生成（这种情况在并发时可能出现）
		if existView != nil {
			// 重新查询所有视图并生成新名称
			allViews, _, err := db.List[model.MvView](tx, 0, 0, map[string]any{
				"table_schema_id": req.TableSchemaID,
				"type":            req.Type,
			}, nil)
			if err != nil {
				log.Error("re-query existing views error", zap.Error(err))
				return errorx.InternalServerError("创建失败")
			}
			newName = s.generateUniqueViewName(namePrefix, allViews)
		}

		view := &model.MvView{
			TableSchemaID: req.TableSchemaID,
			Type:          req.Type,
			Name:          newName,
		}

		if err = db.Create(tx, view); err != nil {
			log.Error("create mv view error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}
		newDashboardID = view.ID

		if req.Type == constants.ViewTypeTable {
			if err = db.Create(tx, &model2.MvViewTable{ViewID: view.ID}); err != nil {
				log.Error("create mv view table error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}
		if req.Type == constants.ViewTypeForm {
			if err = db.Create(tx, &model3.MvViewForm{
				ViewID:            view.ID,
				Stats:             datatypes.JSON([]byte(`{}`)),
				FilterConfig:      datatypes.JSON([]byte(`[]`)),
				CycleRemindConfig: datatypes.JSON([]byte(`{}`)),
				Config:            datatypes.JSON([]byte(`[]`)),
			}); err != nil {
				log.Error("create mv view form error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}
		if req.Type == constants.ViewTypeBoard {
			if err = db.Create(tx, &model7.MvViewBoard{ViewID: view.ID}); err != nil {
				log.Error("create mv view board error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return newDashboardID, nil
}

// generateUniqueViewName 生成唯一的视图名称
func (s *MvViewFormService) generateUniqueViewName(prefix string, existingViews []*model.MvView) string {
	// 匹配模式：前缀 + 数字（如：表格1、表格2、表单1等）
	pattern := regexp.MustCompile(fmt.Sprintf(`^%s(\d+)$`, regexp.QuoteMeta(prefix)))

	// 收集已存在的序号
	usedNumbers := make(map[int]bool)
	maxNumber := 0

	for _, view := range existingViews {
		if view == nil {
			continue
		}
		matches := pattern.FindStringSubmatch(view.Name)
		if len(matches) == 2 {
			// 提取序号
			num, err := strconv.Atoi(matches[1])
			if err == nil {
				usedNumbers[num] = true
				if num > maxNumber {
					maxNumber = num
				}
			}
		}
	}

	// 从1开始查找第一个未使用的序号
	candidateNumber := 1
	for {
		if !usedNumbers[candidateNumber] {
			return fmt.Sprintf("%s%d", prefix, candidateNumber)
		}
		candidateNumber++
		// 防止无限循环
		if candidateNumber > 10000 {
			// 如果序号太大，使用最大序号+1
			return fmt.Sprintf("%s%d", prefix, maxNumber+1)
		}
	}
}

func (s *MvViewFormService) UpdateMvView(req *model.UpdateMvViewReq) error {
	var (
		view         *model.MvView
		updateFields []string
		err          error
	)

	view, err = db.Get[model.MvView](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv view error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}

	if view == nil {
		return errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	if req.Name != nil {
		view.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.Description != nil {
		view.Description = *req.Description
		updateFields = append(updateFields, "description")
	}

	if err = db.Update(db.GetDB(), view, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv view error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	return nil
}

func (s *MvViewFormService) DeleteMvView(req *model.DeleteMvViewReq) error {
	view, err := db.Get[model.MvView](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv view error", zap.Error(err))
		return errorx.InternalServerError("操作失败")
	}
	if view == nil {
		return errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err = db.Delete[model.MvView](tx, map[string]any{"id": view.ID}); err != nil {
			log.Error("delete mv view error", zap.Error(err))
			return errorx.InternalServerError("操作失败")
		}

		if view.Type == constants.ViewTypeTable {
			if err = db.Delete[model2.MvViewTable](tx, map[string]any{"view_id": view.ID}); err != nil {
				log.Error("delete mv view table error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}

		if view.Type == constants.ViewTypeForm {
			if err = db.Delete[model3.MvViewForm](tx, map[string]any{"view_id": view.ID}); err != nil {
				log.Error("delete mv view form error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}

		if view.Type == constants.ViewTypeBoard {
			if err = db.Delete[model7.MvViewBoard](tx, map[string]any{"view_id": view.ID}); err != nil {
				log.Error("delete mv view board error", zap.Error(err))
				return errorx.InternalServerError("操作失败")
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *MvViewFormService) GetMvView(req *model.GetMvViewReq) (*model.GetMvViewRes, error) {
	var (
		view *model.MvView
		err  error
	)

	view, err = db.Get[model.MvView](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv view error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if view == nil {
		return nil, errorx.New(errorx.ErrNotFound, "视图不存在")
	}

	res := &model.GetMvViewRes{
		ID:          view.ID,
		Type:        view.Type,
		Name:        view.Name,
		Description: view.Description,
	}

	// 根据视图类型获取对应的配置
	switch view.Type {
	case constants.ViewTypeTable:
		viewTable, err := db.Get[model2.MvViewTable](db.GetDB(), map[string]any{"view_id": view.ID})
		if err != nil {
			log.Error("get mv view table error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		if viewTable != nil {
			tableConfig := &model.MvViewTableConfig{
				RowHeight: viewTable.RowHeight,
			}

			// 解析 JSON 配置字段
			if len(viewTable.FilterConfig) > 0 {
				if err := json.Unmarshal(viewTable.FilterConfig, &tableConfig.FilterConfig); err != nil {
					log.Error("unmarshal filter config error", zap.Error(err))
				}
			}
			if len(viewTable.GroupConfig) > 0 {
				if err := json.Unmarshal(viewTable.GroupConfig, &tableConfig.GroupConfig); err != nil {
					log.Error("unmarshal group config error", zap.Error(err))
				}
			}
			if len(viewTable.SortConfig) > 0 {
				if err := json.Unmarshal(viewTable.SortConfig, &tableConfig.SortConfig); err != nil {
					log.Error("unmarshal sort config error", zap.Error(err))
				}
			}
			if len(viewTable.ColorConfig) > 0 {
				if err := json.Unmarshal(viewTable.ColorConfig, &tableConfig.ColorConfig); err != nil {
					log.Error("unmarshal color config error", zap.Error(err))
				}
			}

			res.TableConfig = tableConfig
		}
	case constants.ViewTypeForm:
		viewForm, err := db.Get[model3.MvViewForm](db.GetDB(), map[string]any{"view_id": view.ID})
		if err != nil {
			log.Error("get mv view form error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		if viewForm != nil {
			formConfig := &model.MvViewFormConfig{
				Name:                  viewForm.Name,
				Description:           viewForm.Description,
				Cover:                 viewForm.Cover,
				Layout:                viewForm.Layout,
				Stats:                 string(viewForm.Stats),
				EnableSharing:         viewForm.EnableSharing,
				EnableAnonymous:       viewForm.EnableAnonymous,
				Filter:                viewForm.Filter,
				EnableNoLogin:         viewForm.EnableNoLogin,
				EnableLimitSubmit:     viewForm.EnableLimitSubmit,
				LimitSubmitType:       viewForm.LimitSubmitType,
				EnableLimitCollect:    viewForm.EnableLimitCollect,
				LimitCollectCount:     viewForm.LimitCollectCount,
				EnableCycleRemind:     viewForm.EnableCycleRemind,
				EnableEditAfterSubmit: viewForm.EnableEditAfterSubmit,
			}

			// 解析 JSON 配置字段
			if len(viewForm.FilterConfig) > 0 {
				if err := json.Unmarshal(viewForm.FilterConfig, &formConfig.FilterConfig); err != nil {
					log.Error("unmarshal filter config error", zap.Error(err))
				}
			}
			if len(viewForm.CycleRemindConfig) > 0 {
				if err := json.Unmarshal(viewForm.CycleRemindConfig, &formConfig.CycleRemindConfig); err != nil {
					log.Error("unmarshal cycle remind config error", zap.Error(err))
				}
			}
			if len(viewForm.Config) > 0 {
				if err := json.Unmarshal(viewForm.Config, &formConfig.Config); err != nil {
					log.Error("unmarshal config error", zap.Error(err))
				}
			}

			res.FormConfig = formConfig
		}
	case constants.ViewTypeBoard:
		viewBoard, err := db.Get[model7.MvViewBoard](db.GetDB(), map[string]any{"view_id": view.ID})
		if err != nil {
			log.Error("get mv view board error", zap.Error(err))
			return nil, errorx.InternalServerError("获取失败")
		}
		if viewBoard != nil {
			boardConfig := &model.MvViewBoardConfig{
				ShowFieldTitle: viewBoard.ShowFieldTitle,
			}

			// 解析 JSON 配置字段
			if len(viewBoard.FilterConfig) > 0 {
				if err := json.Unmarshal(viewBoard.FilterConfig, &boardConfig.FilterConfig); err != nil {
					log.Error("unmarshal filter config error", zap.Error(err))
				}
			}
			if len(viewBoard.GroupConfig) > 0 {
				if err := json.Unmarshal(viewBoard.GroupConfig, &boardConfig.GroupConfig); err != nil {
					log.Error("unmarshal group config error", zap.Error(err))
				}
			}
			if len(viewBoard.SortConfig) > 0 {
				if err := json.Unmarshal(viewBoard.SortConfig, &boardConfig.SortConfig); err != nil {
					log.Error("unmarshal sort config error", zap.Error(err))
				}
			}

			res.BoardConfig = boardConfig
		}
	}

	return res, nil
}

func (s *MvViewFormService) SwitchActiveView(userId string, req *model.SwitchActiveViewReq) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		// 删除指定 tableSchemaId 对应的所有激活视图记录
		if err := db.Delete[model6.MvViewActive](tx, map[string]any{
			"table_schema_id": req.TableSchemaID,
			"user_id":         userId,
		}); err != nil {
			log.Error("delete existing active views error", zap.Error(err))
			return errorx.InternalServerError("切换激活视图失败")
		}

		// 新增新的激活视图记录
		activeView := &model6.MvViewActive{
			TableSchemaID: req.TableSchemaID,
			ViewID:        req.ViewID,
			UserID:        userId,
		}

		if err := db.Create(tx, activeView); err != nil {
			log.Error("create active view error", zap.Error(err))
			return errorx.InternalServerError("切换激活视图失败")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
