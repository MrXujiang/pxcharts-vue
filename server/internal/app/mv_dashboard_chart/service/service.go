package service

import (
	"mvtable/internal/app/mv_dashboard_chart/model"
	"mvtable/internal/pkg/errorx"
	"mvtable/internal/storage/db"
	"mvtable/pkg/log"

	"go.uber.org/zap"
)

type MvDashboardChartService struct{}

func NewMvDashboardChartService() *MvDashboardChartService {
	return &MvDashboardChartService{}
}

// CreateMvDashboardChart 创建仪表盘图表
func (s *MvDashboardChartService) CreateMvDashboardChart(userId string, req *model.CreateMvDashboardChartReq) (string, error) {
	chart := &model.MvDashboardChart{
		DashboardID:   req.DashboardID,
		TableSchemaID: req.TableSchemaID,
		Field1ID:      req.Field1ID,
		Field2ID:      req.Field2ID,
		Title:         req.Title,
		Type:          req.Type,
		Config:        MapToJSON(req.Config),
	}

	if err := db.Create(db.GetDB(), chart); err != nil {
		log.Error("create mv dashboard chart error", zap.Error(err))
		return "", errorx.InternalServerError("创建失败")
	}

	return chart.ID, nil
}

// UpdateMvDashboardChart 更新仪表盘图表
func (s *MvDashboardChartService) UpdateMvDashboardChart(req *model.UpdateMvDashboardChartReq) error {
	chart, err := db.Get[model.MvDashboardChart](db.GetDB(), map[string]any{"id": req.ID})
	if err != nil {
		log.Error("get mv dashboard chart error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}
	if chart == nil {
		return errorx.New(errorx.ErrNotFound, "图表不存在")
	}

	updateFields := make([]string, 0)

	if req.DashboardID != nil {
		chart.DashboardID = *req.DashboardID
		updateFields = append(updateFields, "dashboard_id")
	}
	if req.TableSchemaID != nil {
		chart.TableSchemaID = *req.TableSchemaID
		updateFields = append(updateFields, "table_schema_id")
	}
	if req.Field1ID != nil {
		chart.Field1ID = *req.Field1ID
		updateFields = append(updateFields, "field1_id")
	}
	if req.Field2ID != nil {
		chart.Field2ID = *req.Field2ID
		updateFields = append(updateFields, "field2_id")
	}
	if req.Title != nil {
		chart.Title = *req.Title
		updateFields = append(updateFields, "title")
	}
	if req.Type != nil {
		chart.Type = *req.Type
		updateFields = append(updateFields, "type")
	}
	if req.Config != nil {
		chart.Config = *req.Config
		updateFields = append(updateFields, "config")
	}

	if len(updateFields) == 0 {
		return nil
	}

	if err = db.Update(db.GetDB(), chart, map[string]any{"id": req.ID}, updateFields...); err != nil {
		log.Error("update mv dashboard chart error", zap.Error(err))
		return errorx.InternalServerError("更新失败")
	}

	return nil
}

// DeleteMvDashboardChart 删除仪表盘图表
func (s *MvDashboardChartService) DeleteMvDashboardChart(req *model.DeleteMvDashboardChartReq) error {
	if err := db.Delete[model.MvDashboardChart](db.GetDB(), map[string]any{"id": req.ID}); err != nil {
		log.Error("delete mv dashboard chart error", zap.Error(err))
		return errorx.InternalServerError("删除失败")
	}

	return nil
}
