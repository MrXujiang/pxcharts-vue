package service

import (
	"encoding/json"
	"fmt"
	"mvtable/internal/app/mv_dashboard/model"
	chartModel "mvtable/internal/app/mv_dashboard_chart/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MvDashboardService struct{}

func NewMvDashboardService() *MvDashboardService {
	return &MvDashboardService{}
}

// CreateMvDashboard 创建仪表盘
func (s *MvDashboardService) CreateMvDashboard(req *model.CreateMvDashboardReq) (string, error) {
	dashboard := &model.MvDashboard{
		ProjectID: req.ProjectID,
		FolderID:  req.FolderID,
		Name:      req.Name,
		BgType:    req.BgType,
		BgColor:   req.BgColor,
		BgImage:   req.BgImage,
	}

	if err := db.Create(db.GetDB(), dashboard); err != nil {
		log.Error("create mv dashboard error", zap.Error(err))
		return "", errorx.InternalServerError("创建失败")
	}

	return dashboard.ID, nil
}

// UpdateMvDashboard 更新仪表盘
func (s *MvDashboardService) UpdateMvDashboard(req *model.UpdateMvDashboardReq) error {
	dashboard, err := db.Get[model.MvDashboard](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv dashboard error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if dashboard == nil {
		return errorx.New(errorx.ErrNotFound, "仪表盘不存在")
	}

	updateFields := make([]string, 0)

	if req.Name != nil {
		dashboard.Name = *req.Name
		updateFields = append(updateFields, "name")
	}
	if req.BgType != nil {
		dashboard.BgType = *req.BgType
		updateFields = append(updateFields, "bg_type")
	}
	if req.BgColor != nil {
		dashboard.BgColor = *req.BgColor
		updateFields = append(updateFields, "bg_color")
	}
	if req.BgImage != nil {
		dashboard.BgImage = *req.BgImage
		updateFields = append(updateFields, "bg_image")
	}
	if req.FolderID != nil {
		dashboard.FolderID = *req.FolderID
		updateFields = append(updateFields, "folder_id")
	}

	if len(updateFields) == 0 {
		return nil
	}

	if err = db.Update(db.GetDB(), dashboard, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv dashboard error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

// DeleteMvDashboard 删除仪表盘
func (s *MvDashboardService) DeleteMvDashboard(req *model.DeleteMvDashboardReq) error {
	if err := db.Delete[model.MvDashboard](db.GetDB(), map[string]any{"id": req.ID}); err != nil {
		log.Error("delete mv dashboard error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}

	return nil
}

// GetDashboard 获取仪表盘及其图表
func (s *MvDashboardService) GetDashboard(req *model.GetDashboardReq) (*model.GetDashboardRes, error) {
	var (
		dashboard *model.MvDashboard
		charts    []*chartModel.MvDashboardChart
		err       error
	)

	dashboard, err = db.Get[model.MvDashboard](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv dashboard error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}
	if dashboard == nil {
		return nil, errorx.New(errorx.ErrNotFound, "仪表盘不存在")
	}

	charts, _, err = db.List[chartModel.MvDashboardChart](db.GetDB(), 0, 0, map[string]any{"dashboard_id": req.ID}, []string{"created_at ASC"})
	if err != nil {
		log.Error("list mv dashboard chart error", zap.Error(err))
		return nil, errorx.InternalServerError("获取失败")
	}

	items := make([]model.DashboardChartItem, len(charts))
	for i, v := range charts {
		var cfg any
		if len(v.Config) > 0 {
			if err := json.Unmarshal(v.Config, &cfg); err != nil {
				log.Error("unmarshal dashboard chart config error", zap.Error(err))
				return nil, errorx.InternalServerError("获取失败")
			}
		}
		items[i] = model.DashboardChartItem{
			ID:            v.ID,
			TableSchemaID: v.TableSchemaID,
			Field1ID:      v.Field1ID,
			Field2ID:      v.Field2ID,
			Title:         v.Title,
			Type:          v.Type,
			Config:        cfg,
		}
	}

	return &model.GetDashboardRes{
		MvDashboard: *dashboard,
		Charts:      items,
	}, nil
}

// CopyMvDashboard 复制仪表盘
func (s *MvDashboardService) CopyMvDashboard(req *model.CopyMvDashboardReq) (*model.CopyMvDashboardRes, error) {
	// 查询原仪表盘
	sourceDashboard, err := db.Get[model.MvDashboard](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get source dashboard error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}
	if sourceDashboard == nil {
		return nil, errorx.New(errorx.ErrNotFound, "源仪表盘不存在")
	}

	// 查询原仪表盘的所有图表
	sourceCharts, _, err := db.List[chartModel.MvDashboardChart](db.GetDB(), 0, 0, map[string]any{"dashboard_id": req.ID}, []string{"created_at ASC"})
	if err != nil {
		log.Error("list source dashboard charts error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}

	// 确定新仪表盘名称
	newName := sourceDashboard.Name
	if req.Name != nil && *req.Name != "" {
		newName = *req.Name
	} else {
		// 自动生成名称：原名称_副本
		newName = fmt.Sprintf("%s_副本", sourceDashboard.Name)
	}

	// 确定新仪表盘所在文件夹：保持与原仪表盘相同
	newFolderID := sourceDashboard.FolderID

	// 校验项目下是否已存在同名仪表盘
	existDashboard, err := db.Get[model.MvDashboard](db.GetDB(), map[string]any{
		"project_id": sourceDashboard.ProjectID,
		"name":       newName,
	})
	if err != nil {
		log.Error("check exist dashboard error", zap.Error(err))
		return nil, errorx.InternalServerError("复制失败")
	}
	if existDashboard != nil {
		// 如果已存在同名，尝试添加序号
		counter := 1
		for {
			candidateName := fmt.Sprintf("%s_副本%d", sourceDashboard.Name, counter)
			existDashboard, err = db.Get[model.MvDashboard](db.GetDB(), map[string]any{
				"project_id": sourceDashboard.ProjectID,
				"name":       candidateName,
			})
			if err != nil {
				log.Error("check exist dashboard error", zap.Error(err))
				return nil, errorx.InternalServerError("复制失败")
			}
			if existDashboard == nil {
				newName = candidateName
				break
			}
			counter++
			// 防止无限循环
			if counter > 1000 {
				return nil, errorx.InternalServerError("无法生成唯一名称")
			}
		}
	}

	var newDashboardID string

	// 使用事务复制仪表盘和图表
	err = db.Transaction(func(tx *gorm.DB) error {
		// 创建新仪表盘
		newDashboard := &model.MvDashboard{
			ProjectID: sourceDashboard.ProjectID,
			FolderID:  newFolderID,
			Name:      newName,
			BgType:    sourceDashboard.BgType,
			BgColor:   sourceDashboard.BgColor,
			BgImage:   sourceDashboard.BgImage,
		}
		if err := db.Create(tx, newDashboard); err != nil {
			log.Error("create copied dashboard error", zap.Error(err))
			return errorx.InternalServerError("复制失败")
		}
		newDashboardID = newDashboard.ID

		// 复制所有图表
		if len(sourceCharts) > 0 {
			newCharts := make([]*chartModel.MvDashboardChart, 0, len(sourceCharts))
			for _, chart := range sourceCharts {
				newChart := &chartModel.MvDashboardChart{
					DashboardID:   newDashboardID,
					TableSchemaID: chart.TableSchemaID,
					Field1ID:      chart.Field1ID,
					Field2ID:      chart.Field2ID,
					Title:         chart.Title,
					Type:          chart.Type,
					Config:        chart.Config,
				}
				newCharts = append(newCharts, newChart)
			}
			if err := db.CreateBatch(tx, newCharts); err != nil {
				log.Error("create copied dashboard charts error", zap.Error(err))
				return errorx.InternalServerError("复制失败")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &model.CopyMvDashboardRes{ID: newDashboardID}, nil
}
